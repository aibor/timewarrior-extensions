// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package twext

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	numConfigParts     = 2
	configKeySeparator = "."
)

// ConfigKey is an e complete config key string.
type ConfigKey string

// Well known ConfigKeys.
const (
	ConfigKeyVerbose      ConfigKey = "verbose"
	ConfigKeyDebug        ConfigKey = "debug"
	ConfigKeyConfirmation ConfigKey = "confirmation"
)

// NewConfigKey composes a new [ConfigKey].
func NewConfigKey(s ...string) ConfigKey {
	return ConfigKey(strings.Join(s, configKeySeparator))
}

func (k ConfigKey) String() string {
	return string(k)
}

// ConfigValue is a config value string.
type ConfigValue string

func (v ConfigValue) String() string {
	return string(v)
}

// Bool returns true if the [ConfigValue] matches on of the defined values
// indication trueness.
func (v ConfigValue) Bool() bool {
	// True boolean values as described in
	// https://timewarrior.net/docs/api/#guidelines
	boolValues := []string{"on", "1", "yes", "y", "true"}

	return slices.Contains(boolValues, string(v))
}

// Int tries to parse the [ConfigValue] as integer. It returns an error if
// the string can not be parsed as integer.
func (v ConfigValue) Int() (int, error) {
	i, err := strconv.Atoi(v.String())
	if err != nil {
		return 0, fmt.Errorf("atoi: %w", err)
	}

	return i, nil
}

// Duration tries to parse the [ConfigValue] as [time.Duration]. It returns an
// error if the string can not be parsed as [time.Duration]. See
// [time.ParseDuration] for the supported format.
func (v ConfigValue) Duration() (time.Duration, error) {
	d, err := time.ParseDuration(v.String())
	if err != nil {
		return 0, fmt.Errorf("parse: %w", err)
	}

	return d, nil
}

// Config is a collection of configuration directives.
type Config map[ConfigKey]ConfigValue

type stringReader interface {
	ReadString(delimiter byte) (string, error)
}

func readConfig(reader stringReader) (Config, error) {
	config := make(Config)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("read line: %w", err)
		}

		line = strings.TrimRight(line, "\n")
		if line == "" {
			break
		}

		key, value, err := readConfigLine(line)
		if err != nil {
			return nil, fmt.Errorf("read config line [%s]: %w", line, err)
		}

		config[key] = value
	}

	if len(config) < 1 {
		return nil, ErrConfigEmpty
	}

	return config, nil
}

func readConfigLine(line string) (ConfigKey, ConfigValue, error) {
	configLine := strings.Split(line, ": ")
	if len(configLine) != numConfigParts {
		return "", "", ErrConfigInvalidLine
	}

	return ConfigKey(configLine[0]), ConfigValue(configLine[1]), nil
}
