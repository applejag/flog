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
	"github.com/posener/complete"
	"github.com/willabides/kongplete"
)

type verbosityLevel int

var args struct {
	MinLevel       loglevel.Level   `name:"min" short:"s" default:"info" help:"Omit logs below specified severity (exclusive)" predictor:"severity"`
	MaxLevel       loglevel.Level   `name:"max" short:"S" default:"none" help:"Omit logs above specified severity (exclusive)" predictor:"severity"`
	MinTime        string           `name:"since" short:"t" help:"Omit logs timestamped before a specific time (or relative time period ago) [Not yet implemented]"`
	MaxTime        string           `name:"before" short:"t" help:"Omit logs timestamped after a specific time (or relative time period ago) [Not yet implemented]"`
	ExcludedLevels []loglevel.Level `name:"exclude" short:"e" help:"Omit logs of specified severity (can be specified multiple times)" predictor:"severity"`
	IncludedLevels []loglevel.Level `name:"include" short:"i" help:"Omit logs of severity not specified with this flag (can be specified multiple times)" predictor:"severity"`
	Paths          []string         `arg optional type:"existingfile" help:"File(s) to read logs from. Uses STDIN if unspecified" predictor:"file"`
	Quiet          bool             `name:"quiet" short:"q" help:"Omit the 'omitted logs' messages. Shorthand for --verbose=0."`
	Verbose        verbosityLevel   `name:"verbose" short:"v" default:"1" type:"counter" help:"Enable verbose output (can be specified up to 2 times, ex: --verbose=2 or -vv)"`
	Version        kong.VersionFlag `help:"Show the version of the program and then exit."`

	CompletionsInstall   bool `name:"completions-install" help:"Install shell completions" xor:"completions"`
	CompletionsUninstall bool `name:"completions-uninstall" help:"Uninstall shell completions" xor:"completions"`

	LicenseConditions bool `name:"license-c" help:"Show the programs license conditions and then exit. (Warn: a lot of text)"`
	LicenseWarranty   bool `name:"license-w" help:"Show the programs warranty and then exit."`
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
	parser := kong.Must(&args,
		kong.Name("flog"),
		kong.Description("Use flog to filter logs on their serverity (even multiline logs), with automatic detection of log formats.\n\n${licenseNotice}"),
		kong.Help(flogHelp),
		kong.Vars{
			"version":       versionNotice,
			"licenseNotice": LicenceNotice,
		},
		kong.TypeMapper(reflect.TypeOf(loglevel.Undefined), levelMapper{}))

	kongplete.Complete(parser,
		kongplete.WithPredictor("file", complete.PredictFiles("*")),
		kongplete.WithPredictor("severity", complete.PredictSet(predictLevelSuggestions...)),
	)

	kctx, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	if args.LicenseConditions {
		showLicenseConditionsThenExit()
	} else if args.LicenseWarranty {
		showLicenseWarrantyThenExit()
	}

	log.SetHandler(console.New(os.Stderr, "flog: "))
	setLoggingLevel(args.Quiet, args.Verbose)

	if args.CompletionsUninstall {
		if err := kongplete.UninstallShellCompletions(kctx); err != nil {
			log.WithError(err).Error("Failed to uninstall shell completions.")
		} else {
			log.Info("Uninstalled. You might have to restart your shell for it to take affect.")
		}
		os.Exit(0)
	} else if args.CompletionsInstall {
		if err := kongplete.InstallShellCompletions(kctx); err != nil {
			log.WithError(err).Error("Failed to uninstall shell completions.")
		} else {
			log.Info("Installed. You might have to restart your shell for it to take affect.")
		}
		os.Exit(0)
	}

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
