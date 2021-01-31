package main

import (
	"strings"

	"github.com/jilleJr/flog/pkg/logparser"
)

func ExamplePrinter_nlog_pure() {
	input := `2021-01-31 17:33:54.3326|TRACE|Program|Sample
2021-01-31 17:33:54.3443|DEBUG|Program|Sample
2021-01-31 17:33:54.3443|INFO|Program|Sample
2021-01-31 17:33:54.3443|WARN|Program|Sample
2021-01-31 17:33:54.3443|ERROR|Program|Sample
2021-01-31 17:33:54.3443|FATAL|Program|Sample`

	r := strings.NewReader(input)
	p := logparser.NewIOParser(r)
	printer := NewConsolePrinter(&p)

	for printer.Next() {
	}

	// Output:
	// 2021-01-31 17:33:54.3443|WARN|Program|Sample
	// 2021-01-31 17:33:54.3443|ERROR|Program|Sample
	// 2021-01-31 17:33:54.3443|FATAL|Program|Sample
}

func ExamplePrinter_nlog_multiline() {
	input := `2021-01-31 17:33:54.3326|TRACE|Program|Sample
2021-01-31 17:33:54.3443|DEBUG|Program|Sample
2021-01-31 17:33:54.3443|INFO|Program|Sample
2021-01-31 17:33:54.3443|WARN|Program|Sample
	some other text
	this still counts as the WARN message
2021-01-31 17:33:54.3443|ERROR|Program|Sample
2021-01-31 17:33:54.3443|FATAL|Program|Sample`

	r := strings.NewReader(input)
	p := logparser.NewIOParser(r)
	printer := NewConsolePrinter(&p)

	for printer.Next() {
	}

	// Output:
	// 2021-01-31 17:33:54.3443|WARN|Program|Sample
	//	some other text
	//	this still counts as the WARN message
	// 2021-01-31 17:33:54.3443|ERROR|Program|Sample
	// 2021-01-31 17:33:54.3443|FATAL|Program|Sample
}

func ExamplePrinter_nlog_ansi() {
	input := `2021-01-31 18:37:05.1550|TRACE|Program|Sample
2021-01-31 18:37:05.1714|DEBUG|Program|Sample
2021-01-31 18:37:05.1714|INFO|Program|Sample
[35m2021-01-31 18:37:05.1714|WARN|Program|Sample[0m
[33m2021-01-31 18:37:05.1714|ERROR|Program|Sample[0m
[31m2021-01-31 18:37:05.1714|FATAL|Program|Sample[0m`

	r := strings.NewReader(input)
	p := logparser.NewIOParser(r)
	printer := NewConsolePrinter(&p)

	for printer.Next() {
	}

	// Output:
	// [35m2021-01-31 18:37:05.1714|WARN|Program|Sample[0m
	// [33m2021-01-31 18:37:05.1714|ERROR|Program|Sample[0m
	// [31m2021-01-31 18:37:05.1714|FATAL|Program|Sample[0m
}
