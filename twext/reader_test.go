// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext_test

import (
	"bytes"
	_ "embed"
	"errors"
	"testing"

	"github.com/aibor/timewarrior-extensions/twext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/input_config_only.txt
	testInputConfigOnly []byte
	//go:embed testdata/input_invalid_json.txt
	testInputInvalidJSON []byte
	//go:embed testdata/input_valid_entries.txt
	testInputValidEntries []byte
)

func TestReader_Config(t *testing.T) {
	t.Run("call twice", func(t *testing.T) {
		input := bytes.NewReader(testInputValidEntries)
		reader := twext.NewReader(input)

		_, err := reader.ReadConfig()
		require.NoError(t, err)

		_, err = reader.ReadConfig()
		require.ErrorIs(t, err, twext.ErrReaderConfigConsumed)
	})

	tests := []struct {
		name        string
		input       []byte
		expected    twext.Config
		expectedErr error
	}{
		{
			name:        "empty",
			input:       []byte{},
			expectedErr: twext.ErrConfigEmpty,
		},
		{
			name:  "config only",
			input: testInputConfigOnly,
			expected: twext.Config{
				"color":            "on",
				"reports.day.axis": "internal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := bytes.NewReader(tt.input)
			reader := twext.NewReader(input)

			actual, err := reader.ReadConfig()

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestReader_Entries(t *testing.T) {
	t.Run("config not called", func(t *testing.T) {
		input := bytes.NewReader(testInputValidEntries)
		reader := twext.NewReader(input)

		_, err := reader.ReadEntries()
		require.ErrorIs(t, err, twext.ErrReaderConfigNotConsumed)
	})

	tests := []struct {
		name        string
		input       []byte
		expected    twext.Entries
		expectedErr error
	}{
		{
			name:  "config only",
			input: testInputConfigOnly,
		},
		{
			name:        "invalid json",
			input:       testInputInvalidJSON,
			expectedErr: assert.AnError,
		},
		{
			name:  "valid entries",
			input: testInputValidEntries,
			expected: []twext.Entry{
				{
					ID:    3,
					Start: twext.MustParseTime("20240630T102128Z"),
					End:   twext.MustParseTime("20240630T102131Z"),
				},
				{
					ID:    2,
					Start: twext.MustParseTime("20240630T143940Z"),
					End:   twext.MustParseTime("20240630T143943Z"),
				},
				{
					ID:    1,
					Start: twext.MustParseTime("20240630T144010Z"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := bytes.NewReader(tt.input)
			reader := twext.NewReader(input)

			_, err := reader.ReadConfig()
			require.NoError(t, err)

			actual, err := reader.ReadEntries()

			if tt.expectedErr != nil {
				if errors.Is(tt.expectedErr, assert.AnError) {
					require.Error(t, err)
				} else {
					require.ErrorIs(t, err, tt.expectedErr)
				}

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
