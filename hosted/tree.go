package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
)

// Tree is the namespace: paths -> nodes. Introspection is mandatory:
// GET /ns lists everything, making services self-describing (P2/P6).
type Tree struct {
	mu    sync.RWMutex
	nodes map[string]*Node
}

func NewTree() *Tree { return &Tree{nodes: map[string]*Node{}} }

func (t *Tree) Mount(path string, n *Node) {
	t.mu.Lock()
	defer t.mu.Unlock()
	n.path = path
	t.nodes[path] = n
}

func (t *Tree) get(path string) (*Node, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	n, ok := t.nodes[path]
	return n, ok
}

func (t *Tree) HandleNamespace(w http.ResponseWriter, r *http.Request) {
	t.mu.RLock()
	paths := make([]string, 0, len(t.nodes))
	for p := range t.nodes {
		paths = append(paths, p)
	}
	t.mu.RUnlock()
	sort.Strings(paths)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paths)
}

func (t *Tree) HandleRead(w http.ResponseWriter, path string) {
	n, ok := t.get(path)
	if !ok {
		http.Error(w, "not found: "+path, http.StatusNotFound)
		return
	}
	data, err := n.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (t *Tree) HandleWrite(w http.ResponseWriter, r *http.Request, path string) {
	n, ok := t.get(path)
	if !ok {
		http.Error(w, "not found: "+path, http.StatusNotFound)
		return
	}
	buf := make([]byte, r.ContentLength)
	r.Body.Read(buf)
	if err := n.Write(buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "ok %s (%d bytes)\n", path, len(buf))
}

// HandleWatch streams node events as Server-Sent Events (SSE),
// previewing the kernel event fabric's delivery over 9P++.
func (t *Tree) HandleWatch(w http.ResponseWriter, r *http.Request, path string) {
	n, ok := t.get(path)
	if !ok {
		http.Error(w, "not found: "+path, http.StatusNotFound)
		return
	}
	fl, _ := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	ch := n.Subscribe()
	defer n.Unsubscribe(ch)

	fmt.Fprintf(w, "event hello\ndata watching %s\n\n", path)
	if fl != nil {
		fl.Flush()
	}
	for {
		select {
		case <-r.Context().Done():
			return
		case ev := <-ch:
			fmt.Fprintf(w, "event update\ndata %s\n\n", ev)
			if fl != nil {
				fl.Flush()
			}
		}
	}
}
