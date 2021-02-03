package loglevel

import "strings"

type Level int

const (
	LevelUndefined Level = iota
	Trace
	Debug
	Information
	Warning
	Error
	Critical
	Fatal
	Panic
)

func (lvl Level) String() string {
	switch lvl {
	case Trace:
		return "LevelTrace"
	case Debug:
		return "LevelDebug"
	case Information:
		return "LevelInformation"
	case Warning:
		return "LevelWarning"
	case Error:
		return "LevelError"
	case Critical:
		return "LevelCritical"
	case Fatal:
		return "LevelFatal"
	case Panic:
		return "LevelPanic"
	}
	return "LevelUndefined"
}

func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "trc", "trce", "trac", "trace":
		return Trace

	case "dbg", "debu", "dbug", "debg", "debug":
		return Debug

	case "inf", "info", "information":
		return Information

	case "warn", "warning":
		return Warning

	case "err", "erro", "error", "fail":
		return Error

	case "crit", "critical":
		return Critical

	case "fata", "fatal":
		return Fatal

	case "panic":
		return Panic
	}

	return LevelUndefined
}
