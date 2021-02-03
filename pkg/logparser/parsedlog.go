package logparser

import (
	"time"

	"github.com/jilleJr/flog/pkg/loglevel"
)

type ParsedLog struct {
	Level     loglevel.Level
	String    string
	Timestamp time.Time
}
