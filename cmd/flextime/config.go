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
	configKeyPrefix      = "flextime"
	configKeyDailyTarget = "time_per_day"
	defaultDailyTarget   = "8h"
	configKeyOffsetTotal = "offset_total"
	defaultOffsetTotal   = "0"
)

type config struct {
	dailyTarget time.Duration
	offset      time.Duration
	debug       bool
	verbose     bool
}

func readConfig(reader *twext.Reader) (config, error) {
	rawCfg, err := reader.ReadConfig()
	if err != nil {
		return config{}, fmt.Errorf("read config section: %w", err)
	}

	dailyTarget, err := configReadDuration(
		rawCfg,
		configKeyDailyTarget,
		defaultDailyTarget,
	)
	if err != nil {
		return config{}, fmt.Errorf("get daily target: %w", err)
	}

	offset, err := configReadDuration(
		rawCfg,
		configKeyOffsetTotal,
		defaultOffsetTotal,
	)
	if err != nil {
		return config{}, fmt.Errorf("get total offset: %w", err)
	}

	cfg := config{
		dailyTarget: dailyTarget,
		offset:      offset,
		debug:       rawCfg[twext.ConfigKeyDebug].Bool(),
		verbose:     rawCfg[twext.ConfigKeyVerbose].Bool(),
	}

	return cfg, nil
}

func configReadDuration(
	twConfig twext.Config,
	cfgKey string,
	defValue string,
) (time.Duration, error) {
	key := twext.NewConfigKey(configKeyPrefix, cfgKey)

	cfgValue, exists := twConfig[key]
	if !exists {
		cfgValue = twext.ConfigValue(defValue)
	}

	duration, err := cfgValue.Duration()
	if err != nil {
		return 0, fmt.Errorf("convert to duration: %w", err)
	}

	return duration, nil
}
