package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// Registry is the set of live declarations, introspectable like a namespace.
type Registry struct {
	mu    sync.Mutex
	decls map[string]*Decl
}

func newRegistry() *Registry { return &Registry{decls: map[string]*Decl{}} }

func (r *Registry) get(name string) (*Decl, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, ok := r.decls[name]
	return d, ok
}

func (r *Registry) put(d *Decl) {
	r.mu.Lock()
	r.decls[d.name] = d
	r.mu.Unlock()
}

func (r *Registry) list() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	names := make([]string, 0, len(r.decls))
	for n := range r.decls {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

type declRequest struct {
	Name    string `json:"name"`
	Desired int    `json:"desired"`
}

func main() {
	reg := newRegistry()

	// governance loop: continuously converge every declaration
	go func() {
		for range time.Tick(500 * time.Millisecond) {
			for _, name := range reg.list() {
				if d, ok := reg.get(name); ok {
					d.reconcile()
				}
			}
		}
	}()

	mux := http.NewServeMux()

	// POST /decl  {"name":"svc","desired":3}
	mux.HandleFunc("/decl", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req declRequest
			dec := json.NewDecoder(r.Body)
			if err := dec.Decode(&req); err != nil || req.Name == "" {
				http.Error(w, "bad body", http.StatusBadRequest)
				return
			}
			if req.Desired < 0 {
				req.Desired = 0
			}
			reg.put(newDecl(req.Name, req.Desired))
			fmt.Fprintf(w, "declared %s → %d desired\n", req.Name, req.Desired)
		case http.MethodGet:
			fmt.Fprintln(w, strings.Join(reg.list(), "\n"))
		default:
			http.Error(w, "method", http.StatusMethodNotAllowed)
		}
	})

	// /decl's detailed status, kill, or SSE watch
	mux.HandleFunc("/decl/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/decl/")
		name = strings.Trim(name, "/")
		if name == "" {
			http.Error(w, "missing decl name", http.StatusBadRequest)
			return
		}

		if strings.HasSuffix(name, "/watch") {
			name = strings.TrimSuffix(name, "/watch")
			d, ok := reg.get(name)
			if !ok {
				http.Error(w, "no such decl", http.StatusNotFound)
				return
			}
			streamEvents(w, r, d)
			return
		}
		if strings.HasSuffix(name, "/kill") {
			name = strings.TrimSuffix(name, "/kill")
			d, ok := reg.get(name)
			if !ok {
				http.Error(w, "no such decl", http.StatusNotFound)
				return
			}
			d.kill()
			fmt.Fprintf(w, "killed one instance of %s; reconciler will respawn\n", name)
			return
		}

		d, ok := reg.get(name)
		if !ok {
			http.Error(w, "no such decl", http.StatusNotFound)
			return
		}
		alive, wanted := d.levels()
		status := map[string]interface{}{
			"name":    name,
			"desired": wanted,
			"alive":   len(alive),
			"instances": alive,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	addr := ":9091"
	log.Printf("reconciler listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func streamEvents(w http.ResponseWriter, r *http.Request, d *Decl) {
	fl, _ := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	ch := d.subscribe()
	defer d.unsubscribe(ch)
	fmt.Fprintf(w, "event hello\ndata watching %s\n\n", d.name)
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
