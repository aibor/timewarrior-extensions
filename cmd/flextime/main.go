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

//nolint:cyclop
func printSums(p *printer, daySums daySums) error {
	_, err := fmt.Fprintln(p.writer)
	if err != nil {
		return fmt.Errorf("write header newline: %w", err)
	}

	err = p.write("date", "actual", "target", "diff")
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	totalSum := p.cfg.offset

	var totalTarget time.Duration

	if p.cfg.offset != 0 && p.cfg.verbose {
		err := p.writeTime("offset", p.cfg.offset, 0)
		if err != nil {
			return fmt.Errorf("write offset: %w", err)
		}
	}

	for day, daySum := range daySums.Sorted() {
		date, err := time.Parse(time.DateOnly, day)
		if err != nil {
			return fmt.Errorf("parse date: %w", err)
		}

		dayTarget := p.cfg.timeTargets.targetFor(date.Weekday())

		totalSum += daySum
		totalTarget += dayTarget

		if !p.cfg.verbose {
			continue
		}

		err = p.writeTime(day, daySum, dayTarget)
		if err != nil {
			return fmt.Errorf("write day [%s]: %w", day, err)
		}
	}

	err = p.writeTime("total", totalSum, totalTarget)
	if err != nil {
		return fmt.Errorf("write totalSum: %w", err)
	}

	err = p.writer.Flush()
	if err != nil {
		return fmt.Errorf("flush: %w", err)
	}

	return nil
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

	return printSums(
		newPrinter(outW, cfg),
		cfg.aggregationStrategy.Aggregate(entries.All()),
	)
}

func main() {
	err := run(os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
