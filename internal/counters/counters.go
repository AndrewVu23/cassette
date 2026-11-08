// Package counters reads a process's CPU counters at turn boundaries.
package counters

// Sample is one reading of a process's cumulative counters.
type Sample struct {
	Instructions  uint64 `json:"instr"`
	Cycles        uint64 `json:"cycles"`
	PInstructions uint64 `json:"p_instr"`  // part of Instructions run on P-cores
	PCycles       uint64 `json:"p_cycles"` // part of Cycles run on P-cores
	EnergyNJ      uint64 `json:"energy_nj"`
	UserNs        uint64 `json:"user_ns"`
	SysNs         uint64 `json:"sys_ns"`
	ChildUserNs   uint64 `json:"child_user_ns"` // only children that already exited
	ChildSysNs    uint64 `json:"child_sys_ns"`
	FootprintB    uint64 `json:"footprint_bytes"`
}
