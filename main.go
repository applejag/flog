package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/apex/log"
	"github.com/jilleJr/flog/internal/apex/handlers/console"
	"github.com/jilleJr/flog/pkg/loglevel"
	"github.com/jilleJr/flog/pkg/logparser"
)

type verbosityLevel int

var args struct {
	MinLevel       loglevel.Level   `name:"min" short:"s" default:"info" help:"Omit logs below specified severity (exclusive)"`
	MaxLevel       loglevel.Level   `name:"max" short:"S" default:"none" help:"Omit logs above specified severity (exclusive)"`
	MinTime        string           `name:"since" short:"t" help:"Omit logs timestamped before a specific time (or relative time period ago) [Not yet implemented]"`
	MaxTime        string           `name:"before" short:"t" help:"Omit logs timestamped after a specific time (or relative time period ago) [Not yet implemented]"`
	ExcludedLevels []loglevel.Level `name:"exclude" short:"e" help:"Omit logs of specified severity (can be specified multiple times)"`
	IncludedLevels []loglevel.Level `name:"include" short:"i" help:"Omit logs of severity not specified with this flag (can be specified multiple times)"`
	Paths          []string         `arg optional type:"existingfile" help:"File(s) to read logs from. Uses STDIN if unspecified"`
	Quiet          bool             `name:"quiet" short:"q" help:"Omit the 'omitted logs' messages. Shorthand for --verbose=0."`
	Verbose        verbosityLevel   `name:"verbose" short:"v" default:"1" type:"counter" help:"Enable verbose output (can be specified up to 2 times, ex: --verbose=2 or -vv)"`
	Version        kong.VersionFlag `help:"Show the text \"${version}\" and then exit."`
}

type LogFilter struct {
	MinLevel      loglevel.Level
	MaxLevel      loglevel.Level
	BlacklistMask loglevel.Level
	WhitelistMask loglevel.Level
}

var loggingLevel log.Level

func setLoggingLevel(quiet bool, v verbosityLevel) {
	if quiet || v <= 0 {
		loggingLevel = log.ErrorLevel
	} else if v == 1 {
		loggingLevel = log.InfoLevel
	} else {
		loggingLevel = log.DebugLevel
	}
	log.SetLevel(loggingLevel)
}

func main() {
	kong.Parse(&args,
		kong.Name("flog"),
		kong.Description("Use flog to filter logs on their serverity (even multiline logs), with automatic detection of log formats."),
		kong.Help(flogHelp),
		kong.Vars{
			"version": appVersion,
		},
		kong.TypeMapper(reflect.TypeOf(loglevel.Undefined), levelMapper{}))

	log.SetHandler(console.New(os.Stderr, "flog: "))
	setLoggingLevel(args.Quiet, args.Verbose)

	filter := LogFilter{
		MinLevel:      args.MinLevel,
		MaxLevel:      args.MaxLevel,
		BlacklistMask: sliceOfArgsAsBitmask(args.ExcludedLevels),
		WhitelistMask: sliceOfArgsAsBitmask(args.IncludedLevels),
	}

	log.WithFields(log.Fields{
		"MinLevel":      filter.MinLevel,
		"MaxLevel":      filter.MaxLevel,
		"WhitelistMask": filter.WhitelistMask,
		"BlacklistMask": filter.BlacklistMask,
	}).Debugf("Parsed filter")

	if len(args.Paths) > 0 {
		for _, path := range args.Paths {
			printLogsFromFile(path, filter)
		}
	} else {
		printLogsFromIO("STDIN", os.Stdin, filter)
	}
}

func printLogsFromFile(path string, filter LogFilter) {
	if file, err := os.Open(path); err != nil {
		fmt.Printf("ERR: Failed to open file: %s: %v\n", path, err)
		os.Exit(1)
	} else {
		defer file.Close()
		printLogsFromIO(file.Name(), file, filter)
	}
}

func printLogsFromIO(name string, r io.Reader, filter LogFilter) {
	p := logparser.NewIOParser(r)

	printer := NewConsolePrinter(name, &p, filter)
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
