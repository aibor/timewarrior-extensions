// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

// AggregationKey defines the interface for types that can be used as
// [Aggregation] keys.
type AggregationKey interface {
	cmp.Ordered
}

// Aggregation is a map of any values aggregated by a common [AggregationKey].
type Aggregation[K AggregationKey, V any] map[K]V

// SortedKeys returns the sorted list of [AggregationKey]s.
func (g Aggregation[K, V]) SortedKeys() []K {
	return slices.Sorted(maps.Keys(g))
}

// Sorted returns an iterator that iterates the [Aggregation] sorted by
// [AggregationKey]s.
func (g Aggregation[K, V]) Sorted() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, key := range g.SortedKeys() {
			if !yield(key, g[key]) {
				return
			}
		}
	}
}

// AggregationKeyFunc returns the aggregation key for the given entry.
type AggregationKeyFunc[K AggregationKey] func(Entry) K

// AggregationValueFunc adds an [Entry] to an aggregation result.
//
// The returned "newResult" is used as input "result" for the next iteration.
// The value can be a single scalar value or a slice.
type AggregationValueFunc[V any] func(result V, entry Entry) (newResult V)

// Aggregate aggregates entries into a map with user defined keys and values.
//
// It iterates the given entries and returns an [Aggregation] map with keys
// created by the given [AggregationKeyFunc] and values by the given
// [AggregationValueFunc] for each entry.
//
// Skipping entries is possible by returning the input result unaltered in
// the [AggregationValueFunc].
func Aggregate[K AggregationKey, V any](
	entries EntryIterator,
	keyFn AggregationKeyFunc[K],
	valueFn AggregationValueFunc[V],
) Aggregation[K, V] {
	aggregated := map[K]V{}

	for entry := range entries {
		key := keyFn(entry)
		aggregated[key] = valueFn(aggregated[key], entry)
	}

	return aggregated
}
