// Package cli holds the command dispatch of the fix-factory binary.
package cli

import (
	"fmt"
	"io"

	"github.com/pedromctech/fix-factory/internal/monitor"
	"github.com/pedromctech/fix-factory/internal/repo"
	"github.com/pedromctech/fix-factory/internal/ui"
)

const usage = "usage: fix-factory [add <owner/name> | status | help]  (no args: interactive mode)\n"

// Run executes the CLI with args (without the program name) and returns the
// exit code: 0 ok, 1 error, 2 usage.
func Run(args []string, in io.Reader, out, errOut io.Writer) int {
	if len(args) == 0 {
		if err := ui.Run(in, out, ui.NewMemory()); err != nil {
			return fail(errOut, 1, fmt.Sprintf("error: %v\n", err))
		}
		return 0
	}
	switch args[0] {
	case "help", "-h", "--help":
		if len(args) != 1 {
			return fail(errOut, 2, usage)
		}
		if _, err := io.WriteString(out, usage); err != nil {
			return 1
		}
		return 0
	case "status":
		if len(args) != 1 {
			return fail(errOut, 2, usage)
		}
		if err := monitor.Render(out, monitor.NewTracker().Snapshot()); err != nil {
			return fail(errOut, 1, fmt.Sprintf("error: %v\n", err))
		}
		return 0
	case "add":
		if len(args) != 2 {
			return fail(errOut, 2, usage)
		}
		r, err := repo.Parse(args[1])
		if err != nil {
			return fail(errOut, 1, fmt.Sprintf("error: %v\n", err))
		}
		if _, err := fmt.Fprintf(out, "added %s\n", r); err != nil {
			return 1
		}
		return 0
	}
	return fail(errOut, 2, usage)
}

// fail writes msg to errOut and returns code. A failed write cannot be
// reported anywhere else, so the exit code is all that is left.
func fail(errOut io.Writer, code int, msg string) int {
	if _, err := io.WriteString(errOut, msg); err != nil {
		return code
	}
	return code
}
