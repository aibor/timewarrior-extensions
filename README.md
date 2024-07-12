<!--
SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>

SPDX-License-Identifier: MIT
-->

# timewarrior-extensions

Collections of timewarrior-extensions and a golang library for parsing
input from timewarrior extension API.

## Extensions

### flextime

Sums time spent per day and calculates the difference to the daily target. The
default is 8 hours. It can be changed by setting the config variable
`flextime.hours_per_day`.

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

The package `twext`  implements basic functions for reading input from
timewarrior as [in their docs](https://timewarrior.net/docs/api/).
All of the extensions in this repository are built based on the library.
```
