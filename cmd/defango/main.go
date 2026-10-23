/*
defango - URL / IP / Email defanging with Golang. Make IoC harmless.

This repository is under MIT License https://github.com/edoardottt/defango/blob/main/LICENSE
*/

// Command defango defangs the IoCs in its standard input, or refangs them.
//
// Usage:
//
//	defango [-r] < input
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/edoardottt/defango"
)

func main() {
	refang := flag.Bool("r", false, "refang instead of defanging")

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: defango [-r] < input")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(2) //nolint:mnd // usage error, as in the flag package
	}

	if err := run(os.Stdin, os.Stdout, *refang); err != nil {
		fmt.Fprintln(os.Stderr, "defango:", err)
		os.Exit(1)
	}
}

// run defangs, or refangs, r into w line by line.
func run(r io.Reader, w io.Writer, refang bool) error {
	in := bufio.NewReader(r)
	out := bufio.NewWriter(w)

	for {
		line, readErr := in.ReadString('\n')

		if refang {
			line = defango.Refang(line)
		} else {
			line = defango.Text(line)
		}

		if _, err := out.WriteString(line); err != nil {
			return err
		}

		if errors.Is(readErr, io.EOF) {
			return out.Flush()
		}

		if readErr != nil {
			return readErr
		}
	}
}
