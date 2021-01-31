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
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, usageDoc)
	os.Exit(2)
}

const usageDoc = `Find Go packages via pkg.go.dev.
usage:
        gofind [<flag> ...] <query> ...

Flags
        -raw   don't apply any formatting if set
`

func main() {
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

	run(query, *rawFlag)
}

func run(query string, raw bool) {
	modules, err := search(query)
	check(err)
	for _, mod := range modules {
		if raw {
			fmt.Println(mod.moduleName + "\t" + mod.synopsis + "\t" + mod.info)
			continue
		}
		check(mod.writeTo(os.Stdout))
	}
}

func search(query string) ([]searchResult, error) {
	v := url.Values{}
	v.Set("q", query)
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
		return nil, fmt.Errorf("HTTP status code error: %d %s", res.StatusCode, res.Status)
	}

	htmlDoc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var results []searchResult
	htmlDoc.Find(".SearchSnippet").Each(func(i int, s *goquery.Selection) {
		moduleName := strings.TrimSpace(s.Find(".SearchSnippet-header").Text())
		synopsis := strings.TrimSpace(s.Find(".SearchSnippet-synopsis").Text())
		info := formatInfo(s.Find(".SearchSnippet-infoLabel").Text())
		results = append(results, searchResult{
			moduleName: moduleName,
			synopsis:   synopsis,
			info:       info,
		})
	})
	return results, nil
}

func formatInfo(info string) string {
	var parts []string
	for _, p := range strings.Split(info, "|") {
		s := strings.SplitN(p, ":", 2)
		if len(s) > 1 {
			label := strings.TrimSpace(s[0])
			value := strings.TrimSpace(s[1])
			parts = append(parts, label+": "+value)
		}
	}
	return strings.Join(parts, " | ")
}

type searchResult struct {
	moduleName string
	synopsis   string
	info       string
}

const (
	punchCardWidth = 80
	indent         = "    "
)

func (s searchResult) writeTo(w io.Writer) error {
	_, err := fmt.Fprintln(w, s.moduleName)
	if err != nil {
		return err
	}
	if s.synopsis != "" {
		doc.ToText(w, s.synopsis, indent, "", punchCardWidth-2*len(indent))
	}
	_, err = fmt.Fprintln(w)
	_, err = fmt.Fprintln(w, indent+s.info)
	_, err = fmt.Fprintln(w)
	return err
}

func check(err error) {
	if err != nil {
		fail(err)
	}
}

func fail(msg ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
