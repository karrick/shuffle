package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/karrick/golf"
)

func main() {
	optHelp := golf.BoolP('h', "help", false, "Print usage information then exit")
	golf.Parse()

	if *optHelp {
		fmt.Fprintf(os.Stderr, "%s [file1 [file2 [fileN ...]]]\n", filepath.Base(os.Args[0]))
		if *optHelp {
			fmt.Fprintln(os.Stderr, "        randomizes a line delimited stream.\n")
			fmt.Fprintln(os.Stderr, "Without filename arguments this reads and shuffles standard input and writes to\nstandard output. With filename arguments, each file is read, shuffled, and\nre-written individually.\n")
			golf.Usage()
		}
		exit(nil)
	}

	if golf.NArg() == 0 {
		exit(standard())
	}

	var rerr error
	for _, pathname := range golf.Args() {
		if err := file(pathname); err != nil {
			if rerr == nil {
				rerr = err
			}
			fmt.Fprintf(os.Stderr, "WARNING: shuffle cannot process: %q: %s", pathname, err)
		}
	}
	exit(rerr)
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func file(pathname string) error {
	fr, err := os.Open(pathname)
	if err != nil {
		return err
	}
	defer fr.Close()

	lines, err := readLines(fr)
	if err != nil {
		return err
	}

	fw, err := os.Create(pathname)
	if err != nil {
		return err
	}
	defer fw.Close()

	return writeLines(fw, shuffle(lines))
}

func standard() error {
	lines, err := readLines(os.Stdin)
	if err != nil {
		return err
	}
	return writeLines(os.Stdout, shuffle(lines))
}

func readLines(ior io.Reader) ([]string, error) {
	var lines []string

	scanner := bufio.NewScanner(ior)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func shuffle(src []string) []string {
	rand.Seed(time.Now().UnixNano())
	dest := make([]string, len(src))

	for i, v := range rand.Perm(len(src)) {
		dest[i] = src[v]
	}

	return dest
}

func writeLines(iow io.Writer, lines []string) error {
	for _, line := range lines {
		_, err := fmt.Fprintf(iow, "%s\n", line)
		if err != nil {
			return err
		}
	}
	return nil
}
