// Command fix-factory is a thin CLI over the fix-factory packages.
package main

import (
	"os"

	"github.com/pedromctech/fix-factory/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
