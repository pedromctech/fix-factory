// Package ui is the interactive terminal interface for adding repos.
package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/pedromctech/fix-factory/internal/repo"
)

// Adder receives the repos the user adds.
type Adder interface {
	Add(repo.Ref) error
}

// Memory is an in-memory Adder that ignores duplicates.
type Memory struct {
	refs []repo.Ref
}

// NewMemory returns an empty Memory.
func NewMemory() *Memory {
	return &Memory{}
}

// Add stores r unless it is already there.
func (m *Memory) Add(r repo.Ref) error {
	for _, e := range m.refs {
		if e == r {
			return nil
		}
	}
	m.refs = append(m.refs, r)
	return nil
}

// List returns the added repos in insertion order.
func (m *Memory) List() []repo.Ref {
	return append([]repo.Ref(nil), m.refs...)
}

// Run reads one repo per line from in and adds each valid one. Invalid lines
// print an error and the loop continues. "quit", "exit" or EOF end it.
func Run(in io.Reader, out io.Writer, adder Adder) error {
	sc := bufio.NewScanner(in)
	if _, err := fmt.Fprintln(out, "Enter a repo (owner/name or GitHub URL), or quit:"); err != nil {
		return fmt.Errorf("write prompt: %w", err)
	}
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch line {
		case "":
			continue
		case "quit", "exit":
			return nil
		}
		ref, err := repo.Parse(line)
		if err != nil {
			if _, werr := fmt.Fprintf(out, "error: %v\n", err); werr != nil {
				return fmt.Errorf("write output: %w", werr)
			}
			continue
		}
		if err := adder.Add(ref); err != nil {
			return fmt.Errorf("add %s: %w", ref, err)
		}
		if _, err := fmt.Fprintf(out, "added %s\n", ref); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	return nil
}
