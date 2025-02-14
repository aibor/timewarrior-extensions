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
	"time"

	"github.com/aibor/timewarrior-extensions/twext"
)

type daySums = twext.Groups[string, time.Duration]

func startDate(entry twext.Entry) string {
	return entry.Start.Format(time.DateOnly)
}

func sumDuration(result time.Duration, entry twext.Entry) time.Duration {
	if !entry.End.IsZero() && !entry.Start.SameDate(&entry.End) {
		log.Printf("entry %d spans multiple days. Skipping.", entry.ID)

		return result
	}

	return result + entry.Duration()
}

func printSums(p *printer, daySums daySums) error {
	_, err := fmt.Fprintln(p.writer)
	if err != nil {
		return fmt.Errorf("write header newline: %w", err)
	}

	err = p.write("date", "actual", "target", "diff")
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	total := p.cfg.offset

	if p.cfg.offset != 0 && p.cfg.verbose {
		err = p.writeTime("offset", p.cfg.offset, 0)
		if err != nil {
			return fmt.Errorf("write offset: %w", err)
		}
	}

	for day, daySum := range daySums.Sorted() {
		total += daySum

		if !p.cfg.verbose {
			continue
		}

		err := p.writeTime(day, daySum, p.cfg.dailyTarget)
		if err != nil {
			return fmt.Errorf("write day [%s]: %w", day, err)
		}
	}

	totalTarget := time.Duration(len(daySums)) * p.cfg.dailyTarget

	err = p.writeTime("total", total, totalTarget)
	if err != nil {
		return fmt.Errorf("write total: %w", err)
	}

	err = p.writer.Flush()
	if err != nil {
		return fmt.Errorf("flush: %w", err)
	}

	return nil
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
		log.Println("cfg - DailyTarget:", cfg.dailyTarget)
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
		twext.Group(entries, startDate, sumDuration),
	)
}

func main() {
	err := run(os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
