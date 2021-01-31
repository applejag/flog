package logparser

import "strings"

type Level int

const (
	LevelUnknown Level = iota
	LevelTrace
	LevelDebug
	LevelInformation
	LevelWarning
	LevelFail
	LevelError
	LevelCritical
	LevelFatal
	LevelPanic
)

func (lvl Level) String() string {
	switch lvl {
	case LevelTrace:
		return "LevelTrace"
	case LevelDebug:
		return "LevelDebug"
	case LevelInformation:
		return "LevelInformation"
	case LevelWarning:
		return "LevelWarning"
	case LevelFail:
		return "LevelFail"
	case LevelError:
		return "LevelError"
	case LevelCritical:
		return "LevelCritical"
	case LevelFatal:
		return "LevelFatal"
	case LevelPanic:
		return "LevelPanic"
	}
	return "LevelUnknown"
}

func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "trc", "trce", "trac", "trace":
		return LevelTrace

	case "dbg", "debu", "dbug", "debg", "debug":
		return LevelDebug

	case "inf", "info", "information":
		return LevelInformation

	case "warn", "warning":
		return LevelWarning

	case "fail":
		return LevelFail

	case "err", "erro", "error":
		return LevelError

	case "crit", "critical":
		return LevelCritical

	case "fata", "fatal":
		return LevelFatal

	case "panic":
		return LevelPanic
	}

	return LevelUnknown
}
