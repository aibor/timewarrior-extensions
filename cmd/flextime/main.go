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

		dayTarget, err := p.cfg.target.targetFor(date.Weekday())
		if err != nil {
			return fmt.Errorf("get weekday target: %w", err)
		}

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
		log.Println("cfg - Target - Monday:", cfg.target.monday)
		log.Println("cfg - Target - Tuesday:", cfg.target.tuesday)
		log.Println("cfg - Target - Wednesday:", cfg.target.wednesday)
		log.Println("cfg - Target - Thursday:", cfg.target.thursday)
		log.Println("cfg - Target - Friday:", cfg.target.friday)
		log.Println("cfg - Target - Saturday:", cfg.target.saturday)
		log.Println("cfg - Target - Sunday:", cfg.target.sunday)
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
