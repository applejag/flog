package main

import (
	"fmt"

	"github.com/jilleJr/flog/pkg/loglevel"
	"github.com/jilleJr/flog/pkg/logparser"
)

type Printer interface {
	Next() bool
}

type consolePrinter struct {
	parser        logparser.Parser
	level         loglevel.Level
	levelsSkipped map[loglevel.Level]int
	skippedAny    bool
}

func NewConsolePrinter(p logparser.Parser, lvl loglevel.Level) Printer {
	return consolePrinter{
		parser:        p,
		level:         lvl,
		levelsSkipped: map[loglevel.Level]int{},
		skippedAny:    false,
	}
}

func (p consolePrinter) Next() bool {
	if !p.parser.Scan() {
		return false
	}
	log := p.parser.ParsedLog()
	if log.Level >= p.level {
		fmt.Println(log.String)
		if p.skippedAny {
			printSkippedLogs(p.levelsSkipped)
			p.levelsSkipped = map[loglevel.Level]int{}
		}
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

const (
	resetAnsi   = "\033[0m"
	skippedAnsi = "\033[90m"
)

func printSkippedLogs(skipped map[loglevel.Level]int) {
	fmt.Println("foo")
}

func getSkippedLevelsSlice(skipepd map[loglevel.Level]int) []string {
	return []string{}
}
