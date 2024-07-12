<!--
SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>

SPDX-License-Identifier: MIT
-->

# timewarrior-extensions

Collections of timewarrior-extensions and a golang library for parsing
input from timewarrior extension API.

## Extensions

### flextime

Sums time spent per day and calculates the difference to the daily goal. The
default is 8 hours. It can be changed by setting the config variable
`flextime.hours_per_day`.

#### Install

Clone the repo. Then build directly into your timewarrior extensions directory:

```
$ git clone github.com/aibor/timewarrior-extensions/cmd/flextime
$ cd flextime
$ go build -o ~/.timewarrior/extensions/flextime ./cmd/flextime
```

## Golang library

The package `twext`  implements basic functions for reading input from
timewarrior as [in their docs](https://timewarrior.net/docs/api/).
All of the extensions in this repository are built based on the library.
```
