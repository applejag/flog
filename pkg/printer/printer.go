// Copyright (C) 2021  Kalle Jillheden
// SPDX-FileCopyrightText: 2021 Kalle Fagerberg
//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package printer

import (
	"fmt"

	"github.com/apex/log"
	"github.com/jilleJr/flog/pkg/loglevel"
	"github.com/jilleJr/flog/pkg/logparser"
)

type Printer interface {
	Next() bool
	PrintOmittedLogs()
}

type consolePrinter struct {
	name          string
	parser        logparser.Reader
	filter        loglevel.Filter
	levelsSkipped map[loglevel.Level]int
	skippedAny    bool
	loggingLevel  log.Level
}

func NewConsolePrinter(name string, p logparser.Reader, filter loglevel.Filter, loggingLevel log.Level) Printer {
	return &consolePrinter{
		name:          name,
		parser:        p,
		filter:        filter,
		levelsSkipped: map[loglevel.Level]int{},
		skippedAny:    false,
		loggingLevel:  loggingLevel,
	}
}

func (p *consolePrinter) Next() bool {
	if !p.parser.Scan() {
		return false
	}
	parsed := p.parser.ParsedLog()
	log.WithFields(log.Fields{
		"message": parsed.String,
		"level":   parsed.Level,
	}).Debugf("Parsed log from: %s", p.name)

	if shouldIncludeLogInOutput(parsed.Level, p.filter) {
		if p.skippedAny {
			p.PrintOmittedLogs()
			p.levelsSkipped = map[loglevel.Level]int{}
			p.skippedAny = false
		}
		fmt.Println(parsed.String)
	} else {
		p.skippedAny = true
		if i, ok := p.levelsSkipped[parsed.Level]; ok {
			p.levelsSkipped[parsed.Level] = i + 1
		} else {
			p.levelsSkipped[parsed.Level] = 1
		}
	}
	return true
}

func shouldIncludeLogInOutput(lvl loglevel.Level, filter loglevel.Filter) bool {
	if filter.WhitelistMask != loglevel.Undefined && filter.WhitelistMask&lvl == loglevel.Undefined {
		return false
	}

	if lvl != loglevel.Unknown && lvl != loglevel.Undefined {
		if filter.MinLevel != loglevel.Undefined && lvl < filter.MinLevel {
			return false
		}

		if filter.MaxLevel != loglevel.Undefined && lvl > filter.MaxLevel {
			return false
		}
	} else if filter.BlacklistMask&loglevel.Unknown > 0 {
		return false
	}

	if filter.BlacklistMask&lvl != loglevel.Undefined {
		return false
	}

	return true
}

func (p *consolePrinter) PrintOmittedLogs() {
	if !p.skippedAny {
		return
	}

	if p.loggingLevel > log.InfoLevel {
		return
	}

	fields := getSkippedLevelsFields(p.levelsSkipped)
	log.WithFields(fields).Infof("Omitted logs from: %s", p.name)
}

const (
	resetAnsi   = "\033[0m"
	skippedAnsi = "\033[90m\033[3m" // gray and italic
)

const printableLevelsLen = 8

var printableLevels = []loglevel.Level{
	loglevel.Trace,
	loglevel.Debug,
	loglevel.Information,
	loglevel.Warning,
	loglevel.Error,
	loglevel.Critical,
	loglevel.Fatal,
	loglevel.Panic,
}

func getSkippedLevelsFields(skipped map[loglevel.Level]int) log.Fields {
	fields := make(log.Fields, len(skipped))
	for lvl, count := range skipped {
		fields[lvl.String()] = count
	}
	return fields
}
