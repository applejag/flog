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

package loglevel

import "strings"

type Level int

const (
	Undefined Level = 0
	Unknown   Level = 1 << (iota - 1)
	Trace
	Debug
	Information
	Warning
	Error
	Critical
	Fatal
	Panic
)

var singularLevels = []Level{
	Unknown,
	Trace,
	Debug,
	Information,
	Warning,
	Error,
	Critical,
	Fatal,
	Panic,
}

func (lvl Level) String() string {
	return lvl.StringDelim('|')
}

func (lvl Level) StringDelim(delim rune) string {
	var b = strings.Builder{}
	for _, singularLevel := range singularLevels {
		if lvl&singularLevel != Undefined {
			if b.Len() > 0 {
				b.WriteRune(delim)
			}
			b.WriteString(singularLevelString(singularLevel))
		}
	}
	if b.Len() == 0 {
		return "Undefined"
	}
	return b.String()
}

func (lvl Level) Levels() []Level {
	var slice []Level
	for _, singularLevel := range singularLevels {
		if lvl&singularLevel != Undefined {
			slice = append(slice, singularLevel)
		}
	}
	return slice
}

func singularLevelString(lvl Level) string {
	switch lvl {
	case Unknown:
		return "Unknown"
	case Trace:
		return "Trace"
	case Debug:
		return "Debug"
	case Information:
		return "Information"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	case Critical:
		return "Critical"
	case Fatal:
		return "Fatal"
	case Panic:
		return "Panic"
	}
	return "Undefined"
}

func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "0", "t", "trc", "trce", "trac", "trace":
		return Trace

	case "1", "d", "dbg", "debu", "dbug", "debg", "debug":
		return Debug

	case "2", "i", "inf", "info", "information":
		return Information

	case "3", "w", "warn", "warning":
		return Warning

	case "4", "e", "err", "erro", "error", "fail":
		return Error

	case "5", "c", "crit", "critical":
		return Critical

	case "6", "f", "fata", "fatal":
		return Fatal

	case "7", "p", "panic":
		return Panic
	}

	return Unknown
}
