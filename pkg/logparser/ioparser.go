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
	txt := p.scanner.Text()
	p.lastLog = ParsedLog{
		String: txt,
		Level:  LevelError,
	}
	return true
}
