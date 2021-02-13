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

package main

import (
	"fmt"
	"testing"

	"github.com/jilleJr/flog/pkg/loglevel"
)

func TestGetSkippedLevelsFields_Empty(t *testing.T) {
	input := map[loglevel.Level]int{}
	want := 0

	got := getSkippedLevelsFields(input)

	if len(got) != want {
		t.Errorf("slice was not empty: %#v", got)
	}
}

func TestShouldIncludeLogInOutput_MinLevel(t *testing.T) {
	var testCases = []struct {
		input    loglevel.Level
		minLevel loglevel.Level
		want     bool
	}{
		{
			input:    loglevel.Fatal,
			minLevel: loglevel.Information,
			want:     true,
		},
		{
			input:    loglevel.Information,
			minLevel: loglevel.Information,
			want:     true,
		},
		{
			input:    loglevel.Debug,
			minLevel: loglevel.Information,
			want:     false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d/input/%v/minLevel/%v", i, tc.input, tc.minLevel), func(t *testing.T) {
			got := shouldIncludeLogInOutput(tc.input, LogFilter{MinLevel: tc.minLevel})
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestShouldIncludeLogInOutput_MaxLevel(t *testing.T) {
	var testCases = []struct {
		input    loglevel.Level
		maxLevel loglevel.Level
		want     bool
	}{
		{
			input:    loglevel.Fatal,
			maxLevel: loglevel.Information,
			want:     false,
		},
		{
			input:    loglevel.Information,
			maxLevel: loglevel.Information,
			want:     true,
		},
		{
			input:    loglevel.Debug,
			maxLevel: loglevel.Information,
			want:     true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d/input/%v/maxLevel/%v", i, tc.input, tc.maxLevel), func(t *testing.T) {
			got := shouldIncludeLogInOutput(tc.input, LogFilter{MaxLevel: tc.maxLevel})
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestShouldIncludeLogInOutput_BlacklistMask(t *testing.T) {
	var testCases = []struct {
		input         loglevel.Level
		blacklistMask loglevel.Level
		want          bool
	}{
		{
			input:         loglevel.Warning,
			blacklistMask: loglevel.Information | loglevel.Panic,
			want:          true,
		},
		{
			input:         loglevel.Information,
			blacklistMask: loglevel.Information | loglevel.Panic,
			want:          false,
		},
		{
			input:         loglevel.Panic,
			blacklistMask: loglevel.Information | loglevel.Panic,
			want:          false,
		},
		{
			input:         loglevel.Debug,
			blacklistMask: loglevel.Information | loglevel.Panic,
			want:          true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d/input/%v/blacklistMask/%v", i, tc.input, tc.blacklistMask), func(t *testing.T) {
			got := shouldIncludeLogInOutput(tc.input, LogFilter{BlacklistMask: tc.blacklistMask})
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestShouldIncludeLogInOutput_WhitelistMask(t *testing.T) {
	var testCases = []struct {
		input         loglevel.Level
		whitelistMask loglevel.Level
		want          bool
	}{
		{
			input:         loglevel.Warning,
			whitelistMask: loglevel.Information | loglevel.Panic,
			want:          false,
		},
		{
			input:         loglevel.Information,
			whitelistMask: loglevel.Information | loglevel.Panic,
			want:          true,
		},
		{
			input:         loglevel.Panic,
			whitelistMask: loglevel.Information | loglevel.Panic,
			want:          true,
		},
		{
			input:         loglevel.Debug,
			whitelistMask: loglevel.Information | loglevel.Panic,
			want:          false,
		},
		{
			input:         loglevel.Debug,
			whitelistMask: loglevel.Undefined,
			want:          true,
		},
		{
			input:         loglevel.Debug,
			whitelistMask: loglevel.Unknown,
			want:          false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d/input/%v/whitelistMask/%v", i, tc.input, tc.whitelistMask), func(t *testing.T) {
			got := shouldIncludeLogInOutput(tc.input, LogFilter{WhitelistMask: tc.whitelistMask})
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
