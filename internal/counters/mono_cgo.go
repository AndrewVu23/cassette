//go:build (darwin || linux) && cgo

package counters

/*
#include <time.h>
*/
import "C"

// MonoNs returns CLOCK_MONOTONIC in ns, so marks line up with other monotonic logs.
func MonoNs() int64 {
	var ts C.struct_timespec
	C.clock_gettime(C.CLOCK_MONOTONIC, &ts)
	return int64(ts.tv_sec)*1e9 + int64(ts.tv_nsec)
}
