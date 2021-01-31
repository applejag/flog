package logparser

import "time"

type ParsedLog struct {
	Level     Level
	String    string
	Timestamp time.Time
}
