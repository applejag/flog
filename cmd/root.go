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

package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/jilleJr/flog/internal/apex/handlers/console"
	"github.com/jilleJr/flog/pkg/flagtype"
	"github.com/jilleJr/flog/pkg/license"
	"github.com/jilleJr/flog/pkg/loglevel"
	"github.com/jilleJr/flog/pkg/logparser"
	"github.com/jilleJr/flog/pkg/printer"
	"github.com/spf13/cobra"
)

var flags struct {
	minLevel       flagtype.LogLevel
	maxLevel       flagtype.LogLevel
	minTime        string
	maxTime        string
	excludedLevels flagtype.LogLevelMask
	includedLevels flagtype.LogLevelMask
	quiet          bool
	verbose        int

	showLicenseConditions bool
	showLicenseWarranty   bool
}

var loggingLevel log.Level

func setLoggingLevel(quiet bool, v int) {
	if quiet || v <= 0 {
		loggingLevel = log.ErrorLevel
	} else if v == 1 {
		loggingLevel = log.InfoLevel
	} else {
		loggingLevel = log.DebugLevel
	}
	log.SetLevel(loggingLevel)
}

var rootCmd = &cobra.Command{
	Use:   "flog [flags] [file1.log [file2.log [file3.log]]]",
	Short: "Filter logs on their serverity (even multiline logs), with automatic detection of log formats",

	Run: func(cmd *cobra.Command, args []string) {

		if flags.showLicenseConditions {
			fmt.Println(license.LicenseConditions)
			return
		} else if flags.showLicenseWarranty {
			fmt.Println(license.LicenseWarranty)
			return
		}

		log.SetHandler(console.New(os.Stderr, "flog: "))
		setLoggingLevel(flags.quiet, flags.verbose)

		filter := loglevel.Filter{
			MinLevel:      flags.minLevel.Level(),
			MaxLevel:      flags.maxLevel.Level(),
			BlacklistMask: flags.excludedLevels.Level(),
			WhitelistMask: flags.includedLevels.Level(),
		}

		log.WithFields(log.Fields{
			"MinLevel":      filter.MinLevel,
			"MaxLevel":      filter.MaxLevel,
			"WhitelistMask": filter.WhitelistMask,
			"BlacklistMask": filter.BlacklistMask,
		}).Debugf("Parsed filter")

		if len(args) > 0 {
			for _, path := range args {
				printLogsFromFile(path, filter)
			}
		} else {
			printLogsFromIO("STDIN", os.Stdin, filter)
		}
	},
}

func Execute(appVersion string) {
	rootCmd.Version = appVersion
	rootCmd.SetVersionTemplate(license.VersionNotice(appVersion))
	rootCmd.Long = fmt.Sprintf(`Use flog to filter logs on their serverity (even multiline logs),
with automatic detection of log formats.

%s
`, license.LicenceNotice(appVersion))
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().VarP(&flags.minLevel, "min", "s", "Omit logs below specified severity (exclusive)")
	rootCmd.Flags().VarP(&flags.maxLevel, "max", "S", "Omit logs above specified severity (exclusive)")
	rootCmd.Flags().StringVarP(&flags.minTime, "since", "t", "", "Omit logs timestamped before a specific time (or relative time period ago) [Not yet implemented]")
	rootCmd.Flags().StringVarP(&flags.minTime, "before", "T", "", "Omit logs timestamped after a specific time (or relative time period ago) [Not yet implemented]")
	rootCmd.Flags().VarP(&flags.excludedLevels, "exclude", "e", "Omit logs of specified severity (can be specified multiple times)")
	rootCmd.Flags().VarP(&flags.includedLevels, "include", "i", "Omit logs of severity not specified with this flag (can be specified multiple times)")
	rootCmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", flags.quiet, "Omit the 'omitted logs' messages. Shorthand for --verbose=0.")
	rootCmd.Flags().CountVarP(&flags.verbose, "verbose", "v", "Enable verbose output (can be specified up to 2 times, ex: --verbose=2 or -vv)")
	rootCmd.Flags().BoolVar(&flags.showLicenseConditions, "license-c", false, "Show the program's license conditions and then exit. (Warn: a lot of text)")
	rootCmd.Flags().BoolVar(&flags.showLicenseWarranty, "license-w", false, "Show the program's warranty and then exit.")
	rootCmd.Flags().Bool("version", false, "Show the program's version and then exit.")
	rootCmd.Flags().Bool("help", false, "Show this help text and then exit.")
}

func printLogsFromFile(path string, filter loglevel.Filter) {
	if file, err := os.Open(path); err != nil {
		fmt.Printf("ERR: Failed to open file: %s: %v\n", path, err)
		os.Exit(1)
	} else {
		defer file.Close()
		printLogsFromIO(file.Name(), file, filter)
	}
}

func printLogsFromIO(name string, r io.Reader, filter loglevel.Filter) {
	logread := logparser.NewIOReader(r)

	p := printer.NewConsolePrinter(name, &logread, filter, loggingLevel)
	ch := setupCloseHandler(p)
	defer close(ch)

	for p.Next() {
	}
}

// Thanks https://golangcode.com/handle-ctrl-c-exit-in-terminal/
// His site shows 404, but the source code is supposed to be found here:
// https://github.com/eddturtle/golangcode-site
func setupCloseHandler(p printer.Printer) chan<- os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func(p printer.Printer) {
		if _, ok := <-ch; ok {
			p.PrintOmittedLogs()
			os.Exit(0)
		}
	}(p)
	return ch
}
