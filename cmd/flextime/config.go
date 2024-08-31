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
	configKeyHoursPerDay = "hours_per_day"
	defaultDailyTarget   = 8 * time.Hour
)

type config struct {
	dailyTarget time.Duration
	debug       bool
	verbose     bool
}

func readConfig(reader *twext.Reader) (config, error) {
	rawCfg, err := reader.ReadConfig()
	if err != nil {
		return config{}, fmt.Errorf("read config section: %w", err)
	}

	dailyTarget, err := configDailyTarget(rawCfg)
	if err != nil {
		return config{}, fmt.Errorf("get daily target: %w", err)
	}

	cfg := config{
		dailyTarget: dailyTarget,
		debug:       rawCfg[twext.ConfigKeyDebug].Bool(),
		verbose:     rawCfg[twext.ConfigKeyVerbose].Bool(),
	}

	return cfg, nil
}

func configDailyTarget(twConfig twext.Config) (time.Duration, error) {
	key := twext.NewConfigKey(configKeyPrefix, configKeyHoursPerDay)

	cfgValue, exists := twConfig[key]
	if !exists {
		return defaultDailyTarget, nil
	}

	intValue, err := cfgValue.Int()
	if err != nil {
		return 0, fmt.Errorf("convert to int: %w", err)
	}

	dailyTarget := time.Duration(intValue) * time.Hour

	return dailyTarget, nil
}
