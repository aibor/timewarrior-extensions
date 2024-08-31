// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

const minutesPerHour = 60

func fmtDuration(d time.Duration) string {
	var prefix string

	if d < 0 {
		prefix = "-"
	}

	return fmt.Sprintf(
		"%s%dh:%02dm",
		prefix,
		int64(d.Abs().Hours()),
		int64(d.Abs().Minutes())%minutesPerHour,
	)
}

const (
	tabPadding = 4
	tabFlags   = tabwriter.AlignRight
)

type printer struct {
	writer *tabwriter.Writer
	cfg    config
}

func newPrinter(w io.Writer, cfg config) *printer {
	return &printer{
		writer: tabwriter.NewWriter(w, 0, 0, tabPadding, ' ', tabFlags),
		cfg:    cfg,
	}
}

func (p *printer) write(handle string, actual, target, diff string) error {
	_, err := fmt.Fprintf(
		p.writer,
		"%s\t%s\t%s\t%s\t\n",
		handle,
		actual,
		target,
		diff,
	)
	if err != nil {
		return fmt.Errorf("fprintf: %w", err)
	}

	return nil
}

func (p *printer) writeTime(handle string, actual, target time.Duration) error {
	diff := actual - target

	return p.write(
		handle,
		fmtDuration(actual),
		fmtDuration(target),
		fmtDuration(diff),
	)
}
