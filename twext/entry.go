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
	End   Time     `json:"end"`
	Tags  []string `json:"tags,omitempty"`
}

// Duration calculates the duration of the [Entry].
//
// If the [Entry] is still active, the current time is used as the end time.
func (e *Entry) Duration() time.Duration {
	end := e.CurrentEnd()

	if end.Before(e.Start.Time) {
		return 0
	}

	return end.Sub(e.Start.Time)
}

// CurrentEnd calculates the actual end of the [Entry].
//
// If the [Entry] is still active, the current time is returned.
// Otherwise, the recorded time is returned.
func (e *Entry) CurrentEnd() *Time {
	if e.IsActive() {
		return &Time{time.Now()}
	}

	return &e.End
}

// IsActive checks if this [Entry] is still active.
func (e *Entry) IsActive() bool {
	return e.End.IsZero()
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

// Filter removes entries that do not match the filter.
func (e Entries) Filter(filter func(entry Entry) bool) Entries {
	entries := make(Entries, 0, len(e))

	for _, entry := range e {
		if filter(entry) {
			entries = append(entries, entry)
		}
	}

	return entries
}

// SplitAtMidnight creates a list of single-day entries.
//
// The [Entry.ID] will be copied when an [Entry] is split, causing multiple
// entries with the same ID to exist.
// The covered sum of durations does not change.
func (e Entries) SplitAtMidnight() Entries {
	entries := make(Entries, 0, len(e))
	for _, entry := range e {
		entries = appendAsSingleDayEntries(entries, entry)
	}

	return entries
}

func appendAsSingleDayEntries(entries Entries, entry Entry) Entries {
	end := entry.CurrentEnd()
	for !entry.Start.SameDate(end) {
		midnight := Time{Time: time.Date(
			entry.Start.Year(),
			entry.Start.Month(),
			entry.Start.Day(),
			0, 0, 0, 0,
			entry.Start.Location(),
		).AddDate(0, 0, 1)}

		entries = append(entries, Entry{
			ID:    entry.ID,
			Start: entry.Start,
			End:   midnight,
			Tags:  entry.Tags,
		})

		entry = Entry{
			ID:    entry.ID,
			Start: midnight,
			End:   entry.End,
			Tags:  entry.Tags,
		}
	}

	entries = append(entries, entry)

	return entries
}
