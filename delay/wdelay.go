//go:build windows

/* delay
Windows code for timing delay. Only included when OS is windows.
*/
package delay

import (
	"syscall"
	"time"
	"unsafe"
)

// precision timing
var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	procFreq    = modkernel32.NewProc("QueryPerformanceFrequency")
	procCounter = modkernel32.NewProc("QueryPerformanceCounter")

	qpcFrequency = getFrequency()
	qpcBase      = getCount()
)

// getFrequency returns frequency in ticks per second.
func getFrequency() int64 {
	var freq int64
	r1, _, _ := syscall.Syscall(procFreq.Addr(), 1, uintptr(unsafe.Pointer(&freq)), 0, 0)
	if r1 == 0 {
		panic("call failed")
	}
	return freq
}

// getCount returns counter ticks.
func getCount() int64 {
	var qpc int64
	syscall.Syscall(procCounter.Addr(), 1, uintptr(unsafe.Pointer(&qpc)), 0, 0)
	return qpc
}

// Now returns current time.Duration with best possible precision.
//
// Now returns time offset from a specific time.
// The values aren't comparable between computer restarts or between computers.
func now() time.Duration {
	return time.Duration(getCount()-qpcBase) * time.Second / (time.Duration(qpcFrequency) * time.Nanosecond)
}

// NowPrecision returns maximum possible precision for Now in nanoseconds.
//func nowPrecision() float64 {
//	return float64(time.Second) / (float64(qpcFrequency) * float64(time.Nanosecond))
//}

// Delay : delay for d nanoseconds
func Delay(d int64) {
	var st, rn int64
	st = int64(now())
	for {
		rn = int64(now())
		if rn-st > d {
			return
		}
	}
}
