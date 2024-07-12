// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: MIT

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

[
{"id":3,"start":"20240629T102128Z","end":"20240629T102131Z"},
{"id":2,"start":"20240630T143940Z","end":"20240630T143943Z"},
{"id":1,"start":"20240630T144010Z","end":"20240630T163943Z"}
]`,
			expectedStdout: `
     date    actual    target     diff
    total      1:59      8:00    -6:00
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
          date    actual    target      diff
    2024-06-29      0:00      8:00     -7:59
    2024-06-30      1:59      8:00     -6:00
         total      1:59     16:00    -14:00
`,
		},
		{
			name: "debug",
			input: `debug: on

[
{"id":3,"start":"20240629T102128Z","end":"20240630T102131Z"}
]`,
			expectedStdout: `
     date    actual    target     diff
    total      0:00      8:00    -8:00
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
