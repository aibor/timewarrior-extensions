// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext_test

import (
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
				Start: twext.MustParseTime(t, "20100203T101530Z"),
				End:   twext.MustParseTime(t, "20100203T142537Z"),
			},
			expected: "4h10m7s",
		},
		{
			name: "days",
			entry: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T101530Z"),
				End:   twext.MustParseTime(t, "20100204T142537Z"),
			},
			expected: "28h10m7s",
		},
		{
			name: "before",
			entry: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T101530Z"),
				End:   twext.MustParseTime(t, "20100202T142537Z"),
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

func TestEntries_SplitAtMidnight(t *testing.T) {
	tests := []struct {
		name     string
		input    twext.Entries
		expected twext.Entries
	}{
		{
			name:     "empty",
			input:    twext.Entries{},
			expected: twext.Entries{},
		},
		{
			name: "only single-day entries",
			input: twext.Entries{
				twext.Entry{
					Start: twext.MustParseTime(t, "20100203T101530Z"),
					End:   twext.MustParseTime(t, "20100203T142537Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime(t, "20100204T092755Z"),
					End:   twext.MustParseTime(t, "20100204T163211Z"),
				},
			},
			expected: twext.Entries{
				twext.Entry{
					Start: twext.MustParseTime(t, "20100203T101530Z"),
					End:   twext.MustParseTime(t, "20100203T142537Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime(t, "20100204T092755Z"),
					End:   twext.MustParseTime(t, "20100204T163211Z"),
				},
			},
		},
		{
			name: "multi-day entries",
			input: twext.Entries{
				twext.Entry{
					Start: twext.MustParseTime(t, "20100203T101530Z"),
					End:   twext.MustParseTime(t, "20100204T142537Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime(t, "20100205T092755Z"),
					End:   twext.MustParseTime(t, "20100207T163211Z"),
				},
			},
			expected: twext.Entries{
				twext.Entry{
					Start: twext.MustParseTime(t, "20100203T101530Z"),
					End:   twext.MustParseTime(t, "20100204T000000Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime(t, "20100204T000000Z"),
					End:   twext.MustParseTime(t, "20100204T142537Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime(t, "20100205T092755Z"),
					End:   twext.MustParseTime(t, "20100206T000000Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime(t, "20100206T000000Z"),
					End:   twext.MustParseTime(t, "20100207T000000Z"),
				},
				twext.Entry{
					Start: twext.MustParseTime(t, "20100207T000000Z"),
					End:   twext.MustParseTime(t, "20100207T163211Z"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("expected slice matches for "+tt.name, func(t *testing.T) {
			actual := tt.input.SplitAtMidnight()
			assert.Equal(t, tt.expected, actual)
		})
		t.Run("duration sum is equal for "+tt.name, func(t *testing.T) {
			var expected, actual time.Duration
			for _, entry := range tt.input {
				expected += entry.Duration()
			}

			for _, splitEntry := range tt.input.SplitAtMidnight() {
				actual += splitEntry.Duration()
			}

			assert.Equal(t, expected, actual)
		})
		t.Run("duration sum is equal for "+tt.name, func(t *testing.T) {
			for _, entry := range tt.input.SplitAtMidnight() {
				assert.LessOrEqual(t, entry.Duration(), 24*time.Hour)
			}
		})
	}
}

func TestEntries_Filter(t *testing.T) {
	tests := []struct {
		name            string
		entries         twext.Entries
		filter          func(twext.Entry) bool
		expectedEntries twext.Entries
	}{
		{
			name: "all",
			entries: twext.Entries{
				{
					Start: twext.MustParseTime(t, "20100203T101530Z"),
					End:   twext.MustParseTime(t, "20100203T142537Z"),
				},
				{
					Start: twext.MustParseTime(t, "20100204T101530Z"),
					End:   twext.MustParseTime(t, "20100204T142537Z"),
				},
			},
			filter: func(_ twext.Entry) bool {
				return true
			},
			expectedEntries: twext.Entries{
				{
					Start: twext.MustParseTime(t, "20100203T101530Z"),
					End:   twext.MustParseTime(t, "20100203T142537Z"),
				},
				{
					Start: twext.MustParseTime(t, "20100204T101530Z"),
					End:   twext.MustParseTime(t, "20100204T142537Z"),
				},
			},
		},
		{
			name: "none",
			entries: twext.Entries{
				{
					Start: twext.MustParseTime(t, "20100203T101530Z"),
					End:   twext.MustParseTime(t, "20100203T142537Z"),
				},
				{
					Start: twext.MustParseTime(t, "20100204T101530Z"),
					End:   twext.MustParseTime(t, "20100204T142537Z"),
				},
			},
			filter: func(_ twext.Entry) bool {
				return false
			},
			expectedEntries: twext.Entries{},
		},
		{
			name: "some",
			entries: twext.Entries{
				{
					ID:    42,
					Start: twext.MustParseTime(t, "20100203T101530Z"),
					End:   twext.MustParseTime(t, "20100203T142537Z"),
				},
				{
					Start: twext.MustParseTime(t, "20100204T101530Z"),
					End:   twext.MustParseTime(t, "20100204T142537Z"),
				},
			},
			filter: func(e twext.Entry) bool {
				return e.ID == 42
			},
			expectedEntries: twext.Entries{
				{
					ID:    42,
					Start: twext.MustParseTime(t, "20100203T101530Z"),
					End:   twext.MustParseTime(t, "20100203T142537Z"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualEntries := tt.entries.Filter(tt.filter)
			assert.Equal(t, tt.expectedEntries, actualEntries)
		})
	}
}
