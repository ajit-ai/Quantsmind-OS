package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Node is the universal Quantis resource: readable, writable, watchable.
// In hosted form it previews what every kernel object will expose via 9P++.
type Node struct {
	path string
	kind string // "", "metrics", "net", "infer"

	mu   sync.Mutex
	data string
	subs map[chan string]bool

	// Event fabric: a bounded replay ring so late subscribers catch up,
	// plus explicit backpressure accounting (slow subscribers are dropped,
	// never allowed to block the writers).
	history     []string // most recent first
	historyCap  int
	dropped     uint64
	lastDropped string
}

// ReplayMax is how many past events a freshly subscribed watcher receives.
const replayBuffer = 32

func newNode(kind string) *Node {
	return &Node{
		kind:       kind,
		subs:       map[chan string]bool{},
		historyCap: replayBuffer,
	}
}

func (n *Node) emit(ev string) {
	n.mu.Lock()
	// append to the replay ring (drop the oldest once full)
	n.history = append([]string{ev}, n.history...)
	if len(n.history) > n.historyCap {
		n.history = n.history[:n.historyCap]
	}
	subs := n.subs
	// snapshot subscriber handles before releasing the lock; delivery
	// happens outside it so a slow consumer never stalls the writer.
	handles := make([]chan string, 0, len(subs))
	for ch := range subs {
		handles = append(handles, ch)
	}
	n.mu.Unlock()

	for _, ch := range handles {
		select {
		case ch <- ev:
		default:
			// Backpressure policy: drop for slow subscribers, but count it
			// so operators can see the pressure (never block the writer).
			n.countDrop(ev)
		}
	}
}

func (n *Node) countDrop(ev string) {
	n.mu.Lock()
	n.dropped++
	n.lastDropped = ev
	n.mu.Unlock()
}

// Subscribe registers a watcher. When replay=true the subscriber first
// receives the recent event history so a late joiner doesn't miss state.
func (n *Node) SubscribeReplay(replay bool) chan string {
	ch := make(chan string, 16)
	n.mu.Lock()
	n.subs[ch] = true
	hist := n.history
	n.mu.Unlock()
	if replay {
		for _, ev := range hist {
			select {
			case ch <- ev:
			default:
			}
		}
	}
	return ch
}

// Subscribe is the lightweight variant: no replay, just live events.
func (n *Node) Subscribe() chan string {
	return n.SubscribeReplay(false)
}

func (n *Node) Unsubscribe(ch chan string) {
	n.mu.Lock()
	delete(n.subs, ch)
	n.mu.Unlock()
}

// Stats summarizes event-fabric pressure for introspection.
func (n *Node) stats() (ring int, dropped uint64, subs int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.history), n.dropped, len(n.subs)
}

func (n *Node) Read() ([]byte, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return []byte(n.data), nil
}

func (n *Node) Write(b []byte) error {
	n.mu.Lock()
	n.data = string(b)
	kind := n.kind
	path := n.path
	n.mu.Unlock()

	if kind == "infer" {
		if strings.Contains(path, "medscan") {
			go func() {
				res := inferMedScan(path, string(b), n)
				// last write wins: cache the structured result on the node
				n.mu.Lock()
				n.data = fmt.Sprintf("triage=%s risk=%d/100 findings=%v\n",
					res.Triage, res.RiskScore, res.Findings)
				n.mu.Unlock()
			}()
			return nil
		}
		go streamTokens(path, string(b), n)
		return nil
	}
	n.emit(fmt.Sprintf("write %s: %s", path, trunc(string(b))))
	return nil
}

func streamTokens(path, prompt string, n *Node) {
	for _, tok := range strings.Fields(prompt) {
		n.emit(fmt.Sprintf("token %s/%s: %s", path, tok, tok))
		time.Sleep(60 * time.Millisecond)
	}
	n.emit(fmt.Sprintf("done %s", path))
}

func trunc(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 80 {
		return s[:80] + "..."
	}
	return s
}
