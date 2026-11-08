//go:build !((darwin || linux) && cgo)

package counters

import "time"

// MonoNs falls back to wall-clock ns where cgo is off (e.g. Windows CI).
// shortcut: not CLOCK_MONOTONIC, so marks can't be lined up with other logs here; measure on macOS or Linux.
func MonoNs() int64 { return time.Now().UnixNano() }
