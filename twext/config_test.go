// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext_test

import (
	"testing"

	"github.com/aibor/timewarrior-extensions/twext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigKey(t *testing.T) {
	tests := []struct {
		name        string
		expectedKey twext.ConfigKey
		input       []string
	}{
		{
			name:        "empty",
			input:       []string{},
			expectedKey: "",
		},
		{
			name:        "single",
			input:       []string{"f"},
			expectedKey: "f",
		},
		{
			name:        "multi",
			input:       []string{"a", "b", "c"},
			expectedKey: "a.b.c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expectedKey.String(), func(t *testing.T) {
			actualKey := twext.NewConfigKey(tt.input...)
			assert.Equal(t, tt.expectedKey, actualKey)
		})
	}
}

func TestConfigValueBool(t *testing.T) {
	tests := []struct {
		input    twext.ConfigValue
		expected bool
	}{
		{input: ""},
		{input: "f"},
		{input: "false"},
		{input: "0"},
		{input: "True"},
		{expected: true, input: "true"},
		{expected: true, input: "on"},
		{expected: true, input: "1"},
		{expected: true, input: "yes"},
		{expected: true, input: "y"},
	}

	for _, tt := range tests {
		t.Run(tt.input.String(), func(t *testing.T) {
			actual := tt.input.Bool()

			if tt.expected {
				assert.True(t, actual)
			} else {
				assert.False(t, actual)
			}
		})
	}
}

func TestConfigValueInt(t *testing.T) {
	tests := []struct {
		input       twext.ConfigValue
		expectedInt int
		invalid     bool
	}{
		{
			input:       "5",
			expectedInt: 5,
		},
		{
			input:       "1234567",
			expectedInt: 1234567,
		},
		{
			input:   "0xff",
			invalid: true,
		},
		{
			input:   "j",
			invalid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input.String(), func(t *testing.T) {
			actualInt, err := tt.input.Int()

			if tt.invalid {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedInt, actualInt)
		})
	}
}
