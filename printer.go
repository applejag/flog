package main

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
	parser        logparser.Parser
	filter        LogFilter
	levelsSkipped map[loglevel.Level]int
	skippedAny    bool
}

func NewConsolePrinter(name string, p logparser.Parser, filter LogFilter) Printer {
	return &consolePrinter{
		name:          name,
		parser:        p,
		filter:        filter,
		levelsSkipped: map[loglevel.Level]int{},
		skippedAny:    false,
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

func shouldIncludeLogInOutput(lvl loglevel.Level, filter LogFilter) bool {
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

	if loggingLevel > log.InfoLevel {
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
