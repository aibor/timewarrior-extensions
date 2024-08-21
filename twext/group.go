// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: MIT

package twext

import (
	"cmp"
	"maps"
	"slices"
)

// GroupKey defines the interface for types that can be used as [Groups] keys.
type GroupKey interface {
	cmp.Ordered
}

// GroupValue defines the interface for types that can be used as [Groups]
// values.
type GroupValue interface {
	any
}

// Groups is a map of any values grouped by a common [GroupKey].
type Groups[K GroupKey, V any] map[K]V

// Keys returns the sorted list of [GroupKey]s.
func (g Groups[K, V]) Keys() []K {
	return slices.Sorted(maps.Keys(g))
}

// GroupKeyFunc returns the group key for the given entry.
type GroupKeyFunc[K GroupKey] func(Entry) K

// GroupValueFunc aggregates entries of a group.
//
// The returned "newResult" is used as input "result" for the next iteration.
// The value can be a single scalar value or a slice.
type GroupValueFunc[V GroupValue] func(result V, entry Entry) (newResult V)

// Group aggregates entries into groups.
//
// It returns [Groups] grouped by [GroupKey]. The return value of the given
// [GroupKeyFunc] is used as mapping key for the processed [Entry]. The value
// depends on the return value of the given [GroupValueFunc].
//
// Skipping entries is possible by returning the input result unaltered in
// the [GroupValueFunc].
func Group[K GroupKey, V any](
	entries Entries,
	keyFn GroupKeyFunc[K],
	valueFn GroupValueFunc[V],
) Groups[K, V] {
	groups := map[K]V{}

	for _, entry := range entries {
		key := keyFn(entry)
		groups[key] = valueFn(groups[key], entry)
	}

	return groups
}
