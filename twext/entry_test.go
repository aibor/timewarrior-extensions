// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext_test

import (
	"slices"
	"testing"
	"time"

	"github.com/aibor/timewarrior-extensions/twext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntry_Duration(t *testing.T) {
	t.Run("incomplete", func(t *testing.T) {
		expected := 1 * time.Hour
		entry := twext.Entry{
			Start: twext.Time{time.Now().Add(-expected)},
		}

		actual := entry.Duration()
		assert.Greater(t, actual, expected)
	})

	tests := []struct {
		name     string
		entry    twext.Entry
		expected string
	}{
		{
			name: "hours",
			entry: twext.Entry{
				Start: twext.MustParseTime("20100203T101530Z"),
				End:   twext.MustParseTime("20100203T142537Z"),
			},
			expected: "4h10m7s",
		},
		{
			name: "days",
			entry: twext.Entry{
				Start: twext.MustParseTime("20100203T101530Z"),
				End:   twext.MustParseTime("20100204T142537Z"),
			},
			expected: "28h10m7s",
		},
		{
			name: "before",
			entry: twext.Entry{
				Start: twext.MustParseTime("20100203T101530Z"),
				End:   twext.MustParseTime("20100202T142537Z"),
			},
			expected: "0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, err := time.ParseDuration(tt.expected)
			require.NoError(t, err)

			actual := tt.entry.Duration()
			assert.Equal(t, expected, actual)
		})
	}
}

func TestEntry_SplitIntoDays(t *testing.T) {
	tests := []struct {
		name     string
		input    twext.Entry
		split    time.Time
		expected []twext.Entry
	}{
		{
			name: "end before split time",
			input: twext.Entry{
				Start: twext.MustParseTime("20100203T101530Z"),
				End:   twext.MustParseTime("20100203T142537Z"),
			},
			split: time.Date(0, 0, 0, 15, 0, 0, 0, time.UTC),
			expected: twext.Entries{
				twext.Entry{
					Start: twext.MustParseTime("20100203T101530Z"),
					End:   twext.MustParseTime("20100203T142537Z"),
				},
			},
		},
		{
			name: "end at split time",
			input: twext.Entry{
				Start: twext.MustParseTime("20100203T101530Z"),
				End:   twext.MustParseTime("20100203T142537Z"),
			},
			split: time.Date(0, 0, 0, 14, 25, 37, 0, time.UTC),
			expected: twext.Entries{
				twext.Entry{
					Start: twext.MustParseTime("20100203T101530Z"),
					End:   twext.MustParseTime("20100203T142537Z"),
				},
			},
		},
		{
			name: "start after split time",
			input: twext.Entry{
				Start: twext.MustParseTime("20100203T101530Z"),
				End:   twext.MustParseTime("20100203T142537Z"),
			},
			split: time.Date(0, 0, 0, 15, 0, 0, 0, time.UTC),
			expected: twext.Entries{
				twext.Entry{
					Start: twext.MustParseTime("20100203T101530Z"),
					End:   twext.MustParseTime("20100203T142537Z"),
				},
			},
		},
		{
			name: "start at split time",
			input: twext.Entry{
				Start: twext.MustParseTime("20100203T101530Z"),
				End:   twext.MustParseTime("20100203T142537Z"),
			},
			split: time.Date(0, 0, 0, 10, 15, 30, 0, time.UTC),
			expected: twext.Entries{
				twext.Entry{
					Start: twext.MustParseTime("20100203T101530Z"),
					End:   twext.MustParseTime("20100203T142537Z"),
				},
			},
		},
		{
			name: "single-day split",
			input: twext.Entry{
				Start: twext.MustParseTime("20100203T101530Z"),
				End:   twext.MustParseTime("20100203T142537Z"),
			},
			split: time.Date(0, 0, 0, 13, 0, 0, 0, time.UTC),
			expected: twext.Entries{
				twext.Entry{
					Start: twext.MustParseTime("20100203T101530Z"),
					End:   twext.MustParseTime("20100203T130000Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime("20100203T130000Z"),
					End:   twext.MustParseTime("20100203T142537Z"),
				},
			},
		},
		{
			name: "multi-day entries noon",
			input: twext.Entry{
				Start: twext.MustParseTime("20100205T092755Z"),
				End:   twext.MustParseTime("20100207T163211Z"),
			},
			split: time.Date(0, 0, 0, 12, 0, 0, 0, time.UTC),
			expected: twext.Entries{
				twext.Entry{
					Start: twext.MustParseTime("20100205T092755Z"),
					End:   twext.MustParseTime("20100205T120000Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime("20100205T120000Z"),
					End:   twext.MustParseTime("20100206T120000Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime("20100206T120000Z"),
					End:   twext.MustParseTime("20100207T120000Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime("20100207T120000Z"),
					End:   twext.MustParseTime("20100207T163211Z"),
				},
			},
		},
		{
			name: "multi-day entries midnight",
			input: twext.Entry{
				Start: twext.MustParseTime("20100205T092755Z"),
				End:   twext.MustParseTime("20100207T163211Z"),
			},
			split: time.Time{},
			expected: twext.Entries{
				twext.Entry{
					Start: twext.MustParseTime("20100205T092755Z"),
					End:   twext.MustParseTime("20100206T000000Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime("20100206T000000Z"),
					End:   twext.MustParseTime("20100207T000000Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime("20100207T000000Z"),
					End:   twext.MustParseTime("20100207T163211Z"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := twext.SplitIntoDays(tt.input, tt.split)
			assert.Equal(t, tt.expected, slices.Collect(actual))
		})
	}
}

func TestEntries_Filter(t *testing.T) {
	tests := []struct {
		name            string
		entries         twext.Entries
		filter          twext.EntryFilter
		expectedEntries twext.Entries
	}{
		{
			name: "all",
			entries: twext.Entries{
				{
					Start: twext.MustParseTime("20100203T101530Z"),
					End:   twext.MustParseTime("20100203T142537Z"),
				},
				{
					Start: twext.MustParseTime("20100204T101530Z"),
					End:   twext.MustParseTime("20100204T142537Z"),
				},
			},
			filter: func(_ twext.Entry) bool {
				return true
			},
			expectedEntries: twext.Entries{
				{
					Start: twext.MustParseTime("20100203T101530Z"),
					End:   twext.MustParseTime("20100203T142537Z"),
				},
				{
					Start: twext.MustParseTime("20100204T101530Z"),
					End:   twext.MustParseTime("20100204T142537Z"),
				},
			},
		},
		{
			name: "none",
			entries: twext.Entries{
				{
					Start: twext.MustParseTime("20100203T101530Z"),
					End:   twext.MustParseTime("20100203T142537Z"),
				},
				{
					Start: twext.MustParseTime("20100204T101530Z"),
					End:   twext.MustParseTime("20100204T142537Z"),
				},
			},
			filter: func(_ twext.Entry) bool {
				return false
			},
		},
		{
			name: "some",
			entries: twext.Entries{
				{
					ID:    42,
					Start: twext.MustParseTime("20100203T101530Z"),
					End:   twext.MustParseTime("20100203T142537Z"),
				},
				{
					Start: twext.MustParseTime("20100204T101530Z"),
					End:   twext.MustParseTime("20100204T142537Z"),
				},
			},
			filter: func(e twext.Entry) bool {
				return e.ID == 42
			},
			expectedEntries: twext.Entries{
				{
					ID:    42,
					Start: twext.MustParseTime("20100203T101530Z"),
					End:   twext.MustParseTime("20100203T142537Z"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualEntries := slices.Collect(tt.filter.Filter(tt.entries.All()))
			assert.EqualValues(t, tt.expectedEntries, actualEntries)
		})
	}
}
