// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"strings"
	"testing"

	"github.com/aibor/timewarrior-extensions/twext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedStdout string
		expectedStderr string
		expectedErr    error
	}{
		{
			name:        "empty input",
			expectedErr: twext.ErrConfigEmpty,
		},
		{
			name: "quiet",
			input: `verbose: off
flextime.hours_per_day: 4
flextime.minutes_offset_total: 60

[
{"id":3,"start":"20240629T102128Z","end":"20240629T102131Z"},
{"id":2,"start":"20240630T143940Z","end":"20240630T143943Z"},
{"id":1,"start":"20240630T144010Z","end":"20240630T163943Z"}
]`,
			expectedStdout: `
     date    actual    target       diff
    total    2h:59m    8h:00m    -5h:00m
`,
		},
		{
			name: "verbose",
			input: `verbose: on

[
{"id":3,"start":"20240629T102128Z","end":"20240629T102131Z"},
{"id":2,"start":"20240630T143940Z","end":"20240630T143943Z"},
{"id":1,"start":"20240630T144010Z","end":"20240630T163943Z"}
]`,
			expectedStdout: `
          date    actual     target        diff
    2024-06-29    0h:00m     8h:00m     -7h:59m
    2024-06-30    1h:59m     8h:00m     -6h:00m
         total    1h:59m    16h:00m    -14h:00m
`,
		},
		{
			name: "positive offset",
			input: `verbose: on
flextime.minutes_offset_total: 90

[
{"id":3,"start":"20240629T102128Z","end":"20240629T102131Z"},
{"id":2,"start":"20240630T143940Z","end":"20240630T143943Z"},
{"id":1,"start":"20240630T144010Z","end":"20240630T163943Z"}
]`,
			expectedStdout: `
          date    actual     target        diff
        offset    1h:30m     0h:00m      1h:30m
    2024-06-29    0h:00m     8h:00m     -7h:59m
    2024-06-30    1h:59m     8h:00m     -6h:00m
         total    3h:29m    16h:00m    -12h:30m
`,
		},
		{
			name: "negative offset",
			input: `verbose: on
flextime.minutes_offset_total: -70

[
{"id":3,"start":"20240629T102128Z","end":"20240629T102131Z"},
{"id":2,"start":"20240630T143940Z","end":"20240630T143943Z"},
{"id":1,"start":"20240630T144010Z","end":"20240630T163943Z"}
]`,
			expectedStdout: `
          date     actual     target        diff
        offset    -1h:10m     0h:00m     -1h:10m
    2024-06-29     0h:00m     8h:00m     -7h:59m
    2024-06-30     1h:59m     8h:00m     -6h:00m
         total     0h:49m    16h:00m    -15h:10m
`,
		},
		{
			name: "debug",
			input: `debug: on

[
{"id":3,"start":"20240629T102128Z","end":"20240630T102131Z"}
]`,
			expectedStdout: `
     date    actual    target       diff
    total    0h:00m    8h:00m    -8h:00m
`,
			expectedStderr: `debug [flextime] - cfg - DailyTarget: 8h0m0s
debug [flextime] - cfg - Debug: true
debug [flextime] - cfg - Verbose: false
debug [flextime] - entry 3 spans multiple days. Skipping.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr strings.Builder

			stdin := strings.NewReader(tt.input)

			err := run(stdin, &stdout, &stderr)
			require.ErrorIs(t, err, tt.expectedErr)

			if tt.expectedErr != nil {
				return
			}

			assert.Equal(t, tt.expectedStdout, stdout.String(), "stdout")
			assert.Equal(t, tt.expectedStderr, stderr.String(), "stderr")
		})
	}
}
