// Package monitor models the phases of a factory run and renders their status.
package monitor

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

// Sentinel errors returned (wrapped in *ParseError) by the parsers.
var (
	ErrUnknownPhase  = errors.New("unknown phase")
	ErrUnknownStatus = errors.New("unknown status")
)

// ParseError reports an input that is not a valid phase or status.
type ParseError struct {
	Input string
	Err   error
}

func (e *ParseError) Error() string { return fmt.Sprintf("%v: %q", e.Err, e.Input) }

// Unwrap returns the sentinel error.
func (e *ParseError) Unwrap() error { return e.Err }

// Phase is one step of a factory run. The zero value is invalid.
type Phase int

const (
	PhaseScan Phase = iota + 1
	PhasePlan
	PhaseImplement
	PhaseTest
	PhaseReview
	PhaseCreatePR
)

var phaseNames = map[Phase]string{
	PhaseScan:      "scan",
	PhasePlan:      "plan",
	PhaseImplement: "implement",
	PhaseTest:      "test",
	PhaseReview:    "review",
	PhaseCreatePR:  "create-pr",
}

// Phases returns all phases in run order.
func Phases() []Phase {
	return []Phase{PhaseScan, PhasePlan, PhaseImplement, PhaseTest, PhaseReview, PhaseCreatePR}
}

// String returns the phase name.
func (p Phase) String() string {
	if n, ok := phaseNames[p]; ok {
		return n
	}
	return "invalid"
}

// ParsePhase validates a phase name.
func ParsePhase(s string) (Phase, error) {
	for p, n := range phaseNames {
		if n == s {
			return p, nil
		}
	}
	return 0, &ParseError{Input: s, Err: ErrUnknownPhase}
}

// Status is the state of a phase. The zero value is invalid.
type Status int

const (
	StatusPending Status = iota + 1
	StatusRunning
	StatusDone
	StatusFailed
)

var statusNames = map[Status]string{
	StatusPending: "pending",
	StatusRunning: "running",
	StatusDone:    "done",
	StatusFailed:  "failed",
}

var statusMarks = map[Status]string{
	StatusPending: " ",
	StatusRunning: "~",
	StatusDone:    "x",
	StatusFailed:  "!",
}

// String returns the status name.
func (s Status) String() string {
	if n, ok := statusNames[s]; ok {
		return n
	}
	return "invalid"
}

// ParseStatus validates a status name.
func ParseStatus(s string) (Status, error) {
	for st, n := range statusNames {
		if n == s {
			return st, nil
		}
	}
	return 0, &ParseError{Input: s, Err: ErrUnknownStatus}
}

// Snapshot is a copy of the status of every phase.
type Snapshot map[Phase]Status

// Tracker holds the status of each phase. It is safe for concurrent use.
type Tracker struct {
	mu     sync.Mutex
	status map[Phase]Status
}

// NewTracker returns a tracker with every phase pending.
func NewTracker() *Tracker {
	t := &Tracker{status: make(map[Phase]Status)}
	for _, p := range Phases() {
		t.status[p] = StatusPending
	}
	return t
}

// Set updates one phase. Invalid phases or statuses are rejected.
func (t *Tracker) Set(p Phase, s Status) error {
	if _, ok := phaseNames[p]; !ok {
		return &ParseError{Input: fmt.Sprint(int(p)), Err: ErrUnknownPhase}
	}
	if _, ok := statusNames[s]; !ok {
		return &ParseError{Input: fmt.Sprint(int(s)), Err: ErrUnknownStatus}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.status[p] = s
	return nil
}

// Snapshot returns a copy of the current state.
func (t *Tracker) Snapshot() Snapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	snap := make(Snapshot, len(t.status))
	for p, s := range t.status {
		snap[p] = s
	}
	return snap
}

// Render writes one line per phase, in order.
func Render(w io.Writer, snap Snapshot) error {
	for _, p := range Phases() {
		s, ok := snap[p]
		if !ok {
			s = StatusPending
		}
		if _, err := fmt.Fprintf(w, "[%s] %-10s %s\n", statusMarks[s], p, s); err != nil {
			return fmt.Errorf("render %s: %w", p, err)
		}
	}
	return nil
}
