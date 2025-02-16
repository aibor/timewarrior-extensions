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
	configKeyOffsetTotal = "minutes_offset_total"
	defaultOffsetTotal   = 0
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
		configKeyHoursPerDay,
		defaultDailyTarget,
		time.Hour,
	)
	if err != nil {
		return config{}, fmt.Errorf("get daily target: %w", err)
	}

	offset, err := configReadDuration(
		rawCfg,
		configKeyOffsetTotal,
		defaultOffsetTotal,
		time.Minute,
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
	defValue time.Duration,
	unit time.Duration,
) (time.Duration, error) {
	key := twext.NewConfigKey(configKeyPrefix, cfgKey)

	cfgValue, exists := twConfig[key]
	if !exists {
		return defValue, nil
	}

	intValue, err := cfgValue.Int()
	if err != nil {
		return 0, fmt.Errorf("convert to int: %w", err)
	}

	duration := time.Duration(intValue) * unit

	return duration, nil
}
