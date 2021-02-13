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

// Package console implements a development-friendly textual handler for apex.
//
// This was originally forked from
// https://github.com/apex/log/blob/f0aad53/handlers/text/text.go
package console

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/apex/log"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr, "")

// start time.
var start = time.Now()

// colors.
const (
	italic        = "\033[3m"
	red           = "\033[31m"
	yellow        = "\033[33m"
	blue          = "\033[34m"
	gray          = "\033[90m"
	grayAndItalic = "\033[90m\033[3m"
)

// Colors mapping.
var Colors = [...]string{
	log.DebugLevel: grayAndItalic,
	log.InfoLevel:  blue,
	log.WarnLevel:  yellow,
	log.ErrorLevel: red,
	log.FatalLevel: red,
}

var PrefixColor = grayAndItalic
var MessageColor = grayAndItalic

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  "INFO",
	log.WarnLevel:  "WARN",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

// Handler implementation.
type Handler struct {
	mu     sync.Mutex
	Writer io.Writer
	prefix string
}

// New handler.
func New(w io.Writer, prefix string) *Handler {
	return &Handler{
		Writer: w,
		prefix: prefix,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]
	names := e.Fields.Names()

	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Fprintf(h.Writer, "%s%s\033[0m%s%5s:\033[0m %s%-25s\033[0m",
		PrefixColor,
		h.prefix,
		color,
		level,
		MessageColor,
		e.Message)

	for _, name := range names {
		fmt.Fprintf(h.Writer, " %s%s\033[0m=%v\033[0m", color, name, e.Fields.Get(name))
	}

	fmt.Fprintln(h.Writer)

	return nil
}
