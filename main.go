package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
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

//--min, -s <severity>      Filter out everything below a specific severity (exclusive)
//--max, -S <severity>      Filter out everthing above specific severity (exclusive)
//--since, -t <time>        Filter out everything before a specific time (or relative time period ago)
//--before, -T <time>       Filter out everything after a specific time (or relative time period ago)
//--exclude, -e <severity>  Filter out a specific severity

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

func printLogsFromFile(path string, minLevel logparser.Level) {
	if file, err := os.Open(path); err != nil {
		fmt.Printf("ERR: Failed to open file: %s\n", path)
		os.Exit(1)
	} else {
		defer file.Close()
		printLogsFromIO(file, minLevel)
	}
}

func printLogsFromIO(r io.Reader, minLevel logparser.Level) {
	p := logparser.NewIOParser(r)

	printer := NewConsolePrinter(&p, minLevel)

	for printer.Next() {
	}
}

func parseLevelArg(s string) logparser.Level {
	switch strings.ToLower(s) {
	case "t", "tra", "trac", "trce", "trace":
		return logparser.LevelTrace

	case "d", "deb", "dbg", "debu", "debg", "dbug", "debug":
		return logparser.LevelDebug

	case "i", "inf", "info", "information":
		return logparser.LevelInformation

	case "w", "wrn", "warn", "warning":
		return logparser.LevelWarning

	case "fail", "e", "err", "erro", "errr", "error":
		return logparser.LevelError

	case "c", "crt", "crit", "critical":
		return logparser.LevelCritical

	case "f", "fata", "fatl", "fatal":
		return logparser.LevelFatal

	case "p", "pan", "pnc", "pani", "panic":
		return logparser.LevelPanic
	}

	return logparser.LevelUndefined
}
