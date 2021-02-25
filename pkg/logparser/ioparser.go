// Filter multiline logs based on the log's severity
// Copyright (C) 2021  Kalle Jillheden
//
// flog is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// flog is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
