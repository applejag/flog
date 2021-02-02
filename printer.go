package main

import (
	"fmt"

	"github.com/jilleJr/flog/pkg/logparser"
)

type Printer interface {
	Next() bool
}

type consolePrinter struct {
	parser logparser.Parser
	level  logparser.Level
}

func NewConsolePrinter(p logparser.Parser, lvl logparser.Level) Printer {
	return consolePrinter{
		parser: p,
		level:  lvl,
	}
}

func (p consolePrinter) Next() bool {
	if !p.parser.Scan() {
		return false
	}
	log := p.parser.ParsedLog()
	if log.Level >= p.level {
		fmt.Println(log.String)
	}
	return true
}
