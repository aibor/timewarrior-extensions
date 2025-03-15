// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext

import (
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"slices"
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

// SplitIntoDays splits the entry in multiple entries, one per day, at the
// given split time.
//
// It splits the [Entry] into days at the given clock [time.Time] of each day.
// The date parts of that given split [time.Time] are ignored. To split at
// midnight, just pass an empty value time.Time{}.
func SplitIntoDays(entry Entry, splitClock time.Time) EntryIterator {
	return func(yield func(Entry) bool) {
		splitTime := setClock(entry.Start.Time, splitClock)
		if entry.Start.Compare(splitTime) >= 0 {
			splitTime = splitTime.AddDate(0, 0, 1)
		}

		for entry.CurrentEnd().Compare(splitTime) > 0 {
			before := entry
			before.End = Time{splitTime}
			splitTime = splitTime.AddDate(0, 0, 1)
			entry.Start = before.End

			if !yield(before) {
				return
			}
		}

		if !yield(entry) {
			return
		}
	}
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

// All returns an [EntryIterator] over all entries.
func (e Entries) All() EntryIterator {
	return slices.Values(e)
}

// EntryIterator is a single value iterator for [Entry]s.
type EntryIterator = iter.Seq[Entry]

// EntryFilter is a function that returns true for [Entry]s that pass the
// filter and false for [Entry]s that should be ignored.
type EntryFilter func(Entry) bool

// Filter returns an [EntryIterator] for all [Entry]s that pass the
// [EntryFilter].
func (f EntryFilter) Filter(entries EntryIterator) EntryIterator {
	return func(yield func(Entry) bool) {
		for entry := range entries {
			if !f(entry) {
				continue
			}

			if !yield(entry) {
				return
			}
		}
	}
}

// SplitAtMidnight creates a list of single-day entries.
//
// The [Entry.ID] will be copied when an [Entry] is split, causing multiple
// entries with the same ID to exist.
// The covered sum of durations does not change.
func SplitAtMidnight(entries EntryIterator) EntryIterator {
	return func(yield func(Entry) bool) {
		for entry := range entries {
			for e := range SplitIntoDays(entry, time.Time{}) {
				if !yield(e) {
					return
				}
			}
		}
	}
}
