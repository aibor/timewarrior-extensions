// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: MIT

package twext

import "errors"

var (
	// ErrDateUnmarshalNotString is returned in case the value given to the
	// JSON unmarshaller is not a string.
	ErrDateUnmarshalNotString = errors.New("input is not a JSON string")

	// ErrConfigInvalidLine is returned if the line can not be split into a
	// key and value part.
	ErrConfigInvalidLine = errors.New("config line has invalid format")

	// ErrConfigEmpty is returned if the config section is empty.
	ErrConfigEmpty = errors.New("config is empty")

	// ErrReaderConfigConsumed is returned in case the config section is
	// tried to be read again.
	ErrReaderConfigConsumed = errors.New("config section already consumed")

	// ErrReaderConfigNotConsumed is returned if the entries are tried to be
	// read before the config section has been read.
	ErrReaderConfigNotConsumed = errors.New("config section not read yet")
)
