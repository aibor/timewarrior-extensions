// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Entry is a timewarrior entry that covers a single recorded time interval.
type Entry struct {
	ID    int      `json:"id"`
	Start Time     `json:"start"`
	End   Time     `json:"end,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

// Duration calculates the duration of the [Entry].
//
// If the [Entry] is still active, the current time is used as the end time.
func (e *Entry) Duration() time.Duration {
	end := e.End
	if end.IsZero() {
		end.Time = time.Now()
	}

	if end.Before(e.Start.Time) {
		return 0
	}

	return end.Sub(e.Start.Time)
}

// Entries is a list of [Entry]s.
type Entries []Entry

// readEntries parses the given reader with JSON data in timewarrior format
// into a list of [Entry]s.
func readEntries(reader io.Reader) (Entries, error) {
	var entries Entries

	jsonReader := json.NewDecoder(reader)

	err := jsonReader.Decode(&entries)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("decode json: %w", err)
	}

	return entries, nil
}
