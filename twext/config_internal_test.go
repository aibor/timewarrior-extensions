// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfigLine(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedKey   ConfigKey
		expectedValue ConfigValue
		expectedErr   error
	}{
		{
			name:        "empty",
			expectedErr: ErrConfigInvalidLine,
		},
		{
			name:        "newline",
			input:       "\n",
			expectedErr: ErrConfigInvalidLine,
		},
		{
			name:        "missing colon",
			input:       "key value",
			expectedErr: ErrConfigInvalidLine,
		},
		{
			name:        "missing space",
			input:       "key:value",
			expectedErr: ErrConfigInvalidLine,
		},
		{
			name:          "simple key vale",
			input:         "key: value",
			expectedKey:   "key",
			expectedValue: "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualKey, actualValue, err := readConfigLine(tt.input)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedKey, actualKey, "key")
			assert.Equal(t, tt.expectedValue, actualValue, "value")
		})
	}
}

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedConfig Config
		expectedErr    error
	}{
		{
			name:           "empty",
			expectedConfig: Config{},
			expectedErr:    ErrConfigEmpty,
		},
		{
			name:           "newline",
			input:          "\n",
			expectedConfig: Config{},
			expectedErr:    ErrConfigEmpty,
		},
		{
			name:  "one line",
			input: "key: value\n",
			expectedConfig: Config{
				"key": "value",
			},
		},
		{
			name:  "one line with blank",
			input: "key: value\n\nsome: thing\n\n",
			expectedConfig: Config{
				"key": "value",
			},
		},
		{
			name:  "multi line with blank",
			input: "key: value\nsome: thing\n\nother: stuff",
			expectedConfig: Config{
				"key":  "value",
				"some": "thing",
			},
		},
		{
			name:        "garbage",
			input:       "some thing\n\n",
			expectedErr: ErrConfigInvalidLine,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := bytes.NewBufferString(tt.input)
			actualConfig, err := readConfig(input)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedConfig, actualConfig)
		})
	}
}
