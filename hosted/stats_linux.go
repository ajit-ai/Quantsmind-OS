//go:build linux

package main

import (
	"os"
	"strconv"
	"strings"
)

type linuxCpuSampler struct {
	started                 bool
	prevTotal, prevIdle     uint64
}

var lcpu linuxCpuSampler

func sampleSystem() (cpuPct int, memPct int, source string) {
	// --- CPU from /proc/stat ("cpu  user nice system idle iowait irq softirq steal")
	b, err := os.ReadFile("/proc/stat")
	if err != nil {
		return -1, -1, "linux"
	}
	lines := strings.Split(string(b), "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "cpu ") {
		return -1, -1, "linux"
	}
	f := strings.Fields(lines[0])
	var vals [8]uint64
	for i := 1; i < len(f) && i < 9; i++ {
		v, _ := strconv.ParseUint(f[i], 10, 64)
		vals[i-1] = v
	}
	idle := vals[3] + vals[4] // idle + iowait
	total := uint64(0)
	for i := 0; i < 8; i++ {
		total += vals[i]
	}
	cpuPct = 0
	if lcpu.started {
		dTotal := total - lcpu.prevTotal
		dIdle := idle - lcpu.prevIdle
		if dTotal > 0 {
			cpuPct = int(float64(dTotal-dIdle)*100/float64(dTotal) + 0.5)
			if cpuPct > 100 {
				cpuPct = 100
			}
			if cpuPct < 0 {
				cpuPct = 0
			}
		}
	}
	lcpu.prevTotal, lcpu.prevIdle, lcpu.started = total, idle, true

	// --- Memory from /proc/meminfo
	mi, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return cpuPct, -1, "linux"
	}
	totalMem, avail := uint64(0), uint64(0)
	for _, line := range strings.Split(string(mi), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		v, _ := strconv.ParseUint(fields[1], 10, 64)
		switch fields[0] {
		case "MemTotal:":
			totalMem = v
		case "MemAvailable:":
			avail = v
		}
	}
	memPct = -1
	if totalMem > 0 && avail > 0 {
		memPct = int(float64(totalMem-avail)*100/float64(totalMem) + 0.5)
	}
	return cpuPct, memPct, "linux"
}
