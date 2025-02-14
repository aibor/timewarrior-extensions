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

type entriesTransformation func(entries twext.Entries) twext.Entries

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
	entries twext.Entries,
) twext.Groups[K, V] {
	return twext.Group(
		s.transform(entries),
		s.keyFn,
		s.valueFn,
	)
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
			name:    strategy,
			keyFn:   startDate,
			valueFn: sumDuration,
			transform: func(entries twext.Entries) twext.Entries {
				return entries.Filter(onlySingleDays)
			},
		}, nil
	case "into-start-date":
		return &aggregationStrategy[string, time.Duration]{
			name:      strategy,
			keyFn:     startDate,
			valueFn:   sumDuration,
			transform: func(entries twext.Entries) twext.Entries { return entries },
		}, nil
	case "into-end-date":
		return &aggregationStrategy[string, time.Duration]{
			name:      strategy,
			keyFn:     endDate,
			valueFn:   sumDuration,
			transform: func(entries twext.Entries) twext.Entries { return entries },
		}, nil
	case "split-at-midnight":
		return &aggregationStrategy[string, time.Duration]{
			name:    strategy,
			keyFn:   startDate,
			valueFn: sumDuration,
			transform: func(entries twext.Entries) twext.Entries {
				return entries.SplitAtMidnight()
			},
		}, nil
	default:
		return nil, fmt.Errorf("%w: %s", errUnknownAggregationStrategy, strategy)
	}
}
