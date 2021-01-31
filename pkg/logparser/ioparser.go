package logparser

import (
	"bufio"
	"io"
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
	if p.lastLog.Level == LevelUnknown {
		p.lastLog.Level = lastLevel
	}
	return true
}
