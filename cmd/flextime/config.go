// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"time"

	"github.com/aibor/timewarrior-extensions/twext"
)

const (
	configKeyPrefix              = "flextime"
	configKeyDailyTarget         = "time_per_day"
	defaultDailyTarget           = "8h"
	configKeyOffsetTotal         = "offset_total"
	defaultOffsetTotal           = "0"
	configKeyAggregationStrategy = "aggregation_strategy"
	defaultAggregationStrategy   = "single-day-only"
)

type config struct {
	dailyTarget         time.Duration
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

	dailyTarget, err := configRead(
		rawCfg,
		configKeyDailyTarget,
		defaultDailyTarget,
		parseDuration,
	)
	if err != nil {
		return config{}, fmt.Errorf("get daily target: %w", err)
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
		dailyTarget:         dailyTarget,
		offset:              offset,
		aggregationStrategy: strategy,
		debug:               rawCfg[twext.ConfigKeyDebug].Bool(),
		verbose:             rawCfg[twext.ConfigKeyVerbose].Bool(),
	}

	return cfg, nil
}

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
