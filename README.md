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
default is 8 hours. If you already have overtime (or undertime), you can have
flextime take it into account by setting a global offset.

#### Configuration

The following configuration keys are supported:

| Key                               | Type     | Default                           | Description                                   |
|-----------------------------------|----------|-----------------------------------|-----------------------------------------------|
| `flextime.time_per_day`           | Duration | `8h`                              | Default daily time target.                    |
| `flextime.time_per_day.<weekday>` | Duration | value of `flextime.time_per_day`  | Weekday specific time target.                 |
| `flextime.offset_total`           | Duration | `0`                               | Time spent or lacking from a previous period. |
| `flextime.aggregation_strategy`   | Enum     | `single-day-only`                 | Strategy to use for aggregating the entries.  |
| `verbose`                         | Bool     | true                              | Print daily sums.                             |
| `debug`                           | Bool     | false                             | Enable debug output.                          |

Durations must be given in a format supported by
[go's time duration parser][go-time-duration].

| Aggregation Strategy | Description                                                                |
|----------------------|----------------------------------------------------------------------------|
| `single-day-only`    | Discard entries spanning multiple days.                                    |
| `into-start-date`    | Count entries spanning multiple days for the day that the entry starts on. |
| `into-end-date`      | Count entries spanning multiple days for the day that the entry end on.    |
| `split-at-midnight`  | Split entries spanning over multiple days at midnight.                     |

##### Example

This is an example configuration for a 35-hour week
in which Fridays are a half workday and weekends off.
Additionally, 3 hours and 23 minutes of pre-existing overtime are taken into account,
and work is counted for the day it happened on.
```
flextime.time_per_day 7h30m
flextime.time_per_day.friday 5h
flextime.time_per_day.saturday 0h
flextime.time_per_day.sunday 0h
flextime.offset_total 3h23m
flextime.aggregation_strategy split-at-midnight
```

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
