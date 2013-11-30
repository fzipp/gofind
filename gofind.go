// Copyright 2013 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gofind searches for Go packages via godoc.org.
//
// Usage:
//         gofind [<flag> ...] <query> ...
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/doc"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const help = `Find Go packages via godoc.org.
Usage: gofind [<flag> ...] <query> ...`

var raw = flag.Bool("raw", false, "don't apply any formatting if true")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, help)
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}
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
		line := scanner.Text()
		if *raw {
			fmt.Println(line)
			continue
		}
		PackageFrom(line).WriteTo(os.Stdout)
	}
	if err := scanner.Err(); err != nil {
		exitError(err)
	}
}

type Package struct {
	Path, Synopsis string
}

func PackageFrom(line string) (pkg Package) {
	s := strings.SplitAfterN(line, " ", 2)
	if len(s) > 0 {
		pkg.Path = s[0]
	}
	if len(s) > 1 {
		pkg.Synopsis = s[1]
	}
	return
}

const (
	punchCardWidth = 80
	indent         = "    "
)

func (pkg Package) WriteTo(w io.Writer) {
	fmt.Fprintln(w, pkg.Path)
	if pkg.Synopsis != "" {
		doc.ToText(w, pkg.Synopsis, indent, "", punchCardWidth-2*len(indent))
	}
	fmt.Println()
}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
