//go:build darwin && cgo

package counters

/*
#include <libproc.h>
#include <mach/mach_time.h>
#include <os/log.h>
#include <os/signpost.h>

static os_log_t poi;

// A "Points of Interest" signpost, so Instruments traces show each mark on their own timeline.
static void signpost(int seq, int end) {
	if (!poi) poi = os_log_create("cassette", OS_LOG_CATEGORY_POINTS_OF_INTEREST);
	if (end) os_signpost_event_emit(poi, OS_SIGNPOST_ID_EXCLUSIVE, "end", "seq %d", seq);
	else os_signpost_event_emit(poi, OS_SIGNPOST_ID_EXCLUSIVE, "req", "seq %d", seq);
}
*/
import "C"

import "unsafe"

// Read returns pid's counters via proc_pid_rusage, which needs no root for
// our own processes. ok is false once the process is gone.
func Read(pid int) (s Sample, ok bool) {
	var ri C.struct_rusage_info_v6
	if C.proc_pid_rusage(C.int(pid), C.RUSAGE_INFO_V6, (*C.rusage_info_t)(unsafe.Pointer(&ri))) != 0 {
		return s, false
	}
	// CPU times are in Mach ticks (about 41.67 ns on Apple Silicon); convert to ns.
	var tb C.mach_timebase_info_data_t
	C.mach_timebase_info(&tb)
	ns := func(t C.uint64_t) uint64 { return uint64(t) * uint64(tb.numer) / uint64(tb.denom) }
	return Sample{
		Instructions:  uint64(ri.ri_instructions),
		Cycles:        uint64(ri.ri_cycles),
		PInstructions: uint64(ri.ri_pinstructions),
		PCycles:       uint64(ri.ri_pcycles),
		EnergyNJ:      uint64(ri.ri_energy_nj),
		UserNs:        ns(ri.ri_user_time),
		SysNs:         ns(ri.ri_system_time),
		ChildUserNs:   ns(ri.ri_child_user_time),
		ChildSysNs:    ns(ri.ri_child_system_time),
		FootprintB:    uint64(ri.ri_phys_footprint),
	}, true
}

// Signpost marks a boundary for Instruments (xctrace) recordings. Call it from one goroutine at a time.
func Signpost(seq int, end bool) {
	e := 0
	if end {
		e = 1
	}
	C.signpost(C.int(seq), C.int(e))
}
