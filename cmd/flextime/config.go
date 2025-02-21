// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/aibor/timewarrior-extensions/twext"
)

const (
	configKeyPrefix          = "flextime"
	configKeyDailyTarget     = "time_per_day"
	defaultDailyTarget       = "8h"
	configKeyMondayTarget    = "time_per_monday"
	configKeyTuesdayTarget   = "time_per_tuesday"
	configKeyWednesdayTarget = "time_per_wednesday"
	configKeyThursdayTarget  = "time_per_thursday"
	configKeyFridayTarget    = "time_per_friday"
	configKeySaturdayTarget  = "time_per_saturday"
	configKeySundayTarget    = "time_per_sunday"
	defaultWeekdayTarget     = "-1s"
	configKeyOffsetTotal     = "offset_total"
	defaultOffsetTotal       = "0"
)

type weekdayTargets struct {
	monday    time.Duration
	tuesday   time.Duration
	wednesday time.Duration
	thursday  time.Duration
	friday    time.Duration
	saturday  time.Duration
	sunday    time.Duration
}

var errInvalidWeekday = errors.New("unknown weekday")

func (t weekdayTargets) targetFor(weekday time.Weekday) (time.Duration, error) {
	switch weekday {
	case time.Monday:
		return t.monday, nil
	case time.Tuesday:
		return t.tuesday, nil
	case time.Wednesday:
		return t.wednesday, nil
	case time.Thursday:
		return t.thursday, nil
	case time.Friday:
		return t.friday, nil
	case time.Saturday:
		return t.saturday, nil
	case time.Sunday:
		return t.sunday, nil
	default:
		return 0, fmt.Errorf("%w: %d", errInvalidWeekday, weekday)
	}
}

type config struct {
	target  weekdayTargets
	offset  time.Duration
	debug   bool
	verbose bool
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

	mondayTarget, err := configReadDuration(
		rawCfg,
		configKeyMondayTarget,
		defaultWeekdayTarget,
	)
	if err != nil {
		return config{}, fmt.Errorf("get daily target: %w", err)
	}
	if mondayTarget < 0 {
		mondayTarget = dailyTarget
	}

	tuesdayTarget, err := configReadDuration(
		rawCfg,
		configKeyTuesdayTarget,
		defaultWeekdayTarget,
	)
	if err != nil {
		return config{}, fmt.Errorf("get daily target: %w", err)
	}
	if tuesdayTarget < 0 {
		tuesdayTarget = dailyTarget
	}

	wednesdayTarget, err := configReadDuration(
		rawCfg,
		configKeyWednesdayTarget,
		defaultWeekdayTarget,
	)
	if err != nil {
		return config{}, fmt.Errorf("get daily target: %w", err)
	}
	if wednesdayTarget < 0 {
		wednesdayTarget = dailyTarget
	}

	thursdayTarget, err := configReadDuration(
		rawCfg,
		configKeyThursdayTarget,
		defaultWeekdayTarget,
	)
	if err != nil {
		return config{}, fmt.Errorf("get daily target: %w", err)
	}
	if thursdayTarget < 0 {
		thursdayTarget = dailyTarget
	}

	fridayTarget, err := configReadDuration(
		rawCfg,
		configKeyFridayTarget,
		defaultWeekdayTarget,
	)
	if err != nil {
		return config{}, fmt.Errorf("get daily target: %w", err)
	}
	if fridayTarget < 0 {
		fridayTarget = dailyTarget
	}

	saturdayTarget, err := configReadDuration(
		rawCfg,
		configKeySaturdayTarget,
		defaultWeekdayTarget,
	)
	if err != nil {
		return config{}, fmt.Errorf("get daily target: %w", err)
	}
	if saturdayTarget < 0 {
		saturdayTarget = dailyTarget
	}

	sundayTarget, err := configReadDuration(
		rawCfg,
		configKeySundayTarget,
		defaultWeekdayTarget,
	)
	if err != nil {
		return config{}, fmt.Errorf("get daily target: %w", err)
	}
	if sundayTarget < 0 {
		sundayTarget = dailyTarget
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
		target: weekdayTargets{
			monday:    mondayTarget,
			tuesday:   tuesdayTarget,
			wednesday: wednesdayTarget,
			thursday:  thursdayTarget,
			friday:    fridayTarget,
			saturday:  saturdayTarget,
			sunday:    sundayTarget,
		},
		offset:  offset,
		debug:   rawCfg[twext.ConfigKeyDebug].Bool(),
		verbose: rawCfg[twext.ConfigKeyVerbose].Bool(),
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
