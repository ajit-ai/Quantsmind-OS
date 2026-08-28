//go:build darwin || freebsd || openbsd || netbsd

package main

// Real collectors for these OSes not yet implemented; -1 signals the
// caller to substitute simulated values (marked clearly as simulated).
func sampleSystem() (cpuPct int, memPct int, source string) {
	return -1, -1, "simulated"
}
