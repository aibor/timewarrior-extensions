// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: MIT

package twext

import (
	"fmt"
	"time"
)

// DateFmt is the date format used by timewarrior.
const DateFmt = "20060102T150405Z07"

// Time extends [time.Time] with a custom functionality.
type Time struct {
	time.Time
}

// ParseTime parses a time from a string according to the timewarrior format.
func ParseTime(s string) (Time, error) {
	t, err := time.Parse(DateFmt, s)
	if err != nil {
		return Time{}, fmt.Errorf("parse time: %w", err)
	}

	return Time{t}, nil
}

// UnmarshalJSON unmarshals timestamp strings from the timewarrior format.
func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) == 0 ||
		string(data) == `""` ||
		string(data) == "null" {
		return nil
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return ErrDateUnmarshalNotString
	}

	data = data[len(`"`) : len(data)-len(`"`)]

	parsed, err := ParseTime(string(data))
	if err != nil {
		return err
	}

	*t = parsed

	return nil
}

// SameDate compares two [Time] values and returns true if they are the same
// date.
func (t *Time) SameDate(o *Time) bool {
	yt, mt, dt := t.Date()
	yo, mo, do := o.Date()

	return yt == yo && mt == mo && dt == do
}
