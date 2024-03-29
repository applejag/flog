// Copyright (C) 2021  Kalle Jillheden
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

package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/jilleJr/flog/pkg/loglevel"
)

func flogHelp(options kong.HelpOptions, ctx *kong.Context) error {
	if err := kong.DefaultHelpPrinter(options, ctx); err != nil {
		return err
	}

	fmt.Println(`
Severities:
  Undefined      0, n, nil, null, none, unde, undefined
  Unknown        1, ?, u, ukwn, unkn, unknown
  Trace          2, t, tra, trac, trce, trace
  Debug          3, d, deb, debu, debg, dbug, debug
  Information    4, i, inf, info, information
  Warning        5, w, wrn, warn, warning
  Error          6, e, err, erro, errr, error
  Critical       7, c, crt, crit, critical
  Fatal          8, f, fata, fatl, fatal
  Panic          9, p, pan, pnc, pani, panic`)

	return nil
}

type levelMapper struct {
}

func (m levelMapper) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	ttype := ctx.Scan.Peek().Type
	switch ttype {
	case kong.FlagValueToken, kong.UntypedToken:
		token := ctx.Scan.Pop()
		switch v := token.Value.(type) {
		case string:
			if lvl, err := parseLevelString(v); err != nil {
				return fmt.Errorf("failed to parse severity: %w", err)
			} else {
				target.Set(reflect.ValueOf(lvl))
				return nil
			}
		case int:
			if v >= int(loglevel.Unknown) && v <= int(loglevel.Panic) {
				target.Set(reflect.ValueOf(v))
				return nil
			} else {
				return fmt.Errorf("severity is out of range %d...%d: %d", int(loglevel.Unknown), int(loglevel.Panic), v)
			}
		default:
			return fmt.Errorf("expected string or int, got: %T", reflect.TypeOf(v))
		}

	case kong.ShortFlagTailToken:
		s := ctx.Scan.Pop().String()
		return fmt.Errorf("severity must be specified, but got the following tail flags: %s", strings.Join(strings.Split(s, ""), ", "))

	default:
		return fmt.Errorf("expected severity, got: %v", ttype)
	}
}

var predictLevelSuggestions = []string{
	"0", "nil", "null", "none", "unde", "undefined",
	"1", "u", "?", "ukwn", "unkn", "unknown",
	"2", "t", "tra", "trac", "trce", "trace",
	"3", "d", "deb", "dbg", "debu", "debg", "dbug", "debug",
	"4", "i", "inf", "info", "information",
	"5", "w", "wrn", "warn", "warning",
	"6", "fail", "e", "err", "erro", "errr", "error",
	"7", "c", "crt", "crit", "critical",
	"8", "f", "fata", "fatl", "fatal",
	"9", "p", "pan", "pnc", "pani", "panic",
}

func parseLevelString(s string) (loglevel.Level, error) {
	switch strings.ToLower(s) {
	case "0", "nil", "null", "none", "unde", "undefined":
		return loglevel.Undefined, nil

	case "1", "u", "?", "ukwn", "unkn", "unknown":
		return loglevel.Unknown, nil

	case "2", "t", "tra", "trac", "trce", "trace":
		return loglevel.Trace, nil

	case "3", "d", "deb", "dbg", "debu", "debg", "dbug", "debug":
		return loglevel.Debug, nil

	case "4", "i", "inf", "info", "information":
		return loglevel.Information, nil

	case "5", "w", "wrn", "warn", "warning":
		return loglevel.Warning, nil

	case "6", "fail", "e", "err", "erro", "errr", "error":
		return loglevel.Error, nil

	case "7", "c", "crt", "crit", "critical":
		return loglevel.Critical, nil

	case "8", "f", "fata", "fatl", "fatal":
		return loglevel.Fatal, nil

	case "9", "p", "pan", "pnc", "pani", "panic":
		return loglevel.Panic, nil
	}

	return loglevel.Undefined, fmt.Errorf("unknown: %s", s)
}
