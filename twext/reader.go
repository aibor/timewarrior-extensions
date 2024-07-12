// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: MIT

package twext

import (
	"bufio"
	"io"
)

// Reader parses timewarrior extension input data.
//
// The format is expected as described here:
// https://timewarrior.net/docs/api/#input-format
//
// After creating a [NewReader], call [Reader.ReadConfig] first and then
// [Reader.ReadEntries] to get the actual time data.
type Reader struct {
	reader *bufio.Reader

	configRead bool
}

// NewReader creates a new [Reader] object.
//
// It does not read any data from the given reader yet.
func NewReader(r io.ReadSeeker) *Reader {
	return &Reader{
		reader: bufio.NewReader(r),
	}
}

// ReadConfig reads the config section of the input.
//
// It must be called before calling [Reader.ReadEntries].
func (r *Reader) ReadConfig() (Config, error) {
	if r.configRead {
		return nil, ErrReaderConfigConsumed
	}

	r.configRead = true

	return readConfig(r.reader)
}

// ReadEntries reads the list of timewarrior entries.
//
// It returns [ErrReaderConfigNotConsumed] if the configuration section of the
// input data has not been read yet. Call [Reader.ReadConfig] beforehand.
func (r *Reader) ReadEntries() (Entries, error) {
	if !r.configRead {
		return nil, ErrReaderConfigNotConsumed
	}

	return readEntries(r.reader)
}
