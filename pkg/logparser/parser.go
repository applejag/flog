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
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/jilleJr/flog/pkg/loglevel"
	"gopkg.in/guregu/null.v3"
)

type ResultType byte

const (
	ResultNoMatch ResultType = iota
	ResultMatch
	ResultMatchMayContinue
)

type Parser interface {
	Parse(line string) (ParsedLog, ResultType)
}

func ParseUsingAnyParser(line string) ParsedLog {
	for _, parser := range defaultParsers {
		log, result := parser.Parse(line)
		if result == ResultMatch || result == ResultMatchMayContinue {
			return log
		}
	}
	return ParsedLog{String: line, Level: loglevel.Undefined}
}

type RegExParser struct {
	Expression     *regexp.Regexp
	TimeLayout     string
	GroupTimestamp int
	GroupLevel     int
}

func (p RegExParser) Parse(line string) (ParsedLog, ResultType) {
	stripped := stripansi.Strip(line)
	matches := p.Expression.FindStringSubmatch(stripped)
	matchesLen := len(matches)
	if matchesLen == 0 {
		return ParsedLog{}, ResultNoMatch
	}
	log := ParsedLog{
		String: line,
	}
	if p.GroupTimestamp > 0 && p.GroupTimestamp < matchesLen {
		log.Timestamp = parseTime(matches[p.GroupTimestamp], p.TimeLayout)
	}
	if p.GroupLevel > 0 && p.GroupLevel < matchesLen {
		log.Level = loglevel.ParseLevel(matches[p.GroupLevel])
	}
	return log, ResultMatch
}

type JSONParser struct{}

func (p JSONParser) Parse(line string) (ParsedLog, ResultType) {
	var obj map[string]interface{}
	if json.Unmarshal([]byte(line), &obj) != nil {
		return ParsedLog{}, ResultNoMatch
	}
	log := ParsedLog{
		Timestamp: parseTime(readJSONTimestamp(obj), ""),
		Level:     loglevel.ParseLevel(readJSONLevel(obj)),
		String:    line,
	}
	return log, ResultMatch
}

func readJSONLevel(obj map[string]interface{}) string {
	var msg string
	var ok bool
	if msg, ok = tryMapValueString(obj, "level"); ok {
		return msg
	} else if msg, ok = tryMapValueString(obj, "lvl"); ok {
		return msg
	} else if msg, ok = tryMapValueString(obj, "severity"); ok {
		return msg
	}
	return ""
}

func readJSONTimestamp(obj map[string]interface{}) string {
	var msg string
	var ok bool
	if msg, ok = tryMapValueString(obj, "timestamp"); ok {
		return msg
	} else if msg, ok = tryMapValueString(obj, "date"); ok {
		return msg
	} else if msg, ok = tryMapValueString(obj, "datetime"); ok {
		return msg
	}
	return ""
}

func tryMapValueString(obj map[string]interface{}, key string) (string, bool) {
	if value, ok := obj[key]; ok {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

const dateTimeRegex = `\d{4}-\d\d?-\d\d?(?:[ ·T]\d\d?[:.]\d\d?(?:[:.]\d+(?:\.\d+)?)?(?:Z|[+-]?\d{2}:?\d{2})?)?`

func compileRegexp(value string) *regexp.Regexp {
	value = strings.ReplaceAll(value, "{date}", dateTimeRegex)
	return regexp.MustCompile(value)
}

var defaultParsers = []Parser{
	JSONParser{},
	RegExParser{
		// 2021-01-31 17:33:54.3326|TRACE|Program|Sample
		Expression:     compileRegexp(`^({date})\|(\w+)\|.*$`),
		GroupTimestamp: 1,
		GroupLevel:     2,
	},
	RegExParser{
		// time="2021-01-31T19:04:01+01:00" level=trace msg="A walrus appears" animal=walrus
		Expression:     compileRegexp(`^time="({date})"[\s·]level="?(\w+)["\s$].*$`),
		GroupTimestamp: 1,
		GroupLevel:     2,
	},
	RegExParser{
		// WARN[0000] A walrus appears            animal=walrus
		// fail: Program[0]
		Expression: compileRegexp(`^(\w{4})[\[:].*$`),
		GroupLevel: 1,
	},
	RegExParser{
		// I0204 09:00:44.662471       i health.go:55] Starting MySQL health checker...
		Expression:     compileRegexp(`^(\w)(\d{4} \d\d:\d\d:\d\d(\.?\d+)?)\s+.*$`),
		TimeLayout:     "0102 15:04:05.999999999",
		GroupTimestamp: 2,
		GroupLevel:     1,
	},
	RegExParser{
		// Jun-18 14:50+0200 [DEBUG | TEST | wharf-core/main.go:23] Sample  hello=world
		Expression:     compileRegexp(`^([a-zA-Z0-9:+ \-]+) \[(\w+).*$`),
		TimeLayout:     "Jan-02 15:04Z0700",
		GroupTimestamp: 1,
		GroupLevel:     2,
	},
}

var timeLayouts = []string{
	"2006-01-02 15:04:05.999999999Z07:00",
	"2006-01-02 15:04:05.999999999-0700",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02·15:04:05.999999999Z07:00",
	"2006-01-02·15:04:05.999999999-0700",
	"2006-01-02·15:04:05.999999999",
	"2006-01-02 15.04.05.999999999Z07:00",
	"2006-01-02 15.04.05.999999999-0700",
	"2006-01-02 15.04.05.999999999",
	"2006-01-02·15.04.05.999999999Z07:00",
	"2006-01-02·15.04.05.999999999-0700",
	"2006-01-02·15.04.05.999999999",
	time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"
	"2006-01-02T15:04:05.999999999-0700",
	"2006-01-02T15:04:05.999999999",

	"15:04:05.999999999",
	"15:04",
	time.Kitchen, // "3:04PM"

	time.ANSIC,    // "Mon Jan _2 15:04:05 2006"
	time.UnixDate, // "Mon Jan _2 15:04:05 MST 2006"
	time.RFC822,   // "02 Jan 06 15:04 MST"
	time.RFC822Z,  // "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	time.RFC850,   // "Monday, 02-Jan-06 15:04:05 MST"
	time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	time.RubyDate, // "Mon Jan 02 15:04:05 -0700 2006"
}

func parseTime(value, preferredLayout string) null.Time {
	if value == "" {
		return null.Time{}
	}
	if preferredLayout != "" {
		if t, err := time.Parse(preferredLayout, value); err == nil {
			return null.TimeFrom(timeDefaults(t))
		}
	}
	for _, layout := range timeLayouts {
		if t, err := time.Parse(layout, value); err == nil {
			return null.TimeFrom(timeDefaults(t))
		}
	}
	return null.Time{}
}

func timeDefaults(t time.Time) time.Time {
	if t.Year() == 0 {
		t = time.Date(time.Now().Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	}
	return t
}
