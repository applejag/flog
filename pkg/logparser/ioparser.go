package logparser

import (
	"bufio"
	"io"

	"github.com/jilleJr/flog/pkg/loglevel"
)

type IOParser struct {
	scanner *bufio.Scanner
	lastLog ParsedLog
}

func NewIOParser(r io.Reader) IOParser {
	return IOParser{
		scanner: bufio.NewScanner(r),
	}
}

func (p *IOParser) ParsedLog() ParsedLog {
	return p.lastLog
}

func (p *IOParser) Scan() bool {
	if !p.scanner.Scan() {
		return false
	}
	lastLevel := p.lastLog.Level
	p.lastLog = parseLog(p.scanner.Text())
	if p.lastLog.Level == loglevel.Undefined || p.lastLog.Level == loglevel.Unknown {
		p.lastLog.Level = lastLevel
	}
	return true
}
