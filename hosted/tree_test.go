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
