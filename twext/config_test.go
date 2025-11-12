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

func TestConfigKeySubKey(t *testing.T) {
	tests := []struct {
		name   string
		key    twext.ConfigKey
		prefix twext.ConfigKey
		found  bool
		subKey twext.ConfigKey
	}{
		{
			name:   "empty prefixes empty",
			key:    "",
			prefix: "",
			found:  true,
			subKey: "",
		},
		{
			name:   "empty is not a prefix for something",
			key:    "a",
			prefix: "",
			found:  false,
			subKey: "a",
		},
		{
			name:   "empty prefixes empty namespace",
			key:    ".a",
			prefix: "",
			found:  true,
			subKey: "a",
		},
		{
			name:   "something is not a prefix for empty",
			key:    "",
			prefix: "z",
			found:  false,
			subKey: "",
		},
		{
			name:   "a prefixes a",
			key:    "a",
			prefix: "a",
			found:  true,
			subKey: "",
		},
		{
			name:   "a prefixes a.b",
			key:    "a.b",
			prefix: "a",
			found:  true,
			subKey: "b",
		},
		{
			name:   "b is not a prefix for a.b",
			key:    "a.b",
			prefix: "b",
			found:  false,
			subKey: "a.b",
		},
		{
			name:   "part is not a prefix for partial",
			key:    "partial",
			prefix: "part",
			found:  false,
			subKey: "partial",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subKey, found := tt.key.SubKey(tt.prefix)
			assert.Equalf(
				t, tt.found, found,
				"Does '%s' have the prefix '%s'?",
				tt.key, tt.prefix,
			)
			assert.Equal(t, tt.subKey, subKey)
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

func TestConfigValueDuration(t *testing.T) {
	tests := []struct {
		input            twext.ConfigValue
		expectedDuration time.Duration
		invalid          bool
	}{
		{
			input:            "5s",
			expectedDuration: 5 * time.Second,
		},
		{
			input:            "8h",
			expectedDuration: 8 * time.Hour,
		},
		{
			input:            "0",
			expectedDuration: 0,
		},
		{
			input:   "5",
			invalid: true,
		},
		{
			input:   "j",
			invalid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input.String(), func(t *testing.T) {
			actualDuration, err := tt.input.Duration()

			if tt.invalid {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedDuration, actualDuration)
		})
	}
}
