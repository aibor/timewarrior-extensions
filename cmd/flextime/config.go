// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aibor/timewarrior-extensions/twext"
)

const numberOfWeekdays = 7

const (
	configKeyPrefix              = "flextime"
	configKeyTimeTarget          = "time_per_day"
	defaultTimeTarget            = "8h"
	configKeyOffsetTotal         = "offset_total"
	defaultOffsetTotal           = "0"
	configKeyAggregationStrategy = "aggregation_strategy"
	defaultAggregationStrategy   = "single-day-only"
)

type timeTargets struct {
	weekdays        map[time.Weekday]time.Duration
	defaultDuration time.Duration
}

func (t timeTargets) targetFor(weekday time.Weekday) time.Duration {
	if target, exists := t.weekdays[weekday]; exists {
		return target
	}

	return t.defaultDuration
}

func (t timeTargets) String() string {
	str := &strings.Builder{}
	_, _ = fmt.Fprintf(str, "Default: %s", t.defaultDuration)

	for day, duration := range t.weekdays {
		_, _ = fmt.Fprintf(str, " %s: %s", day, duration)
	}

	return str.String()
}

type config struct {
	timeTargets         timeTargets
	offset              time.Duration
	aggregationStrategy *aggregationStrategy[string, time.Duration]
	debug               bool
	verbose             bool
}

func readConfig(reader *twext.Reader) (config, error) {
	rawCfg, err := reader.ReadConfig()
	if err != nil {
		return config{}, fmt.Errorf("read config section: %w", err)
	}

	target, err := readTimeTargetConfig(rawCfg)
	if err != nil {
		return config{}, fmt.Errorf("get time target: %w", err)
	}

	offset, err := configRead(
		rawCfg,
		configKeyOffsetTotal,
		defaultOffsetTotal,
		parseDuration,
	)
	if err != nil {
		return config{}, fmt.Errorf("get total offset: %w", err)
	}

	strategy, err := configRead(
		rawCfg,
		configKeyAggregationStrategy,
		defaultAggregationStrategy,
		parseAggregationStrategy,
	)
	if err != nil {
		return config{}, fmt.Errorf("get aggregation strategy: %w", err)
	}

	cfg := config{
		timeTargets:         target,
		offset:              offset,
		aggregationStrategy: strategy,
		debug:               rawCfg[twext.ConfigKeyDebug].Bool(),
		verbose:             rawCfg[twext.ConfigKeyVerbose].Bool(),
	}

	return cfg, nil
}

func readTimeTargetConfig(twConfig twext.Config) (timeTargets, error) {
	var err error

	targets := timeTargets{
		weekdays: make(map[time.Weekday]time.Duration, numberOfWeekdays),
	}

	targets.defaultDuration, err = configRead(
		twConfig,
		configKeyTimeTarget,
		defaultTimeTarget,
		parseDuration,
	)
	if err != nil {
		return timeTargets{}, fmt.Errorf("get default target: %w", err)
	}

	for day := range time.Weekday(numberOfWeekdays) {
		subKey := strings.ToLower(day.String())
		key := twext.NewConfigKey(configKeyPrefix, configKeyTimeTarget, subKey)
		cfgValue, exists := twConfig[key]

		if !exists {
			continue
		}

		targets.weekdays[day], err = parseDuration(cfgValue)
		if err != nil {
			return timeTargets{}, fmt.Errorf("get target for %s: %w", day, err)
		}
	}

	return targets, nil
}

//nolint:ireturn,nolintlint
func configRead[R any](
	twConfig twext.Config,
	cfgKey string,
	defValue twext.ConfigValue,
	parseValue func(twext.ConfigValue) (R, error),
) (R, error) {
	key := twext.NewConfigKey(configKeyPrefix, cfgKey)

	cfgValue, exists := twConfig[key]
	if !exists {
		cfgValue = defValue
	}

	return parseValue(cfgValue)
}

func parseDuration(value twext.ConfigValue) (time.Duration, error) {
	duration, err := value.Duration()
	if err != nil {
		return 0, fmt.Errorf("convert to duration: %w", err)
	}

	return duration, nil
}

func parseAggregationStrategy(
	value twext.ConfigValue,
) (*aggregationStrategy[string, time.Duration], error) {
	strategy, err := createAggregationStrategy(value.String())
	if err != nil {
		return nil, fmt.Errorf("create aggregation strategy: %w", err)
	}

	return strategy, nil
}
