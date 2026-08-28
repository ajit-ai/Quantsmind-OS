//go:build windows

package main

import (
	"sync"
	"syscall"
	"unsafe"
)

type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

type filetime struct {
	Low  uint32
	High uint32
}

var (
	procGlobalMemoryStatusEx = syscall.NewLazyDLL("kernel32.dll").NewProc("GlobalMemoryStatusEx")
	procGetSystemTimes       = syscall.NewLazyDLL("kernel32.dll").NewProc("GetSystemTimes")
)

func memUsedPercent() int {
	m := &memoryStatusEx{Length: uint32(unsafe.Sizeof(memoryStatusEx{}))}
	r1, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(m)))
	if r1 == 0 {
		return -1
	}
	return int(m.MemoryLoad)
}

type cpuSampler struct {
	mu                      sync.Mutex
	started                 bool
	prevIdle, prevKernel    uint64
	prevUser                uint64
}

var sampler cpuSampler

// sample returns current CPU usage percent via GetSystemTimes deltas.
func (s *cpuSampler) sample() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	var idle, kernel, user filetime
	r1, _, _ := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&idle)),
		uintptr(unsafe.Pointer(&kernel)),
		uintptr(unsafe.Pointer(&user)),
	)
	if r1 == 0 {
		return -1
	}
	toU64 := func(ft filetime) uint64 { return uint64(ft.High)<<32 | uint64(ft.Low) }
	idleN := toU64(idle)
	kerN := toU64(kernel)
	usrN := toU64(user)

	if !s.started {
		s.prevIdle, s.prevKernel, s.prevUser = idleN, kerN, usrN
		s.started = true
		return 0
	}

	total := (kerN + usrN) - (s.prevKernel + s.prevUser)
	idleD := idleN - s.prevIdle
	s.prevIdle, s.prevKernel, s.prevUser = idleN, kerN, usrN

	if total <= 0 {
		return 0
	}
	busy := total - idleD
	pct := int(float64(busy)/float64(total)*100 + 0.5)
	if pct > 100 {
		pct = 100
	}
	if pct < 0 {
		pct = 0
	}
	return pct
}

func sampleSystem() (cpuPct int, memPct int, source string) {
	return sampler.sample(), memUsedPercent(), "real"
}
