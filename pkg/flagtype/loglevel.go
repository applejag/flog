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

package flagtype

import (
	"fmt"

	"github.com/jilleJr/flog/pkg/loglevel"
	"github.com/spf13/pflag"
)

type LogLevel loglevel.Level

// Ensure it conforms to the interface
var lvl LogLevel = 0
var _ pflag.Value = &lvl

func (lvl *LogLevel) Level() loglevel.Level {
	return loglevel.Level(*lvl)
}

func (lvl *LogLevel) String() string {
	return lvl.Level().String()
}

func (lvl *LogLevel) Set(str string) error {
	newLvl := loglevel.ParseLevel(str)
	if newLvl == loglevel.Unknown {
		return fmt.Errorf("unknown log level: %q", str)
	}
	*lvl = LogLevel(newLvl)
	return nil
}

func (lvl *LogLevel) Type() string {
	return "loglevel"
}
