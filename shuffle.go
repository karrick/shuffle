package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

func main() {
	var err error
	if len(os.Args) > 1 {
		err = file(os.Args[1])
	} else {
		err = standard()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
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
