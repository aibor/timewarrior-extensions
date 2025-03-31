// SPDX-FileCopyrightText: 2025 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetClock(t *testing.T) {
	testZone := time.FixedZone("test", -7)

	tests := []struct {
		name     string
		base     time.Time
		clock    time.Time
		expected time.Time
	}{
		{
			name: "empty",
		},
		{
			name:     "midnight",
			base:     time.Date(2015, 3, 1, 13, 4, 15, 0, time.UTC),
			clock:    time.Time{},
			expected: time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "ignore clock date",
			base:     time.Date(2015, 3, 1, 13, 4, 15, 0, time.UTC),
			clock:    time.Date(2016, 4, 2, 14, 8, 24, 6, testZone),
			expected: time.Date(2015, 3, 1, 14, 8, 24, 6, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := setClock(tt.base, tt.clock)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
