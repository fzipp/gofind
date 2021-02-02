# gofind

[![PkgGoDev](https://pkg.go.dev/badge/github.com/fzipp/gofind)](https://pkg.go.dev/github.com/fzipp/gofind)
![Build Status](https://github.com/fzipp/gofind/workflows/build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fzipp/gofind)](https://goreportcard.com/report/github.com/fzipp/gofind)

Gofind searches for Go modules via [pkg.go.dev](https://pkg.go.dev).

## Installation

```
go get github.com/fzipp/gofind
```

## Usage

```
gofind [-a] [-raw] query ...

Flags:
    -a     load all search results if set, not just the first 10 results
    -raw   don't apply any formatting if set
```

## Examples

Search for packages providing logging functionality:

```
$ gofind logging
log
    Package log implements a simple logging package.

    Version: go1.15.6 | Published: Dec  3, 2020 | Imported by: 203356 | License: BSD-3-Clause

github.com/sirupsen/logrus
    Package logrus is a structured logger for Go, completely API compatible
    with the standard library logger.

    Version: v1.7.0 | Published: May 28, 2020 | Imported by: 46315 | License: MIT

github.com/ethereum/go-ethereum/log
    Package log15 provides an opinionated, simple toolkit for best-practice
    logging that is both human and machine readable.

    Version: v1.9.25 | Published: Dec 11, 2020 | Imported by: 8625 | Licenses: Apache-2.0, GPL-3.0
...
```

Search for multiple terms:

```
$ gofind go cloud
```

Search for an exact match:

```
$ gofind "go cloud"
```

Combine searches:

```
$ gofind yaml OR json
```

### Tip

In case you want automatic paging if the output doesn't fit on one screen,
you can add the following function to your shell profile
(e.g. ~/.bash_profile) on Unix or Linux systems:

```
# Automatically page gofind output if it doesn't fit on one screen.
gofind() {
  command gofind "$@" | less -X -F
}
```

## License

This project is free and open source software licensed under the
[BSD 3-Clause License](LICENSE).
