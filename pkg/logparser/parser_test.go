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

package logparser

import (
	"testing"
	"time"

	"github.com/jilleJr/flog/pkg/loglevel"
	"gopkg.in/guregu/null.v3"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name  string
		line  string
		level loglevel.Level
		time  null.Time
	}{
		{
			name:  "nlog",
			line:  `2021-01-31 17:33:54.3326|TRACE|Program|Sample`,
			level: loglevel.Trace,
			time:  null.TimeFrom(time.Date(2021, 1, 31, 17, 33, 54, 332600, time.UTC)),
		},
		{
			name:  "logrus",
			line:  `time="2021-01-31T19:04:01+01:00" level=info msg="A walrus appears" animal=walrus`,
			level: loglevel.Information,
			time:  null.TimeFrom(time.Date(2021, 1, 31, 19, 04, 01, 0, time.FixedZone("", 60*60))),
		},
		{
			name: "logrus_ansi",
			line: `[33mWARN[0m[0000] A walrus appears            [33manimal[0m=walrus`,
			level: loglevel.Warning,
			time:  null.Time{},
		},
		{
			name:  "dotnet",
			line:  `fail: Program[0]`,
			level: loglevel.Error,
			time:  null.Time{},
		},
		{
			// https://github.com/jilleJr/flog/issues/8
			name:  "klog",
			line:  `I0204 09:00:44.662471       i health.go:55] Starting MySQL health checker...`,
			level: loglevel.Information,
			time:  null.TimeFrom(time.Date(time.Now().Year(), 2, 4, 9, 0, 44, 662471, time.UTC)),
		},
		{
			name:  "iver-wharf/wharf-core",
			line:  `Jun-18 14:50+0200 [DEBUG | TEST | wharf-core/main.go:23] Sample  hello=world`,
			level: loglevel.Debug,
			time:  null.TimeFrom(time.Date(time.Now().Year(), 6, 18, 14, 50, 0, 0, time.FixedZone("", 60*60*2))),
		},
		{
			name:  "json",
			line:  `{"level":"debug","timestamp":"2021-06-05T23:50:00Z","message":"foo bar"}`,
			level: loglevel.Debug,
			time:  null.TimeFrom(time.Date(2021, 6, 5, 23, 50, 0, 0, time.UTC)),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			log := ParseUsingAnyParser(tc.line)
			if tc.level != log.Level {
				t.Errorf("wrong log level\nwanted: %s\ngot:    %s", tc.level, log.Level)
			}
			if !timeEquals(tc.time, log.Timestamp) {
				t.Errorf("wrong time\nwanted: %s\ngot:    %s", nullTimeString(tc.time), nullTimeString(log.Timestamp))
			}
		})
	}
}

func nullTimeString(t null.Time) string {
	if t.Valid {
		return t.Time.Format(time.RFC3339)
	} else {
		return "<null>"
	}
}

func timeEquals(t1, t2 null.Time) bool {
	return nullTimeString(t1) == nullTimeString(t2)
}
