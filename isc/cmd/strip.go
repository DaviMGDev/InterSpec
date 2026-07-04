// Command strip removes comments from InterSpec (.is) files.
package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/davi/isc/stripper"
)

// RunStrip is the entry point for the "strip" command.
func RunStrip(args []string) error {
	fs := flag.NewFlagSet("strip", flag.ContinueOnError)
	output := fs.String("o", "", "output file (default: stdout)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: isc strip [flags] [file]\n")
		fmt.Fprintf(os.Stderr, "\nStrip all comments from an InterSpec (.is) file.\n")
		fmt.Fprintf(os.Stderr, "If no file is given, reads from stdin.\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Determine input source
	var r io.Reader
	var fileName string

	if fs.NArg() == 0 {
		r = os.Stdin
		fileName = "stdin"
	} else {
		fileName = fs.Arg(0)
		f, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("opening input file: %w", err)
		}
		defer f.Close()
		r = f
	}

	// Determine output destination
	var w io.Writer
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	} else {
		w = os.Stdout
	}

	s := stripper.New(fileName)
	if err := s.Strip(r, w); err != nil {
		return fmt.Errorf("stripping comments: %w", err)
	}

	return nil
}
