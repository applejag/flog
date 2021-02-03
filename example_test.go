package main

import (
	"strings"

	"github.com/jilleJr/flog/pkg/loglevel"
	"github.com/jilleJr/flog/pkg/logparser"
)

func ExamplePrinter_nlog_text() {
	input := `
2021-01-31 17:33:54.3326|TRACE|Program|Sample
2021-01-31 17:33:54.3443|DEBUG|Program|Sample
2021-01-31 17:33:54.3443|INFO|Program|Sample
2021-01-31 17:33:54.3443|WARN|Program|Sample
2021-01-31 17:33:54.3443|ERROR|Program|Sample
2021-01-31 17:33:54.3443|FATAL|Program|Sample`

	r := strings.NewReader(input)
	p := logparser.NewIOParser(r)
	printer := NewConsolePrinter(&p, LogFilter{MinLevel: loglevel.Warning})

	for printer.Next() {
	}

	// Output:
	// [90m[3mflog: Omitted 1 Trace, 1 Debug, 1 Information.[0m
	// 2021-01-31 17:33:54.3443|WARN|Program|Sample
	// 2021-01-31 17:33:54.3443|ERROR|Program|Sample
	// 2021-01-31 17:33:54.3443|FATAL|Program|Sample
}

func ExamplePrinter_nlog_text_multiline() {
	input := `
2021-01-31 17:33:54.3326|TRACE|Program|Sample
2021-01-31 17:33:54.3443|DEBUG|Program|Sample
2021-01-31 17:33:54.3443|INFO|Program|Sample
2021-01-31 17:33:54.3443|WARN|Program|Sample
	some other text
	this still counts as the WARN message
2021-01-31 17:33:54.3443|ERROR|Program|Sample
2021-01-31 17:33:54.3443|FATAL|Program|Sample`

	r := strings.NewReader(input)
	p := logparser.NewIOParser(r)
	printer := NewConsolePrinter(&p, LogFilter{MinLevel: loglevel.Warning})

	for printer.Next() {
	}

	// Output:
	// [90m[3mflog: Omitted 1 Trace, 1 Debug, 1 Information.[0m
	// 2021-01-31 17:33:54.3443|WARN|Program|Sample
	//	some other text
	//	this still counts as the WARN message
	// 2021-01-31 17:33:54.3443|ERROR|Program|Sample
	// 2021-01-31 17:33:54.3443|FATAL|Program|Sample
}

func ExamplePrinter_nlog_ansi() {
	input := `
2021-01-31 18:37:05.1550|TRACE|Program|Sample
2021-01-31 18:37:05.1714|DEBUG|Program|Sample
2021-01-31 18:37:05.1714|INFO|Program|Sample
[35m2021-01-31 18:37:05.1714|WARN|Program|Sample[0m
[33m2021-01-31 18:37:05.1714|ERROR|Program|Sample[0m
[31m2021-01-31 18:37:05.1714|FATAL|Program|Sample[0m`

	r := strings.NewReader(input)
	p := logparser.NewIOParser(r)
	printer := NewConsolePrinter(&p, LogFilter{MinLevel: loglevel.Warning})

	for printer.Next() {
	}

	// Output:
	// [90m[3mflog: Omitted 1 Trace, 1 Debug, 1 Information.[0m
	// [35m2021-01-31 18:37:05.1714|WARN|Program|Sample[0m
	// [33m2021-01-31 18:37:05.1714|ERROR|Program|Sample[0m
	// [31m2021-01-31 18:37:05.1714|FATAL|Program|Sample[0m
}

func ExamplePrinter_logrus_text() {
	input := `
time="2021-01-31T19:04:01+01:00" level=trace msg="A walrus appears" animal=walrus
time="2021-01-31T19:04:01+01:00" level=debug msg="A walrus appears" animal=walrus
time="2021-01-31T19:04:01+01:00" level=info msg="A walrus appears" animal=walrus
time="2021-01-31T19:04:01+01:00" level=warning msg="A walrus appears" animal=walrus
time="2021-01-31T19:04:01+01:00" level=error msg="A walrus appears" animal=walrus
time="2021-01-31T19:04:01+01:00" level=fatal msg="A walrus appears" animal=walrus`

	r := strings.NewReader(input)
	p := logparser.NewIOParser(r)
	printer := NewConsolePrinter(&p, LogFilter{MinLevel: loglevel.Warning})

	for printer.Next() {
	}

	// Output:
	// [90m[3mflog: Omitted 1 Trace, 1 Debug, 1 Information.[0m
	// time="2021-01-31T19:04:01+01:00" level=warning msg="A walrus appears" animal=walrus
	// time="2021-01-31T19:04:01+01:00" level=error msg="A walrus appears" animal=walrus
	// time="2021-01-31T19:04:01+01:00" level=fatal msg="A walrus appears" animal=walrus
}

func ExamplePrinter_logrus_ansi() {
	input := `
[37mTRAC[0m[0000] A walrus appears                              [37manimal[0m=walrus
[37mDEBU[0m[0000] A walrus appears                              [37manimal[0m=walrus
[36mINFO[0m[0000] A walrus appears                              [36manimal[0m=walrus
[33mWARN[0m[0000] A walrus appears                              [33manimal[0m=walrus
[31mERRO[0m[0000] A walrus appears                              [31manimal[0m=walrus
[31mFATA[0m[0000] A walrus appears                              [31manimal[0m=walrus`

	r := strings.NewReader(input)
	p := logparser.NewIOParser(r)
	printer := NewConsolePrinter(&p, LogFilter{MinLevel: loglevel.Warning})

	for printer.Next() {
	}

	// Output:
	// [90m[3mflog: Omitted 1 Trace, 1 Debug, 1 Information.[0m
	// [33mWARN[0m[0000] A walrus appears                              [33manimal[0m=walrus
	// [31mERRO[0m[0000] A walrus appears                              [31manimal[0m=walrus
	// [31mFATA[0m[0000] A walrus appears                              [31manimal[0m=walrus
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
	p := logparser.NewIOParser(r)
	printer := NewConsolePrinter(&p, LogFilter{MinLevel: loglevel.Warning})

	for printer.Next() {
	}

	// Output:
	// [90m[3mflog: Omitted 1 Trace, 1 Debug, 1 Information.[0m
	// [33mWARN[0m[0000] A walrus appears
	//	Some
	//	Multiline                              [33manimal[0m=walrus
	// [31mERRO[0m[0000] A walrus appears                              [31manimal[0m=walrus
	// [31mFATA[0m[0000] A walrus appears                              [31manimal[0m=walrus
}
