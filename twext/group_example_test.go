// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext_test

import (
	"fmt"
	"time"

	"github.com/aibor/timewarrior-extensions/twext"
)

func parseTime(s string) twext.Time {
	t, err := twext.ParseTime(s)
	if err != nil {
		panic(err)
	}

	return t
}

func ExampleGroup_slices() {
	entries := twext.Entries{
		twext.Entry{
			ID:    3,
			Start: parseTime("20100629T080000Z"),
			End:   parseTime("20100629T140000Z"),
		},
		twext.Entry{
			ID:    2,
			Start: parseTime("20100629T150000Z"),
			End:   parseTime("20100629T180000Z"),
		},
		twext.Entry{
			ID:    1,
			Start: parseTime("20100630T080000Z"),
			End:   parseTime("20100630T160000Z"),
		},
	}

	groups := twext.Group(
		entries,
		func(entry twext.Entry) string {
			return entry.Start.Format(time.DateOnly)
		},
		func(result []int, entry twext.Entry) []int {
			return append(result, entry.ID)
		},
	)

	for day, group := range groups.Sorted() {
		fmt.Println(day, group)
	}
	// Output:
	// 2010-06-29 [3 2]
	// 2010-06-30 [1]
}

func ExampleGroup_reduce() {
	entries := twext.Entries{
		twext.Entry{
			ID:    3,
			Start: parseTime("20100629T080000Z"),
			End:   parseTime("20100629T140000Z"),
		},
		twext.Entry{
			ID:    2,
			Start: parseTime("20100629T150000Z"),
			End:   parseTime("20100629T180000Z"),
		},
		twext.Entry{
			ID:    1,
			Start: parseTime("20100630T080000Z"),
			End:   parseTime("20100630T160000Z"),
		},
	}

	groups := twext.Group(
		entries,
		func(entry twext.Entry) string {
			return entry.Start.Format(time.DateOnly)
		},
		func(result time.Duration, entry twext.Entry) time.Duration {
			return result + entry.Duration()
		},
	)

	for day, group := range groups.Sorted() {
		fmt.Println(day, group)
	}
	// Output:
	// 2010-06-29 9h0m0s
	// 2010-06-30 8h0m0s
}
