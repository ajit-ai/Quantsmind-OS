package main

import "testing"

func TestNodeReadWrite(t *testing.T) {
	n := newNode("")
	n.Write([]byte("hello"))
	data, err := n.Read()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("expected 'hello', got %q", string(data))
	}
}

func TestNodeWriteEmitsEvent(t *testing.T) {
	n := newNode("")
	ch := n.Subscribe()
	defer n.Unsubscribe(ch)

	if err := n.Write([]byte("ping")); err != nil {
		t.Fatal(err)
	}
	ev := <-ch
	if ev == "" {
		t.Fatal("expected a watch event after write")
	}
}

func TestNodeMultipleSubscribers(t *testing.T) {
	n := newNode("")
	c1 := n.Subscribe()
	c2 := n.Subscribe()
	defer n.Unsubscribe(c1)
	defer n.Unsubscribe(c2)

	n.Write([]byte("fanout"))
	<-c1
	<-c2
}

func TestTreeMountAndIntrospection(t *testing.T) {
	tr := NewTree()
	tr.Mount("/sys/cpu", newNode("metrics"))
	tr.Mount("/dev/net0", newNode("net"))

	if _, ok := tr.get("/sys/cpu"); !ok {
		t.Fatal("expected /sys/cpu to be mountable")
	}
	if _, ok := tr.get("/dev/net0"); !ok {
		t.Fatal("expected /dev/net0 to be mountable")
	}
	if n, ok := tr.get("/missing"); ok || n != nil {
		t.Fatal("expected /missing to be absent")
	}
}

func TestTreeNeverMutatesAfterMount(t *testing.T) {
	tr := NewTree()
	tr.Mount("/a", newNode(""))
	n, ok := tr.get("/a")
	if !ok || n == nil {
		t.Fatal("expected to find mounted node")
	}
	if n.path != "/a" {
		t.Fatalf("node path mismatch: %q", n.path)
	}
}

func TestReplaySeedsLateSubscriber(t *testing.T) {
	n := newNode("")
	// emit a few events before anyone subscribes
	n.emit("event-1")
	n.emit("event-2")
	n.emit("event-3")

	ch := n.SubscribeReplay(true)
	defer n.Unsubscribe(ch)

	// replay should deliver the recent history, newest first
	got := map[string]bool{}
	for i := 0; i < 3; i++ {
		got[<-ch] = true
	}
	if !got["event-3"] || !got["event-2"] || !got["event-1"] {
		t.Fatalf("replay missing recent events: %v", got)
	}
}

func TestReplayRingIsBounded(t *testing.T) {
	n := newNode("")
	for i := 0; i < replayBuffer+50; i++ {
		n.emit("ev")
	}
	ring, _, _ := n.stats()
	if ring > replayBuffer {
		t.Fatalf("replay ring exceeded cap: %d > %d", ring, replayBuffer)
	}
}

func TestBackpressureCountsDropped(t *testing.T) {
	n := newNode("")
	// a full channel = a slow subscriber that can't keep up
	ch := n.Subscribe()
	defer n.Unsubscribe(ch)
	for i := 0; i < cap(ch)+5; i++ {
		n.emit("busy")
	}
	_, dropped, _ := n.stats()
	if dropped == 0 {
		t.Fatal("expected backpressure to count dropped events")
	}
}

func TestFabricStatsAggregate(t *testing.T) {
	tr := NewTree()
	a := newNode("")
	b := newNode("")
	tr.Mount("/a", a)
	tr.Mount("/b", b)

	// subscribe to a so it has a watcher; b stays idle
	chA := a.Subscribe()
	defer a.Unsubscribe(chA)
	a.emit("x")

	fs := tr.FabricStats()
	if fs.Watchers != 1 {
		t.Fatalf("expected 1 watcher aggregated, got %d", fs.Watchers)
	}
	if fs.ReplayRing < 1 {
		t.Fatalf("expected replayed events in aggregate ring, got %d", fs.ReplayRing)
	}
	_ = b
}

