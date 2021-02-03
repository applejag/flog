package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/jilleJr/flog/pkg/loglevel"
	"github.com/jilleJr/flog/pkg/logparser"
)

var cli struct {
	MinLevel       string   `name:"min" short:"s" default:"info" help:"Filter out logs below specified severity (exclusive)" enum:"t,tra,trac,trce,trace,d,deb,debu,debg,dbug,debug,i,inf,info,information,w,wrn,warn,warning,fail,e,err,erro,errr,error,c,crt,crit,critical,f,fata,fatl,fatal,p,pan,pnc,pani,panic"`
	MaxLevel       string   `name:"max" short:"S" default:"none" help:"Filter out logs above specified severity (exclusive)" enum:"t,tra,trac,trce,trace,d,deb,debu,debg,dbug,debug,i,inf,info,information,w,wrn,warn,warning,fail,e,err,erro,errr,error,c,crt,crit,critical,f,fata,fatl,fatal,p,pan,pnc,pani,panic"`
	MinTime        string   `name:"since" short:"t" help:"Filter out logs timestamped before a specific time (or relative time period ago) [Not yet implemented]"`
	MaxTime        string   `name:"before" short:"t" help:"Filter out logs timestamped after a specific time (or relative time period ago) [Not yet implemented]"`
	ExcludedLevels []string `name:"exclude" short:"e" help:"Filter out logs of specified severity (can be specified multiple times)" enum:"u,?,ukwn,unknown,t,tra,trac,trce,trace,d,deb,debu,debg,dbug,debug,i,inf,info,information,w,wrn,warn,warning,fail,e,err,erro,errr,error,c,crt,crit,critical,f,fata,fatl,fatal,p,pan,pnc,pani,panic"`
	IncludedLevels []string `name:"include" short:"i" help:"Filter out logs of severity not specified with this flag (can be specified multiple times)" enum:"u,?,ukwn,unknown,t,tra,trac,trce,trace,d,deb,debu,debg,dbug,debug,i,inf,info,information,w,wrn,warn,warning,fail,e,err,erro,errr,error,c,crt,crit,critical,f,fata,fatl,fatal,p,pan,pnc,pani,panic"`
	Paths          []string `arg optional type:"existingfile" help:"File(s) to read logs from. Uses STDIN if unspecified"`
	Quiet          bool     `name:"quiet" short:"q" help:"Omit the 'omitted logs' messages."`
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
	kong.Parse(&cli)

	filter := LogFilter{
		MinLevel:      parseLevelArg(cli.MinLevel),
		MaxLevel:      parseLevelArg(cli.MinLevel),
		Quiet:         cli.Quiet,
		BlacklistMask: parseLevelArgsAsBitmask(cli.ExcludedLevels),
		WhitelistMask: parseLevelArgsAsBitmask(cli.IncludedLevels),
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

func parseLevelArgsAsBitmask(slice []string) loglevel.Level {
	m := loglevel.Undefined
	for _, lvlStr := range slice {
		m |= parseLevelArg(lvlStr)
	}
	return m
}

func parseLevelArg(s string) loglevel.Level {
	switch strings.ToLower(s) {
	case "u", "?", "ukwn", "unknown":
		return loglevel.Unknown

	case "t", "tra", "trac", "trce", "trace":
		return loglevel.Trace

	case "d", "deb", "dbg", "debu", "debg", "dbug", "debug":
		return loglevel.Debug

	case "i", "inf", "info", "information":
		return loglevel.Information

	case "w", "wrn", "warn", "warning":
		return loglevel.Warning

	case "fail", "e", "err", "erro", "errr", "error":
		return loglevel.Error

	case "c", "crt", "crit", "critical":
		return loglevel.Critical

	case "f", "fata", "fatl", "fatal":
		return loglevel.Fatal

	case "p", "pan", "pnc", "pani", "panic":
		return loglevel.Panic
	}

	return loglevel.Undefined
}
