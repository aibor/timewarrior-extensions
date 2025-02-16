<!--
SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>

SPDX-License-Identifier: GPL-3.0-or-later
-->

# timewarrior-extensions

[![PkgGoDev][pkg-go-dev-badge]][pkg-go-dev]
[![Go Report Card][go-report-card-badge]][go-report-card]
[![Actions][actions-test-badge]][actions-test]

Collections of timewarrior-extensions and a golang library for parsing
input from timewarrior extension API.

## Extensions

### flextime

Sums time spent per day and calculates the difference to the daily target. The
default is 8 hours. It can be changed by setting the config variable
`flextime.time_per_day`. It takes a [go time duration string][go-time-duration].

If you already have overtime, you can have flextime take it into account by
setting the config variable `flextime.offset_total`. It takes a
[go time duration string][go-time-duration].

#### Install

After cloning the repo, build directly into your timewarrior extensions
directory:

```
go build -o ~/.timewarrior/extensions/flextime ./cmd/flextime
```

Then run timew with `flextime` report:

```
timew flextime
```

The output looks like this. `actual` is the actual accounted time. `target` is
the daily/total target. `diff` is the difference between `actual` and `target`.

```
          date    actual    target      diff
    2024-06-30      3:39      8:00     -4:20
    2024-07-08      0:06      8:00     -7:53
         total      3:45     16:00    -12:14
```

## Golang library

The package [twext][pkg-go-dev] implements basic functions for reading input
from timewarrior as described in 
[their docs](https://timewarrior.net/docs/api/). All of the extensions in this
repository are built based on the library.

[pkg-go-dev]:           https://pkg.go.dev/github.com/aibor/timewarrior-extensions/twext
[pkg-go-dev-badge]:     https://pkg.go.dev/badge/github.com/aibor/timewarrior-extensions/twext
[go-report-card]:       https://goreportcard.com/report/github.com/aibor/timewarrior-extensions
[go-report-card-badge]: https://goreportcard.com/badge/github.com/aibor/timewarrior-extensions
[actions-test]:         https://github.com/aibor/timewarrior-extensions/actions/workflows/test.yaml
[actions-test-badge]:   https://github.com/aibor/timewarrior-extensions/actions/workflows/test.yaml/badge.svg?branch=main
[go-time-duration]:     https://pkg.go.dev/time#ParseDuration
