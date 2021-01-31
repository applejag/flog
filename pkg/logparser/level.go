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

func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "trc":
	case "trce":
	case "trac":
	case "trace":
		return LevelTrace

	case "dbg":
	case "dbug":
	case "debg":
	case "debug":
		return LevelDebug

	case "inf":
	case "info":
	case "information":
		return LevelInformation

	case "warn":
	case "warning":
		return LevelWarning

	case "fail":
		return LevelFail

	case "err":
	case "error":
		return LevelError

	case "crit":
	case "critical":
		return LevelCritical

	case "fatal":
		return LevelFatal

	case "panic":
		return LevelPanic
	}

	return LevelUnknown
}
