// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

// print command prints text documents to selected printer.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chenxi2015/winprinters"
)

var (
	copies    = flag.Int("n", 1, "number of copies to print")
	printerId = flag.String("p", findDefaultPrinter(), "printer name or printer index from printer list")
	doList    = flag.Bool("l", false, "list printers")
)

func findDefaultPrinter() string {
	p, err := printers.GetDefault()
	if err != nil {
		return ""
	}
	return p
}

func listPrinters() error {
	printerNames, err := printers.ReadNames()
	if err != nil {
		return err
	}
	defaultPrinter, err := printers.GetDefault()
	if err != nil {
		return err
	}
	for i, p := range printerNames {
		s := " "
		if p == defaultPrinter {
			s = "*"
		}
		fmt.Printf(" %s %d. %s\n", s, i, p)
	}
	return nil
}

func selectPrinter() (string, error) {
	n, err := strconv.Atoi(*printerId)
	if err != nil {
		// must be a printer name
		return *printerId, nil
	}
	printerNames, err := printers.ReadNames()
	if err != nil {
		return "", err
	}
	if n < 0 {
		return "", fmt.Errorf("printer index (%d) cannot be negative", n)
	}
	if n >= len(printerNames) {
		return "", fmt.Errorf("printer index (%d) is too large, there are only %d printers", n, len(printerNames))
	}
	return printerNames[n], nil
}

func printOneDocument(printerName, documentName string, lines []string) error {
	p, err := printers.Open(printerName)
	if err != nil {
		return err
	}
	defer func(p *printers.Printer) {
		_ = p.Close()
	}(p)

	err = p.StartRawDocument(documentName)
	if err != nil {
		return err
	}
	defer func(p *printers.Printer) {
		_ = p.EndDocument()
	}(p)

	err = p.StartPage()
	if err != nil {
		return err
	}

	for _, line := range lines {
		_, _ = fmt.Fprintf(p, "%s\r\n", line)
	}

	return p.EndPage()
}

func printDocument(path string) error {
	if *copies < 0 {
		return fmt.Errorf("number of copies to print (%d) cannot be negative", *copies)
	}

	output, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(output), "\n")

	printerName, err := selectPrinter()
	if err != nil {
		return err
	}

	for i := 0; i < *copies; i++ {
		err := printOneDocument(printerName, path, lines)
		if err != nil {
			return err
		}
	}
	return nil
}

func usage() {
	_, _ = fmt.Fprintln(os.Stderr)
	_, _ = fmt.Fprintf(os.Stderr, "usage: print [-n=<copies>] [-p=<printer>] <file-path-to-print>\n")
	_, _ = fmt.Fprintf(os.Stderr, "       or\n")
	_, _ = fmt.Fprintf(os.Stderr, "       print -l\n")
	_, _ = fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
	os.Exit(1)
}

func exit(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if *doList {
		exit(listPrinters())
	}
	switch len(flag.Args()) {
	case 0:
		_, _ = fmt.Fprintf(os.Stderr, "no document path to print provided\n")
	case 1:
		exit(printDocument(flag.Arg(0)))
	default:
		_, _ = fmt.Fprintf(os.Stderr, "too many parameters provided\n")
	}
	usage()
}
