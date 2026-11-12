// Command fix-factory is a thin CLI over the fix-factory packages.
package main

import (
	"fmt"
	"os"

	"github.com/pedromctech/fix-factory/internal/monitor"
	"github.com/pedromctech/fix-factory/internal/repo"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 && args[0] == "status" {
		if err := monitor.Render(os.Stdout, monitor.NewTracker().Snapshot()); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if len(args) != 2 || args[0] != "add" {
		fmt.Fprintln(os.Stderr, "usage: fix-factory add <owner/name> | status")
		os.Exit(2)
	}
	r, err := repo.Parse(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("added %s\n", r)
}
