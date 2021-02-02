package logparser

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/acarl005/stripansi"
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
			Level:  LevelUndefined,
		}
	}

	stripped := stripansi.Strip(s)
	if loc := datetimeRegex.FindStringIndex(stripped); loc != nil {
		stripped = stripped[loc[1]:]
	}

	var level Level
	if lvls := levelRegex.FindAllString(stripped, 5); lvls != nil {
		for _, lvlStr := range lvls {
			lvlStr = strings.Trim(lvlStr, "|[]():=\"'")
			if lvl := ParseLevel(lvlStr); lvl != LevelUndefined {
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
