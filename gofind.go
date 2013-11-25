// Copyright 2013 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gofind searches for Go packages via godoc.org.
//
// Usage:
//         gofind <query> ...
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const help = `Find Go packages via godoc.org.
Usage: gofind <query> ...`

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, help)
	}
	flag.Parse()

	query := strings.Join(flag.Args(), " ")

	v := url.Values{}
	v.Set("q", query)
	req, err := http.NewRequest("GET", "http://godoc.org/?"+v.Encode(), nil)
	if err != nil {
		exitError(err)
	}
	req.Header.Add("Accept", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		exitError(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		printPackage(parseLine(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		exitError(err)
	}
}

func parseLine(l string) (name, desc string) {
	s := strings.SplitAfterN(l, " ", 2)
	if len(s) > 0 {
		name = s[0]
	}
	if len(s) > 1 {
		desc = s[1]
	}
	return
}

func printPackage(name, desc string) {
	fmt.Println(name)
	if desc != "" {
		fmt.Println("   ", desc)
	}
}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
