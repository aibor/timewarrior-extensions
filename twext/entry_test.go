// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: MIT

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
