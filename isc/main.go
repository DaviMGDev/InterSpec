// isc — InterSpec Compiler
//
// Usage:
//
//	isc strip [flags] [file]   Strip comments from an .is file
package main

import (
	"fmt"
	"os"

	"github.com/davi/isc/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: isc <command> [flags] [args]\n")
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		fmt.Fprintf(os.Stderr, "  strip   Strip comments from an InterSpec (.is) file\n")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "strip":
		if err := cmd.RunStrip(args); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n", command)
		fmt.Fprintf(os.Stderr, "Available commands: strip\n")
		os.Exit(1)
	}
}
