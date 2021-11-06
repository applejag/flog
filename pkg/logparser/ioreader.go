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

package logparser

import (
	"bufio"
	"io"

	"github.com/jilleJr/flog/pkg/loglevel"
)

type IOReader struct {
	scanner *bufio.Scanner
	lastLog ParsedLog
}

func NewIOReader(r io.Reader) IOReader {
	return IOReader{
		scanner: bufio.NewScanner(r),
	}
}

func (p *IOReader) ParsedLog() ParsedLog {
	return p.lastLog
}

func (p *IOReader) Scan() bool {
	if !p.scanner.Scan() {
		return false
	}
	lastLevel := p.lastLog.Level
	p.lastLog = ParseUsingAnyParser(p.scanner.Text())
	if p.lastLog.Level == loglevel.Undefined || p.lastLog.Level == loglevel.Unknown {
		p.lastLog.Level = lastLevel
	}
	return true
}
