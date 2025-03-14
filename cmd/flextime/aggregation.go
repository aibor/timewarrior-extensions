// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aibor/timewarrior-extensions/twext"
)

type entriesTransformation func(entries twext.EntryIterator) twext.EntryIterator

type aggregationStrategy[K twext.GroupKey, V any] struct {
	name      string
	keyFn     twext.GroupKeyFunc[K]
	valueFn   twext.GroupValueFunc[V]
	transform entriesTransformation
}

func (s *aggregationStrategy[K, V]) String() string {
	return s.name
}

func (s *aggregationStrategy[K, V]) Aggregate(
	entries twext.EntryIterator,
) twext.Groups[K, V] {
	if s.transform != nil {
		entries = s.transform(entries)
	}

	return twext.Group(entries, s.keyFn, s.valueFn)
}

func startDate(entry twext.Entry) string {
	return entry.Start.Format(time.DateOnly)
}

func endDate(entry twext.Entry) string {
	return entry.End.Format(time.DateOnly)
}

func sumDuration(result time.Duration, entry twext.Entry) time.Duration {
	return result + entry.Duration()
}

func onlySingleDays(entry twext.Entry) bool {
	sameDate := entry.Start.SameDate(entry.CurrentEnd())
	if !sameDate {
		log.Printf("entry %d spans multiple days. Skipping.", entry.ID)
	}

	return sameDate
}

var errUnknownAggregationStrategy = errors.New("unknown aggregation strategy")

func createAggregationStrategy(
	strategy string,
) (*aggregationStrategy[string, time.Duration], error) {
	switch strategy {
	case "single-day-only":
		return &aggregationStrategy[string, time.Duration]{
			name:      strategy,
			keyFn:     startDate,
			valueFn:   sumDuration,
			transform: twext.EntryFilter(onlySingleDays).Filter,
		}, nil
	case "into-start-date":
		return &aggregationStrategy[string, time.Duration]{
			name:    strategy,
			keyFn:   startDate,
			valueFn: sumDuration,
		}, nil
	case "into-end-date":
		return &aggregationStrategy[string, time.Duration]{
			name:    strategy,
			keyFn:   endDate,
			valueFn: sumDuration,
		}, nil
	case "split-at-midnight":
		return &aggregationStrategy[string, time.Duration]{
			name:      strategy,
			keyFn:     startDate,
			valueFn:   sumDuration,
			transform: twext.SplitAtMidnight,
		}, nil
	default:
		return nil, fmt.Errorf("%w: %s", errUnknownAggregationStrategy, strategy)
	}
}
