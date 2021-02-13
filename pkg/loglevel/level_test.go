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

package loglevel

import (
	"fmt"
	"testing"
)

func TestShouldIncludeLogInOutput_WhitelistMask(t *testing.T) {
	var testCases = []struct {
		input Level
		want  string
	}{
		{
			input: Warning,
			want:  "Warning",
		},
		{
			input: Information | Panic,
			want:  "Information|Panic",
		},
		{
			input: Undefined,
			want:  "Undefined",
		},
		{
			input: Undefined | Debug,
			want:  "Debug",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d/input/%v", i, tc.input), func(t *testing.T) {
			got := tc.input.String()
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
