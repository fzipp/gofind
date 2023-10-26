# gofind

[![PkgGoDev](https://pkg.go.dev/badge/github.com/fzipp/gofind)](https://pkg.go.dev/github.com/fzipp/gofind)
![Build Status](https://github.com/fzipp/gofind/workflows/build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fzipp/gofind)](https://goreportcard.com/report/github.com/fzipp/gofind)

Gofind conveniently searches for Go modules
from the command line
and lists them there
without the need to visit [pkg.go.dev](https://pkg.go.dev)
through a web browser.

## Installation

```
go install github.com/fzipp/gofind@latest
```

## Usage

```
gofind [-a] query ...

Flags:
    -a     load all search results if set, not just the first 10 results
```

## Examples

Search for packages providing logging functionality:

```
$ gofind logging
log (log)
    Package log implements a simple logging package.

    Imported by 369,051 | go1.17.3 published on 5 days ago | BSD-3-Clause

logrus (github.com/sirupsen/logrus)
    Package logrus is a structured logger for Go, completely API compatible
    with the standard library logger.

    Imported by 75,868 | v1.8.1 published on Mar  9, 2021 | MIT

log (github.com/go-kit/kit/log)
    Package log provides a structured logger.

    Imported by 5,625 | v0.12.0 published on Sep 18, 2021 | MIT
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

If you wish to enable automatic paging
when the output doesn't fit on one screen,
you can add the following function to your shell profile
(e.g. ~/.bash_profile) 
on Unix or Linux systems:

```
# Automatically page gofind output if it doesn't fit on one screen.
gofind() {
  command gofind "$@" | less -X -F
}
```

## License

This project is free and open source software licensed under the
[BSD 3-Clause License](LICENSE).
