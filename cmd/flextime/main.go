// SPDX-FileCopyrightText: 2024 Tobias Böhm <code@aibor.de>
//
// SPDX-License-Identifier: GPL-3.0-or-later

// Implements a timewarrior extension for calculating daily sums and the
// difference to the daily target.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/aibor/timewarrior-extensions/twext"
)

type daySums = twext.Groups[string, time.Duration]

func printSums(p *printer, daySums daySums) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			var ok bool
			if err, ok = rec.(error); !ok {
				//nolint:err113
				err = fmt.Errorf("non-error panic: %v", rec)
			}
		}
	}()

	p.writeHeader()

	var totalTarget time.Duration

	totalSum := p.cfg.offset

	if p.cfg.offset != 0 && p.cfg.verbose {
		p.writeTime("offset", p.cfg.offset, 0)
	}

	for day, daySum := range daySums.Sorted() {
		date, err := time.Parse(time.DateOnly, day)
		if err != nil {
			return fmt.Errorf("parse date: %w", err)
		}

		dayTarget := p.cfg.timeTargets.targetFor(date)

		totalSum += daySum
		totalTarget += dayTarget

		if !p.cfg.verbose {
			continue
		}

		p.writeTime(day, daySum, dayTarget)
	}

	p.writeTotals(totalSum, totalTarget)
	p.flush()

	return err
}

func version() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	return buildInfo.Main.Version
}

func run(inR io.ReadSeeker, outW, errW io.Writer) error {
	reader := twext.NewReader(inR)

	cfg, err := readConfig(reader)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	log.SetPrefix("debug [flextime] - ")
	log.SetFlags(0)

	if cfg.debug {
		log.SetOutput(errW)
		log.Println("version:", version())
		log.Println("cfg - Offset:", cfg.offset)
		log.Println("cfg - Target:", cfg.timeTargets)
		log.Println("cfg - AggregationStrategy:", cfg.aggregationStrategy)
		log.Println("cfg - Debug:", cfg.debug)
		log.Println("cfg - Verbose:", cfg.verbose)
	} else {
		log.SetOutput(io.Discard)
	}

	entries, err := reader.ReadEntries()
	if err != nil {
		return fmt.Errorf("read entries: %w", err)
	}

	printer := newPrinter(outW, cfg)
	daySums := cfg.aggregationStrategy.Aggregate(entries.All())

	if err := printSums(printer, daySums); err != nil {
		return fmt.Errorf("print day sums: %w", err)
	}

	return nil
}

func main() {
	if err := run(os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
