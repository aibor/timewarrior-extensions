// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"testing"
	"time"

	"github.com/aibor/timewarrior-extensions/twext"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	name := "my-custom-strategy"
	strategy := aggregationStrategy[string, int]{
		name:      name,
		keyFn:     nil,
		valueFn:   nil,
		transform: nil,
	}

	assert.Equal(t, name, strategy.String())
}

func TestStartDate(t *testing.T) {
	tests := []struct {
		name     string
		input    twext.Entry
		expected string
	}{
		{
			name: "single day",
			input: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T101530Z"),
				End:   twext.MustParseTime(t, "20100203T102030Z"),
			},
			expected: "2010-02-03",
		},
		{
			name: "multi day",
			input: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T101530Z"),
				End:   twext.MustParseTime(t, "20100204T031125Z"),
			},
			expected: "2010-02-03",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, startDate(tt.input))
		})
	}
}

func TestEndDate(t *testing.T) {
	tests := []struct {
		name     string
		input    twext.Entry
		expected string
	}{
		{
			name: "single day",
			input: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T101530Z"),
				End:   twext.MustParseTime(t, "20100203T102030Z"),
			},
			expected: "2010-02-03",
		},
		{
			name: "multi day",
			input: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T101530Z"),
				End:   twext.MustParseTime(t, "20100204T031125Z"),
			},
			expected: "2010-02-04",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, endDate(tt.input))
		})
	}
}

func TestSumDuration(t *testing.T) {
	tests := []struct {
		name     string
		entry    twext.Entry
		base     time.Duration
		expected time.Duration
	}{
		{
			name: "short duration initial",
			entry: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T101530Z"),
				End:   twext.MustParseTime(t, "20100203T102030Z"),
			},
			base:     0,
			expected: 5 * time.Minute,
		},
		{
			name: "short duration accumulates",
			entry: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T101530Z"),
				End:   twext.MustParseTime(t, "20100203T102030Z"),
			},
			base:     15 * time.Minute,
			expected: 20 * time.Minute,
		},
		{
			name: "cross-day duration initial",
			entry: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T234530Z"),
				End:   twext.MustParseTime(t, "20100204T001530Z"),
			},
			base:     0,
			expected: 30 * time.Minute,
		},
		{
			name: "cross-day duration initial",
			entry: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T234530Z"),
				End:   twext.MustParseTime(t, "20100204T001530Z"),
			},
			base:     30 * time.Minute,
			expected: 1 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum := sumDuration(tt.base, tt.entry)
			assert.Equal(t, tt.expected, sum)
		})
	}
}

func TestOnlySingleDays(t *testing.T) {
	tests := []struct {
		name     string
		entry    twext.Entry
		expected bool
	}{
		{
			name: "single day",
			entry: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T101530Z"),
				End:   twext.MustParseTime(t, "20100203T102030Z"),
			},
			expected: true,
		},
		{
			name: "multi day",
			entry: twext.Entry{
				Start: twext.MustParseTime(t, "20100203T234530Z"),
				End:   twext.MustParseTime(t, "20100204T001530Z"),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, onlySingleDays(tt.entry))
		})
	}
}
