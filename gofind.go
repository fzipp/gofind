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

var raw = flag.Bool("raw", false, "don't apply any formatting")

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

	query := strings.Join(flag.Args(), " ")
	modules, err := search(query)
	check(err)
	for _, mod := range modules {
		err = mod.writeTo(os.Stdout)
		check(err)
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
		if *raw {
			fmt.Println(moduleName + "\t" + synopsis)
			return
		}
		results = append(results, searchResult{
			moduleName: moduleName,
			synopsis:   synopsis,
		})
	})
	return results, nil
}

type searchResult struct {
	moduleName string
	synopsis   string
}

const (
	punchCardWidth = 80
	indent         = "    "
)

func (r searchResult) writeTo(w io.Writer) error {
	_, err := fmt.Fprintln(w, r.moduleName)
	if err != nil {
		return err
	}
	if r.synopsis != "" {
		doc.ToText(w, r.synopsis, indent, "", punchCardWidth-2*len(indent))
	}
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
