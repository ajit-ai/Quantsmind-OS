package main

import "testing"

func TestReconcileScalesUp(t *testing.T) {
	d := newDecl("svc", 3)
	d.reconcile()
	if d.count() != 3 {
		t.Fatalf("expected 3 alive, got %d", d.count())
	}
}

func TestReconcileScalesDown(t *testing.T) {
	d := newDecl("svc", 5)
	d.reconcile()
	if d.count() != 5 {
		t.Fatalf("setup: expected 5, got %d", d.count())
	}
	d.desired = 2
	d.reconcile()
	if d.count() != 2 {
		t.Fatalf("expected 2 after scale-down, got %d", d.count())
	}
}

func TestSelfHealAfterKill(t *testing.T) {
	d := newDecl("svc", 3)
	d.reconcile()
	aliveBefore, _ := d.levels()
	if len(aliveBefore) != 3 {
		t.Fatalf("setup: expected 3, got %d", len(aliveBefore))
	}
	// simulate a crash
	d.kill()
	if got := d.count(); got != 2 {
		t.Fatalf("after kill expected 2, got %d", got)
	}
	// reconciler converges back to desired
	d.reconcile()
	if got := d.count(); got != 3 {
		t.Fatalf("after reconcile expected self-heal to 3, got %d", got)
	}
}

func TestInstanceIDsAreUnique(t *testing.T) {
	d := newDecl("svc", 10)
	d.reconcile()
	alive, wanted := d.levels()
	if len(alive) != wanted {
		t.Fatalf("expected %d alive, got %d", wanted, len(alive))
	}
	seen := map[string]bool{}
	for _, id := range alive {
		if seen[id] {
			t.Fatalf("duplicate instance id: %s", id)
		}
		seen[id] = true
	}
}

func TestReconcileIdempotent(t *testing.T) {
	d := newDecl("svc", 4)
	d.reconcile()
	d.reconcile()
	d.reconcile()
	if got := d.count(); got != 4 {
		t.Fatalf("reconcile must be idempotent; got %d", got)
	}
}

func TestKillEmptyDoesNotPanic(t *testing.T) {
	d := newDecl("svc", 0)
	d.kill() // must not panic
	if got := d.count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
