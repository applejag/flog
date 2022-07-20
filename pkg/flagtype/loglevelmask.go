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

type LogLevelMask loglevel.Level

// Ensure it conforms to the interface
var lvlSlice LogLevelMask = 0
var _ pflag.SliceValue = &lvlSlice
var _ pflag.Value = &lvlSlice

func (s *LogLevelMask) Level() loglevel.Level {
	return loglevel.Level(*s)
}

func (s *LogLevelMask) String() string {
	return s.Level().String()
}

func (s *LogLevelMask) Set(str string) error {
	newLvl := loglevel.ParseLevel(str)
	if newLvl == loglevel.Unknown {
		return fmt.Errorf("unknown log level: %q", str)
	}
	*s |= LogLevelMask(newLvl)
	return nil
}

func (s *LogLevelMask) Type() string {
	return "loglevels"
}

// Append adds the specified value to the end of the flag value list.
func (s *LogLevelMask) Append(str string) error {
	newLvl := loglevel.ParseLevel(str)
	if newLvl == loglevel.Unknown {
		return fmt.Errorf("unknown log level: %q", str)
	}
	*s = LogLevelMask(s.Level() | newLvl)
	return nil
}

// Replace will fully overwrite any data currently in the flag value list.
func (s *LogLevelMask) Replace(slice []string) error {
	var newBitmap loglevel.Level
	for _, str := range slice {
		newLvl := loglevel.ParseLevel(str)
		if newLvl == loglevel.Unknown {
			return fmt.Errorf("unknown log level: %q", str)
		}
		newBitmap |= newLvl
	}
	*s = LogLevelMask(newBitmap)
	return nil
}

// GetSlice returns the flag value list as an array of strings.
func (s *LogLevelMask) GetSlice() []string {
	if s == nil {
		return nil
	}
	lvls := s.Level().Levels()
	slice := make([]string, len(lvls))
	for i, lvl := range lvls {
		slice[i] = lvl.String()
	}
	return slice
}
