package logparser

import "strings"

type Level int

const (
	LevelUndefined Level = iota
	LevelTrace
	LevelDebug
	LevelInformation
	LevelWarning
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
	case LevelError:
		return "LevelError"
	case LevelCritical:
		return "LevelCritical"
	case LevelFatal:
		return "LevelFatal"
	case LevelPanic:
		return "LevelPanic"
	}
	return "LevelUndefined"
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

	case "err", "erro", "error", "fail":
		return LevelError

	case "crit", "critical":
		return LevelCritical

	case "fata", "fatal":
		return LevelFatal

	case "panic":
		return LevelPanic
	}

	return LevelUndefined
}
