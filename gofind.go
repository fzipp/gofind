// Copyright 2013 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gofind searches for Go packages via pkg.go.dev.
//
// Usage:
//      gofind [<flag> ...] <query> ...
//
// Flags
//      -raw   don't apply any formatting if set
package main

import (
	"flag"
	"fmt"
	"go/doc"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func usage() {
	_, _ = fmt.Fprint(os.Stderr, usageDoc)
	os.Exit(2)
}

const usageDoc = `Find Go modules via pkg.go.dev.

Usage:
    gofind [-a] [-raw] query ...

Flags:
    -a     load all search results if set, not just the first 10 results
    -raw   don't apply any formatting if set

Examples:
    gofind logging
    gofind -a logging
    gofind go cloud        # Search for multiple terms
    gofind "go cloud"      # Search for an exact match
    gofind yaml OR json    # Combine searches
`

func main() {
	allFlag := flag.Bool("a", false, "load all search results, not just the first 10 results")
	rawFlag := flag.Bool("raw", false, "don't apply any formatting")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}
	args := flag.Args()
	for i, arg := range args {
		if strings.Contains(arg, " ") {
			args[i] = `"` + arg + `"`
		}
	}
	query := strings.Join(args, " ")

	run(query, *allFlag, *rawFlag)
}

func run(query string, all, raw bool) {
	modules, err := search(query, all)
	check(err)
	for _, mod := range modules {
		if raw {
			fmt.Println(mod.modulePath + "\t" + mod.synopsis + "\t" + mod.info)
			continue
		}
		check(mod.writeTo(os.Stdout))
	}
}

func search(query string, all bool) ([]searchResult, error) {
	limit := 10
	if all {
		limit = 100
	}
	return searchLimited(query, limit)
}

func searchLimited(query string, limit int) ([]searchResult, error) {
	v := url.Values{}
	v.Set("q", query)
	v.Set("limit", strconv.Itoa(limit))
	req, err := http.NewRequest("GET", "https://pkg.go.dev/search?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "text/html")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status code error: %d %sr", res.StatusCode, res.Status)
	}

	htmlDoc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return scrapeSearchResults(htmlDoc)
}

func scrapeSearchResults(htmlDoc *goquery.Document) (sr []searchResult, err error) {
	htmlDoc.Find(".SearchSnippet").Each(func(i int, s *goquery.Selection) {
		header := s.Find(".SearchSnippet-headerContainer").Text()
		moduleName := strings.TrimSpace(header[:strings.Index(header, "(")])
		modulePath := trimParens(strings.TrimSpace(s.Find(".SearchSnippet-header-path").Text()))
		synopsis := strings.TrimSpace(s.Find(".SearchSnippet-synopsis").Text())
		info := formatInfo(s.Find(".SearchSnippet-infoLabel").Text())
		sr = append(sr, searchResult{
			moduleName: moduleName,
			modulePath: modulePath,
			synopsis:   synopsis,
			info:       info,
		})
	})
	return sr, err
}

func trimParens(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "("), ")")
}

func formatInfo(info string) string {
	var parts []string
	for _, p := range strings.Split(info, "|") {
		parts = append(parts, strings.TrimSpace(p))
	}
	return strings.Join(parts, " | ")
}

type searchResult struct {
	moduleName string
	modulePath string
	synopsis   string
	info       string
}

const (
	punchCardWidth = 80
	indent         = "    "
)

func (s searchResult) writeTo(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s (%s)\n", s.moduleName, s.modulePath)
	if err != nil {
		return err
	}
	if s.synopsis != "" {
		doc.ToText(w, s.synopsis, indent, "", punchCardWidth-2*len(indent))
	}
	_, err = fmt.Fprintf(w, "\n%s%s\n\n", indent, s.info)
	return err
}

func check(err error) {
	if err != nil {
		fail(err)
	}
}

func fail(msg ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}
