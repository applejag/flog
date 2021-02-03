package loglevel

import "strings"

type Level int

const (
	Undefined Level = iota
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

	return Undefined
}
