// Copyright 2013 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gofind searches for Go packages via godoc.org.
//
// Usage:
//      gofind [<flag> ...] <query> ...
//
// Flags
//      -raw   don't apply any formatting if set
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

func usage() {
	fmt.Fprintf(os.Stderr, usageDoc)
	os.Exit(2)
}

const usageDoc = `Find Go packages via godoc.org.
usage:
        gofind [<flag> ...] <query> ...

Flags
        -raw   don't apply any formatting if set
`

var raw = flag.Bool("raw", false, "don't apply any formatting")

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
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
		packageFrom(line).writeTo(os.Stdout)
	}
	if err := scanner.Err(); err != nil {
		exitError(err)
	}
}

type Package struct {
	path, synopsis string
}

func packageFrom(line string) (pkg Package) {
	s := strings.SplitAfterN(line, " ", 2)
	if len(s) > 0 {
		pkg.path = s[0]
	}
	if len(s) > 1 {
		pkg.synopsis = s[1]
	}
	return
}

const (
	punchCardWidth = 80
	indent         = "    "
)

func (pkg Package) writeTo(w io.Writer) {
	fmt.Fprintln(w, pkg.path)
	if pkg.synopsis != "" {
		doc.ToText(w, pkg.synopsis, indent, "", punchCardWidth-2*len(indent))
	}
	fmt.Fprintln(w)
}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
