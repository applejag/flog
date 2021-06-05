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
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/jilleJr/flog/internal/apex/handlers/console"
	"github.com/jilleJr/flog/pkg/loglevel"
	"github.com/jilleJr/flog/pkg/logparser"
)

func ExamplePrinter_logrus_text() {
	input := `
time="2021-01-31T19:04:01+01:00" level=trace msg="A walrus appears" animal=walrus
time="2021-01-31T19:04:01+01:00" level=debug msg="A walrus appears" animal=walrus
time="2021-01-31T19:04:01+01:00" level=info msg="A walrus appears" animal=walrus
time="2021-01-31T19:04:01+01:00" level=warning msg="A walrus appears" animal=walrus
time="2021-01-31T19:04:01+01:00" level=error msg="A walrus appears" animal=walrus
time="2021-01-31T19:04:01+01:00" level=fatal msg="A walrus appears" animal=walrus`

	r := strings.NewReader(input)
	p := logparser.NewIOReader(r)
	printer := NewConsolePrinter("test", &p, LogFilter{MinLevel: loglevel.Warning})
	log.SetHandler(console.New(os.Stdout, "flog: "))

	for printer.Next() {
	}

	// Output:
	// [90m[3mflog: [0m[34m INFO:[0m [90m[3mOmitted logs from: test  [0m [34mDebug[0m=1[0m [34mInformation[0m=1[0m [34mTrace[0m=1[0m
	// time="2021-01-31T19:04:01+01:00" level=warning msg="A walrus appears" animal=walrus
	// time="2021-01-31T19:04:01+01:00" level=error msg="A walrus appears" animal=walrus
	// time="2021-01-31T19:04:01+01:00" level=fatal msg="A walrus appears" animal=walrus
}

func ExamplePrinter_logrus_ansi_multiline() {
	input := `
[37mTRAC[0m[0000] A walrus appears                              [37manimal[0m=walrus
[37mDEBU[0m[0000] A walrus appears                              [37manimal[0m=walrus
[36mINFO[0m[0000] A walrus appears                              [36manimal[0m=walrus
[33mWARN[0m[0000] A walrus appears
	Some
	Multiline                              [33manimal[0m=walrus
[31mERRO[0m[0000] A walrus appears                              [31manimal[0m=walrus
[31mFATA[0m[0000] A walrus appears                              [31manimal[0m=walrus`

	r := strings.NewReader(input)
	p := logparser.NewIOReader(r)
	printer := NewConsolePrinter("test", &p, LogFilter{MinLevel: loglevel.Warning})
	log.SetHandler(console.New(os.Stdout, "flog: "))

	for printer.Next() {
	}

	// Output:
	// [90m[3mflog: [0m[34m INFO:[0m [90m[3mOmitted logs from: test  [0m [34mDebug[0m=1[0m [34mInformation[0m=1[0m [34mTrace[0m=1[0m
	// [33mWARN[0m[0000] A walrus appears
	//	Some
	//	Multiline                              [33manimal[0m=walrus
	// [31mERRO[0m[0000] A walrus appears                              [31manimal[0m=walrus
	// [31mFATA[0m[0000] A walrus appears                              [31manimal[0m=walrus
}
