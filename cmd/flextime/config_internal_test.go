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
				dates:           map[time.Time]time.Duration{},
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
				dates:           map[time.Time]time.Duration{},
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
				dates: map[time.Time]time.Duration{},
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
		{
			name: "date-specific overrides",
			config: twext.Config{
				"flextime.time_per_day.date.2025-10-31": "1h",
				"flextime.time_per_day.date.2025-11-01": "2h",
			},
			expected: timeTargets{
				dates: map[time.Time]time.Duration{
					time.Date(
						2025, 10, 31,
						0, 0, 0, 0,
						time.UTC,
					): 1 * time.Hour,
					time.Date(
						2025, 11, 1,
						0, 0, 0, 0,
						time.UTC,
					): 2 * time.Hour,
				},
				weekdays:        map[time.Weekday]time.Duration{},
				defaultDuration: 8 * time.Hour,
			},
		},
		{
			name: "invalid date overrides",
			config: twext.Config{
				"flextime.time_per_day.date.friday_the_24st_of_may_2024": "1h",
			},
			errorMsg: "get date for override",
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

func TestTargetFor(t *testing.T) {
	tests := []struct {
		name     string
		config   timeTargets
		day      time.Time
		expected time.Duration
	}{
		{
			name: "just default",
			config: timeTargets{
				dates:           map[time.Time]time.Duration{},
				weekdays:        map[time.Weekday]time.Duration{},
				defaultDuration: 8 * time.Hour,
			},
			day: time.Date(
				2025, 10, 31,
				0, 0, 0, 0,
				time.UTC,
			),
			expected: 8 * time.Hour,
		},
		{
			name: "weekday specific",
			config: timeTargets{
				dates: map[time.Time]time.Duration{},
				weekdays: map[time.Weekday]time.Duration{
					time.Friday: 1 * time.Hour,
				},
				defaultDuration: 8 * time.Hour,
			},
			day: time.Date(
				2025, 10, 31,
				0, 0, 0, 0,
				time.UTC,
			),
			expected: 1 * time.Hour,
		},
		{
			name: "other weekday thus default",
			config: timeTargets{
				dates: map[time.Time]time.Duration{},
				weekdays: map[time.Weekday]time.Duration{
					time.Monday: 1 * time.Hour,
				},
				defaultDuration: 8 * time.Hour,
			},
			day: time.Date(
				2025, 10, 31,
				0, 0, 0, 0,
				time.UTC,
			),
			expected: 8 * time.Hour,
		},
		{
			name: "date specific",
			config: timeTargets{
				dates: map[time.Time]time.Duration{
					time.Date(
						2025, 10, 31,
						0, 0, 0, 0,
						time.UTC,
					): 1 * time.Hour,
				},
				weekdays:        map[time.Weekday]time.Duration{},
				defaultDuration: 8 * time.Hour,
			},
			day: time.Date(
				2025, 10, 31,
				0, 0, 0, 0,
				time.UTC,
			),
			expected: 1 * time.Hour,
		},
		{
			name: "other date thus default",
			config: timeTargets{
				dates: map[time.Time]time.Duration{
					time.Date(
						2025, 11, 1,
						0, 0, 0, 0,
						time.UTC,
					): 1 * time.Hour,
				},
				weekdays:        map[time.Weekday]time.Duration{},
				defaultDuration: 8 * time.Hour,
			},
			day: time.Date(
				2025, 10, 31,
				0, 0, 0, 0,
				time.UTC,
			),
			expected: 8 * time.Hour,
		},
		{
			name: "date specific but dirty request",
			config: timeTargets{
				dates: map[time.Time]time.Duration{
					time.Date(
						2025, 10, 31,
						0, 0, 0, 0,
						time.UTC,
					): 1 * time.Hour,
				},
				weekdays:        map[time.Weekday]time.Duration{},
				defaultDuration: 8 * time.Hour,
			},
			day: time.Date(
				2025, 10, 31,
				19, 23, 18, 123503340,
				time.FixedZone("ACWST", (8*60+45)*60),
			),
			expected: 1 * time.Hour,
		},
		{
			name: "date specific has precedence over weekday specific",
			config: timeTargets{
				dates: map[time.Time]time.Duration{
					time.Date(
						2025, 10, 31,
						0, 0, 0, 0,
						time.UTC,
					): 1 * time.Hour,
				},
				weekdays: map[time.Weekday]time.Duration{
					time.Friday: 4 * time.Hour,
				},
				defaultDuration: 8 * time.Hour,
			},
			day: time.Date(
				2025, 10, 31,
				0, 0, 0, 0,
				time.UTC,
			),
			expected: 1 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.config.targetFor(tt.day))
		})
	}
}
