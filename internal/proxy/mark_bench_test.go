package proxy

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// BenchmarkMark measures one --marks line (counter read + JSON + file write),
// the delay a "req" mark adds before each replayed reply.
func BenchmarkMark(b *testing.B) {
	f, err := os.Create(filepath.Join(b.TempDir(), "marks.jsonl"))
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	// Read another process, as in a real run.
	child := exec.Command("sleep", "60")
	if err := child.Start(); err != nil {
		b.Fatal(err)
	}
	defer child.Process.Kill()

	p := &Proxy{marks: json.NewEncoder(f)}
	p.SetPID(child.Process.Pid)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.mark("req", 0)
	}
}
