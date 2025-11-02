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

func (p *printer) printf(format string, args ...any) {
	_, err := fmt.Fprintf(p.writer, format, args...)
	if err != nil {
		panic(fmt.Errorf("fprintf: %w", err))
	}
}

func (p *printer) flush() {
	err := p.writer.Flush()
	if err != nil {
		panic(fmt.Errorf("flush: %w", err))
	}
}

func (p *printer) write(handle string, actual, target, diff string) {
	p.printf("%s\t%s\t%s\t%s\t\n", handle, actual, target, diff)
}

func (p *printer) writeTime(handle string, actual, target time.Duration) {
	p.write(
		handle,
		fmtDuration(actual),
		fmtDuration(target),
		fmtDuration(actual-target),
	)
}

func (p *printer) writeHeader() {
	p.printf("\n")
	p.write("date", "actual", "target", "diff")
}

func (p *printer) writeTotals(totalSum, totalTarget time.Duration) {
	p.writeTime("total", totalSum, totalTarget)
}
