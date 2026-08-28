package main

import (
	"fmt"
	"sync"
)

// Decl is a desired-state declaration: "name must have `desired` live
// instances". The reconciler loop converges reality to this declaration.
type Decl struct {
	name    string
	desired int

	mu      sync.Mutex
	next    int
	alive   []string
	history []string
	subs    map[chan string]bool
}

func newDecl(name string, desired int) *Decl {
	return &Decl{name: name, desired: desired, subs: map[chan string]bool{}}
}

func (d *Decl) count() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.alive)
}

func (d *Decl) announce(msg string) {
	d.mu.Lock()
	d.history = append(d.history, msg)
	if len(d.history) > 100 {
		d.history = d.history[len(d.history)-100:]
	}
	subs := d.subs
	d.mu.Unlock()
	for ch := range subs {
		select {
		case ch <- msg:
		default:
		}
	}
}

func (d *Decl) levels() (alive []string, wanted int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	cp := make([]string, len(d.alive))
	copy(cp, d.alive)
	return cp, d.desired
}

func (d *Decl) spawn() {
	d.mu.Lock()
	d.next++
	id := fmt.Sprintf("%s-%d", d.name, d.next)
	d.alive = append(d.alive, id)
	d.mu.Unlock()
	d.announce("spawn " + id)
}

func (d *Decl) reap() {
	d.mu.Lock()
	if len(d.alive) == 0 {
		d.mu.Unlock()
		return
	}
	id := d.alive[len(d.alive)-1]
	d.alive = d.alive[:len(d.alive)-1]
	d.mu.Unlock()
	d.announce("reap " + id)
}

// Simulate an instance dying unexpectedly. The reconciler respawns it on
// the next tick — this is the self-healing loop (K8s liveness equivalent).
func (d *Decl) kill() {
	d.mu.Lock()
	if len(d.alive) == 0 {
		d.mu.Unlock()
		return
	}
	id := d.alive[len(d.alive)-1]
	d.alive = d.alive[:len(d.alive)-1]
	d.mu.Unlock()
	d.announce("CRASH " + id)
}

// reconcile converges reality (alive count) to the desired state.
func (d *Decl) reconcile() {
	for {
		alive, wanted := d.levels()
		switch {
		case len(alive) < wanted:
			d.spawn()
		case len(alive) > wanted:
			d.reap()
		default:
			return
		}
	}
}

func (d *Decl) subscribe() chan string {
	ch := make(chan string, 32)
	d.mu.Lock()
	d.subs[ch] = true
	d.mu.Unlock()
	return ch
}

func (d *Decl) unsubscribe(ch chan string) {
	d.mu.Lock()
	delete(d.subs, ch)
	d.mu.Unlock()
}
