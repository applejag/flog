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
}

func NewConsolePrinter(p logparser.Parser) Printer {
	return consolePrinter{
		parser: p,
	}
}

func (p consolePrinter) Next() bool {
	if !p.parser.Scan() {
		return false
	}
	log := p.parser.ParsedLog()
	if log.Level > logparser.LevelInformation {
		fmt.Println(log.String)
	}
	return true
}
