// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"testing"
	"time"

	"github.com/aibor/timewarrior-extensions/twext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadTimeTargetConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   twext.Config
		expected timeTargets
		errorMsg string
	}{
		{
			name:   "empty config",
			config: twext.Config{},
			expected: timeTargets{
				weekdays:        map[time.Weekday]time.Duration{},
				defaultDuration: 8 * time.Hour,
			},
		},
		{
			name: "just default",
			config: twext.Config{
				"flextime.time_per_day": "7h",
			},
			expected: timeTargets{
				weekdays:        map[time.Weekday]time.Duration{},
				defaultDuration: 7 * time.Hour,
			},
		},
		{
			name: "invalid default",
			config: twext.Config{
				"flextime.time_per_day": "a lot",
			},
			errorMsg: "time: invalid duration",
		},
		{
			name: "all weekdays",
			config: twext.Config{
				"flextime.time_per_day.monday":    "1h",
				"flextime.time_per_day.tuesday":   "2h",
				"flextime.time_per_day.wednesday": "3h",
				"flextime.time_per_day.thursday":  "4h",
				"flextime.time_per_day.friday":    "5h",
				"flextime.time_per_day.saturday":  "6h",
				"flextime.time_per_day.sunday":    "7h",
			},
			expected: timeTargets{
				weekdays: map[time.Weekday]time.Duration{
					time.Monday:    1 * time.Hour,
					time.Tuesday:   2 * time.Hour,
					time.Wednesday: 3 * time.Hour,
					time.Thursday:  4 * time.Hour,
					time.Friday:    5 * time.Hour,
					time.Saturday:  6 * time.Hour,
					time.Sunday:    7 * time.Hour,
				},
				defaultDuration: 8 * time.Hour,
			},
		},
		{
			name: "invalid weekday",
			config: twext.Config{
				"flextime.time_per_day.saturday": "nothing",
			},
			errorMsg: "time: invalid duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := readTimeTargetConfig(tt.config)

			if tt.errorMsg != "" {
				require.ErrorContains(t, err, tt.errorMsg)

				return
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}
