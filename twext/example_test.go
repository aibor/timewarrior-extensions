// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext_test

import (
	"fmt"
	"strings"

	"github.com/aibor/timewarrior-extensions/twext"
)

func Example() {
	stdin := strings.NewReader(`color: on
reports.day.axis: internal

[
{"id":3,"start":"20240629T102128Z","end":"20240629T102131Z"},
{"id":2,"start":"20240630T143940Z","end":"20240630T143943Z"},
{"id":1,"start":"20240630T144010Z"}
]
`)

	reader := twext.NewReader(stdin)

	cfg, err := reader.ReadConfig()
	if err != nil {
		panic("cannot read config section: " + err.Error())
	}

	fmt.Println("color:", cfg["color"])

	entries, err := reader.ReadEntries()
	if err != nil {
		panic("cannot read entries: " + err.Error())
	}

	groups := twext.Group(
		entries.All(),
		func(e twext.Entry) int {
			return e.Start.Day()
		},
		func(r int, _ twext.Entry) int {
			return r + 1
		},
	)

	fmt.Println("29:", groups[29])
	fmt.Println("30:", groups[30])

	// Output:
	// color: on
	// 29: 1
	// 30: 2
}
