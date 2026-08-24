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
}

func newNode(kind string) *Node {
	return &Node{kind: kind, subs: map[chan string]bool{}}
}

func (n *Node) emit(ev string) {
	for ch := range n.subs {
		select {
		case ch <- ev:
		default: // slow subscriber: drop rather than block writers
		}
	}
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

func (n *Node) Subscribe() chan string {
	ch := make(chan string, 16)
	n.mu.Lock()
	n.subs[ch] = true
	n.mu.Unlock()
	return ch
}

func (n *Node) Unsubscribe(ch chan string) {
	n.mu.Lock()
	delete(n.subs, ch)
	n.mu.Unlock()
}

func trunc(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 80 {
		return s[:80] + "..."
	}
	return s
}
