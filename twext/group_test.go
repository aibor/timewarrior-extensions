// SPDX-FileCopyrightText: 2025 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/aibor/timewarrior-extensions/twext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroups_SortedKeys(t *testing.T) {
	tests := []struct {
		name         string
		groups       twext.Groups[int, int]
		expectedKeys []int
	}{
		{
			name:   "empty",
			groups: twext.Groups[int, int]{},
		},
		{
			name:         "single",
			groups:       twext.Groups[int, int]{1: 2},
			expectedKeys: []int{1},
		},
		{
			name:         "many",
			groups:       twext.Groups[int, int]{1: 2, 3: 4, 5: 6, 7: 8, 9: 0},
			expectedKeys: []int{1, 3, 5, 7, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualKeys := tt.groups.SortedKeys()
			assert.Equal(t, tt.expectedKeys, actualKeys)
		})
	}
}

func TestGroups_Sorted(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		groups := twext.Groups[int, int]{}
		for range groups.Sorted() {
			require.Fail(t, "loop body should not be called")
		}
	})

	groups := twext.Groups[int, int]{1: 2, 3: 4, 5: 6, 7: 8, 9: 0}

	type testVal struct{ i1, i2 int }

	t.Run("break", func(t *testing.T) {
		expected := []testVal{{1, 2}}

		actual := []testVal{}
		for key, value := range groups.Sorted() {
			actual = append(actual, testVal{key, value})

			break
		}

		assert.Equal(t, expected, actual)
	})

	t.Run("all", func(t *testing.T) {
		expected := []testVal{
			{1, 2},
			{3, 4},
			{5, 6},
			{7, 8},
			{9, 0},
		}

		actual := []testVal{}
		for key, value := range groups.Sorted() {
			actual = append(actual, testVal{key, value})
		}

		assert.Equal(t, expected, actual)
	})
}

func TestGroup(t *testing.T) {
	tests := []struct {
		name           string
		entries        twext.Entries
		keyFn          twext.GroupKeyFunc[int]
		valueFn        twext.GroupValueFunc[int]
		expectedGroups twext.Groups[int, int]
	}{
		{
			name:    "empty",
			entries: twext.Entries{},
			keyFn: func(_ twext.Entry) int {
				return 1
			},
			valueFn: func(_ int, _ twext.Entry) int {
				return 2
			},
			expectedGroups: twext.Groups[int, int]{},
		},
		{
			name: "single",
			entries: twext.Entries{
				{},
			},
			keyFn: func(_ twext.Entry) int {
				return 1
			},
			valueFn: func(_ int, _ twext.Entry) int {
				return 2
			},
			expectedGroups: twext.Groups[int, int]{
				1: 2,
			},
		},
		{
			name: "by tag count",
			entries: twext.Entries{
				{Tags: []string{"a", "b", "c"}},
				{Tags: []string{"a", "d", "e"}},
				{Tags: []string{"a"}},
				{Tags: []string{"z"}},
				{Tags: []string{"z"}},
				{Tags: []string{"a", "z"}},
			},
			keyFn: func(e twext.Entry) int {
				return len(e.Tags)
			},
			valueFn: func(r int, _ twext.Entry) int {
				return r + 1
			},
			expectedGroups: twext.Groups[int, int]{
				1: 3,
				2: 1,
				3: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := twext.Group(tt.entries, tt.keyFn, tt.valueFn)
			assert.Equal(t, tt.expectedGroups, actual)
		})
	}
}

func BenchmarkGroup(b *testing.B) {
	valueFn := func(_ int, _ twext.Entry) int {
		return 0
	}

	benchs := []struct {
		name              string
		keyFn             twext.GroupKeyFunc[int]
		expectedGroupsMax int
	}{
		{
			name: "single group",
			keyFn: func(_ twext.Entry) int {
				return 0
			},
			expectedGroupsMax: 1,
		},
		{
			name: "100 groups",
			keyFn: func(e twext.Entry) int {
				return e.ID % 100
			},
			expectedGroupsMax: 100,
		},
		{
			name: "max groups",
			keyFn: func(e twext.Entry) int {
				return e.ID
			},
		},
	}

	for _, bench := range benchs {
		b.Run(bench.name, func(b *testing.B) {
			for _, m := range []int{3, 4, 5, 6} {
				n := int(math.Pow10(m))

				entries := make(twext.Entries, n)
				for i := range n {
					entries[i].ID = i
				}

				expectedLen := n
				if bench.expectedGroupsMax > 0 && bench.expectedGroupsMax < n {
					expectedLen = bench.expectedGroupsMax
				}

				b.Run(strconv.Itoa(n), func(b *testing.B) {
					actualLen := 0

					for b.Loop() {
						groups := twext.Group(entries, bench.keyFn, valueFn)
						actualLen = len(groups)
					}

					require.Equal(b, expectedLen, actualLen)
				})
			}
		})
	}
}
