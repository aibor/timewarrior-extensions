// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: MIT

package twext

import "testing"

// MustParseTime parses a time from a string in the timewarrior format.
func MustParseTime(tb testing.TB, s string) Time {
	tb.Helper()

	t, err := ParseTime(s)
	if err != nil {
		tb.Fatalf("must parse time: %v", err)
	}

	return t
}
