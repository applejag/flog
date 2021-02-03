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
	MinLevel       string   `name:"min" short:"s" default:"info" help:"Filter out logs below specified severity (exclusive)" enum:"n,non,none,t,tra,trac,trce,trace,d,deb,debu,debg,dbug,debug,i,inf,info,information,w,wrn,warn,warning,fail,e,err,erro,errr,error,c,crt,crit,critical,f,fata,fatl,fatal,p,pan,pnc,pani,panic"`
	MaxLevel       string   `name:"max" short:"S" default:"none" help:"Filter out logs above specified severity (exclusive) [Not yet implemented]" enum:"n,non,none,t,tra,trac,trce,trace,d,deb,debu,debg,dbug,debug,i,inf,info,information,w,wrn,warn,warning,fail,e,err,erro,errr,error,c,crt,crit,critical,f,fata,fatl,fatal,p,pan,pnc,pani,panic"`
	MinTime        string   `name:"since" short:"t" help:"Filter out logs timestamped before a specific time (or relative time period ago) [Not yet implemented]"`
	MaxTime        string   `name:"before" short:"t" help:"Filter out logs timestamped after a specific time (or relative time period ago) [Not yet implemented]"`
	ExcludedLevels []string `name:"exclude" short:"e" help:"Filter out logs of specified severity (can be specified multiple times) [Not yet implemented]" enum:"n,non,none,t,tra,trac,trce,trace,d,deb,debu,debg,dbug,debug,i,inf,info,information,w,wrn,warn,warning,fail,e,err,erro,errr,error,c,crt,crit,critical,f,fata,fatl,fatal,p,pan,pnc,pani,panic"`
	Paths          []string `arg optional type:"existingfile" help:"File(s) to read logs from. Uses STDIN if unspecified"`
}

func main() {
	kong.Parse(&cli)

	minLevel := parseLevelArg(cli.MinLevel)

	if len(cli.Paths) > 0 {
		for _, path := range cli.Paths {
			printLogsFromFile(path, minLevel)
		}
	} else {
		printLogsFromIO(os.Stdin, minLevel)
	}
}

func printLogsFromFile(path string, minLevel loglevel.Level) {
	if file, err := os.Open(path); err != nil {
		fmt.Printf("ERR: Failed to open file: %s: %v\n", path, err)
		os.Exit(1)
	} else {
		defer file.Close()
		printLogsFromIO(file, minLevel)
	}
}

func printLogsFromIO(r io.Reader, minLevel loglevel.Level) {
	p := logparser.NewIOParser(r)

	printer := NewConsolePrinter(&p, minLevel)
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

func parseLevelArg(s string) loglevel.Level {
	switch strings.ToLower(s) {
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

	return loglevel.LevelUndefined
}
