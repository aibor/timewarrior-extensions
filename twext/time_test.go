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

func TestParseTime(t *testing.T) {
	validTime, err := time.Parse(time.RFC3339, "2010-02-01T12:34:56Z")
	require.NoError(t, err)

	tests := []struct {
		name         string
		input        string
		expectedTime twext.Time
		invalid      bool
	}{
		{
			name:    "empty",
			input:   "",
			invalid: true,
		},
		{
			name:    "partial",
			input:   "20100201T123456",
			invalid: true,
		},
		{
			name:    "wrong format",
			input:   "2010-02-01T12:34:56Z",
			invalid: true,
		},
		{
			name:         "valid",
			input:        "20100201T123456Z",
			expectedTime: twext.Time{validTime},
		},
		{
			name:  "null",
			input: "00010101T000000Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := twext.ParseTime(tt.input)

			if tt.invalid {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedTime, actual)
		})
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		expectedDate twext.Time
		invalid      bool
	}{
		{
			name: "empty",
		},
		{
			name:  "empty string",
			input: []byte(`""`),
		},
		{
			name:  "null string",
			input: []byte(`null`),
		},
		{
			name:    "not a string",
			input:   []byte{1},
			invalid: true,
		},
		{
			name:    "invalid date string",
			input:   []byte(`"20342T232Z"`),
			invalid: true,
		},
		{
			name:    "valid date",
			input:   []byte(`20240630T143940Z`),
			invalid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actualDate twext.Time
			err := actualDate.UnmarshalJSON(tt.input)

			if tt.invalid {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedDate, actualDate)
		})
	}
}

func TestTimeSameDate(t *testing.T) {
	tests := []struct {
		name        string
		dateA       twext.Time
		dateB       twext.Time
		expectEqual bool
	}{
		{
			name:        "same timestamp",
			dateA:       twext.MustParseTime(t, "20240630T112233Z"),
			dateB:       twext.MustParseTime(t, "20240630T112233Z"),
			expectEqual: true,
		},
		{
			name:        "same date",
			dateA:       twext.MustParseTime(t, "20240630T112233Z"),
			dateB:       twext.MustParseTime(t, "20240630T223344Z"),
			expectEqual: true,
		},
		{
			name:        "zero values",
			dateA:       twext.MustParseTime(t, "00010101T000000Z"),
			dateB:       twext.MustParseTime(t, "00010101T000000Z"),
			expectEqual: true,
		},
		{
			name:  "different date",
			dateA: twext.MustParseTime(t, "20240630T112233Z"),
			dateB: twext.MustParseTime(t, "20240629T112233Z"),
		},
		{
			name:  "with zero",
			dateA: twext.MustParseTime(t, "20240630T112233Z"),
			dateB: twext.MustParseTime(t, "00010101T000000Z"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualEqual := tt.dateA.SameDate(&tt.dateB)

			if tt.expectEqual {
				assert.True(t, actualEqual)
			} else {
				assert.False(t, actualEqual)
			}
		})
	}
}
