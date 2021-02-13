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
	"regexp"
	"strings"
	"unicode"

	"github.com/acarl005/stripansi"
	"github.com/jilleJr/flog/pkg/loglevel"
)

type Parser interface {
	Scan() bool
	ParsedLog() ParsedLog
}

var datetimeRegex = regexp.MustCompile(`\d{4}-\d\d?-\d\d?(?:[ ·T]\d\d?[:.]\d\d?(?:[:.]\d+(?:\.\d+)?)?(Z|[-+ ]\d\d?([:.]\d{1,3})?)?)?`)
var levelRegex = regexp.MustCompile(`\(\w+\)|\[\w+\]|'\w+'|"\w+"|\w+[\[(|:]|=\w+`)

func parseLog(s string) ParsedLog {
	if len(s) > 0 && unicode.IsSpace(rune(s[0])) {
		return ParsedLog{
			String: s,
			Level:  loglevel.Undefined,
		}
	}

	stripped := stripansi.Strip(s)
	if loc := datetimeRegex.FindStringIndex(stripped); loc != nil {
		stripped = stripped[loc[1]:]
	}

	var level loglevel.Level
	if lvls := levelRegex.FindAllString(stripped, 5); lvls != nil {
		for _, lvlStr := range lvls {
			lvlStr = strings.Trim(lvlStr, "|[]():=\"'")
			if lvl := loglevel.ParseLevel(lvlStr); lvl != loglevel.Unknown {
				level = lvl
				break
			}
		}
	}

	return ParsedLog{
		String: s,
		Level:  level,
	}
}
