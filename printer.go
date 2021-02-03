package main

import (
	"fmt"
	"strings"

	"github.com/jilleJr/flog/pkg/loglevel"
	"github.com/jilleJr/flog/pkg/logparser"
)

type Printer interface {
	Next() bool
	PrintOmittedLogs()
}

type consolePrinter struct {
	parser        logparser.Parser
	filter        LogFilter
	levelsSkipped map[loglevel.Level]int
	skippedAny    bool
}

func NewConsolePrinter(p logparser.Parser, filter LogFilter) Printer {
	return &consolePrinter{
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
	log := p.parser.ParsedLog()
	if shouldIncludeLogInOutput(log.Level, p.filter) {
		if p.skippedAny {
			if p.filter.Quiet == false {
				printSkippedLogs(p.levelsSkipped)
			}
			p.levelsSkipped = map[loglevel.Level]int{}
			p.skippedAny = false
		}
		fmt.Println(log.String)
	} else {
		p.skippedAny = true
		if i, ok := p.levelsSkipped[log.Level]; ok {
			p.levelsSkipped[log.Level] = i + 1
		} else {
			p.levelsSkipped[log.Level] = 1
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
	if p.skippedAny && p.filter.Quiet == false {
		printSkippedLogs(p.levelsSkipped)
	}
}

const (
	resetAnsi   = "\033[0m"
	skippedAnsi = "\033[90m\033[3m" // gray and italic
)

func printSkippedLogs(skipped map[loglevel.Level]int) {
	skippedStrings := getSkippedLevelsSlice(skipped)
	str := strings.Join(skippedStrings, ", ")
	fmt.Print(skippedAnsi)
	fmt.Print("flog: Omitted ")
	fmt.Print(str)
	fmt.Print(".")
	fmt.Print(resetAnsi)
	fmt.Println()
}

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

func getSkippedLevelsSlice(skipped map[loglevel.Level]int) []string {
	levels := make([]string, printableLevelsLen)
	index := 0
	for _, lvl := range printableLevels {
		if num, ok := skipped[lvl]; ok {
			levels[index] = fmt.Sprintf("%d %s", num, lvl.String())
			index++
		}
	}
	return levels[0:index]
}
