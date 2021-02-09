package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/jilleJr/flog/pkg/loglevel"
	"github.com/jilleJr/flog/pkg/logparser"
)

var cli struct {
	MinLevel       loglevel.Level   `name:"min" short:"s" default:"info" help:"Omit logs below specified severity (exclusive)"`
	MaxLevel       loglevel.Level   `name:"max" short:"S" default:"none" help:"Omit logs above specified severity (exclusive)"`
	MinTime        string           `name:"since" short:"t" help:"Omit logs timestamped before a specific time (or relative time period ago) [Not yet implemented]"`
	MaxTime        string           `name:"before" short:"t" help:"Omit logs timestamped after a specific time (or relative time period ago) [Not yet implemented]"`
	ExcludedLevels []loglevel.Level `name:"exclude" short:"e" help:"Omit logs of specified severity (can be specified multiple times)"`
	IncludedLevels []loglevel.Level `name:"include" short:"i" help:"Omit logs of severity not specified with this flag (can be specified multiple times)"`
	Paths          []string         `arg optional type:"existingfile" help:"File(s) to read logs from. Uses STDIN if unspecified"`
	Quiet          bool             `name:"quiet" short:"q" help:"Omit the 'omitted logs' messages."`
}

type LogFilter struct {
	MinLevel      loglevel.Level
	MaxLevel      loglevel.Level
	Quiet         bool
	BlacklistMask loglevel.Level
	WhitelistMask loglevel.Level
}

var logFilter LogFilter

func main() {
	kong.Parse(&cli,
		kong.Help(flogHelp),
		kong.TypeMapper(reflect.TypeOf(loglevel.Undefined), levelMapper{}))

	filter := LogFilter{
		MinLevel:      cli.MinLevel,
		MaxLevel:      cli.MaxLevel,
		Quiet:         cli.Quiet,
		BlacklistMask: sliceOfArgsAsBitmask(cli.ExcludedLevels),
		WhitelistMask: sliceOfArgsAsBitmask(cli.IncludedLevels),
	}

	if len(cli.Paths) > 0 {
		for _, path := range cli.Paths {
			printLogsFromFile(path, filter)
		}
	} else {
		printLogsFromIO(os.Stdin, filter)
	}
}

func printLogsFromFile(path string, filter LogFilter) {
	if file, err := os.Open(path); err != nil {
		fmt.Printf("ERR: Failed to open file: %s: %v\n", path, err)
		os.Exit(1)
	} else {
		defer file.Close()
		printLogsFromIO(file, filter)
	}

}

func printLogsFromIO(r io.Reader, filter LogFilter) {
	p := logparser.NewIOParser(r)

	printer := NewConsolePrinter(&p, filter)
	ch := setupCloseHandler(printer)
	defer close(ch)

	for printer.Next() {
	}
}

// Thanks https://golangcode.com/handle-ctrl-c-exit-in-terminal/
// His site shows 404, but the source code is supposed to be found here:
// https://github.com/eddturtle/golangcode-site
func setupCloseHandler(p Printer) chan<- os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func(p Printer) {
		if _, ok := <-ch; ok {
			p.PrintOmittedLogs()
			os.Exit(0)
		}
	}(p)
	return ch
}

func sliceOfArgsAsBitmask(slice []loglevel.Level) loglevel.Level {
	m := loglevel.Undefined
	for _, lvl := range slice {
		m |= lvl
	}
	return m
}
