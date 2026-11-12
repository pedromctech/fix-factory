package monitor_test

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/pedromctech/fix-factory/internal/monitor"
)

var phaseNames = []string{"scan", "plan", "implement", "test", "review", "create-pr"}

func TestParsePhase(t *testing.T) {
	for _, n := range phaseNames {
		p, err := monitor.ParsePhase(n)
		if err != nil || p.String() != n {
			t.Errorf("ParsePhase(%q) = %v, %v", n, p, err)
		}
	}
	for _, in := range []string{"", "bogus", "Scan", "PLAN", " plan"} {
		_, err := monitor.ParsePhase(in)
		var pe *monitor.ParseError
		if !errors.Is(err, monitor.ErrUnknownPhase) || !errors.As(err, &pe) {
			t.Errorf("ParsePhase(%q) err = %v", in, err)
		}
	}
}

func TestParseStatus(t *testing.T) {
	for _, n := range []string{"pending", "running", "done", "failed"} {
		s, err := monitor.ParseStatus(n)
		if err != nil || s.String() != n {
			t.Errorf("ParseStatus(%q) = %v, %v", n, s, err)
		}
	}
	for _, in := range []string{"", "bogus", "Done", "RUNNING"} {
		_, err := monitor.ParseStatus(in)
		var pe *monitor.ParseError
		if !errors.Is(err, monitor.ErrUnknownStatus) || !errors.As(err, &pe) {
			t.Errorf("ParseStatus(%q) err = %v", in, err)
		}
	}
}

func TestPhasesOrder(t *testing.T) {
	ps := monitor.Phases()
	if len(ps) != 6 {
		t.Fatalf("got %d phases", len(ps))
	}
	for i, p := range ps {
		if p.String() != phaseNames[i] {
			t.Errorf("phase %d = %s, want %s", i, p, phaseNames[i])
		}
	}
}

func mustPhase(t *testing.T, n string) monitor.Phase {
	t.Helper()
	p, err := monitor.ParsePhase(n)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func mustStatus(t *testing.T, n string) monitor.Status {
	t.Helper()
	s, err := monitor.ParseStatus(n)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestNewTrackerAllPending(t *testing.T) {
	snap := monitor.NewTracker().Snapshot()
	pending := mustStatus(t, "pending")
	for _, p := range monitor.Phases() {
		if snap[p] != pending {
			t.Errorf("%s = %s", p, snap[p])
		}
	}
}

func TestSetUpdatesOnlyOnePhase(t *testing.T) {
	tr := monitor.NewTracker()
	if err := tr.Set(mustPhase(t, "plan"), mustStatus(t, "done")); err != nil {
		t.Fatal(err)
	}
	snap := tr.Snapshot()
	for _, p := range monitor.Phases() {
		want := "pending"
		if p.String() == "plan" {
			want = "done"
		}
		if snap[p].String() != want {
			t.Errorf("%s = %s, want %s", p, snap[p], want)
		}
	}
}

func TestSetRejectsInvalid(t *testing.T) {
	tr := monitor.NewTracker()
	var badP monitor.Phase
	var badS monitor.Status
	if err := tr.Set(badP, mustStatus(t, "done")); err == nil {
		t.Error("zero phase accepted")
	}
	if err := tr.Set(mustPhase(t, "scan"), badS); err == nil {
		t.Error("zero status accepted")
	}
	if err := tr.Set(monitor.Phase(99), mustStatus(t, "done")); err == nil {
		t.Error("out of range phase accepted")
	}
	snap := tr.Snapshot()
	for _, p := range monitor.Phases() {
		if snap[p].String() != "pending" {
			t.Errorf("%s changed to %s", p, snap[p])
		}
	}
}

func TestSnapshotIsCopy(t *testing.T) {
	tr := monitor.NewTracker()
	snap := tr.Snapshot()
	scan := mustPhase(t, "scan")
	snap[scan] = mustStatus(t, "failed")
	if tr.Snapshot()[scan].String() != "pending" {
		t.Error("snapshot mutation leaked")
	}
}

func TestRenderMixedState(t *testing.T) {
	tr := monitor.NewTracker()
	_ = tr.Set(mustPhase(t, "scan"), mustStatus(t, "done"))
	_ = tr.Set(mustPhase(t, "plan"), mustStatus(t, "running"))
	_ = tr.Set(mustPhase(t, "test"), mustStatus(t, "failed"))
	var buf bytes.Buffer
	if err := monitor.Render(&buf, tr.Snapshot()); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 6 {
		t.Fatalf("got %d lines:\n%s", len(lines), buf.String())
	}
	want := []string{"done", "running", "pending", "failed", "pending", "pending"}
	for i, l := range lines {
		if !strings.Contains(l, phaseNames[i]) || !strings.HasSuffix(l, want[i]) {
			t.Errorf("line %d = %q", i, l)
		}
	}
}

func TestConcurrentSetAndSnapshot(t *testing.T) {
	tr := monitor.NewTracker()
	p, s := mustPhase(t, "test"), mustStatus(t, "running")
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); _ = tr.Set(p, s) }()
		go func() { defer wg.Done(); _ = tr.Snapshot() }()
	}
	wg.Wait()
}
