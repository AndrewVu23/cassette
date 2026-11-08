//go:build !(darwin && cgo)

package counters

// Read has no reader outside macOS yet.
// shortcut: Linux returns nothing; add perf or cgroup reads when Linux counters are needed.
func Read(pid int) (Sample, bool) { return Sample{}, false }

// Signpost is macOS only (Instruments); a no-op elsewhere.
func Signpost(seq int, end bool) {}
