# flog

Filter logs on their serverity, with automatic detection of log formats

## Sample usage

```sh
# Sample Go app using sirupsen/logrus
$ go run main.go
TRAC[0000] A walrus appears                              animal=walrus
DEBU[0000] A walrus appears                              animal=walrus
INFO[0000] A walrus appears                              animal=walrus
WARN[0000] A walrus appears                              animal=walrus
ERRO[0000] A walrus appears                              animal=walrus
FATA[0000] A walrus appears                              animal=walrus

$ go run main.go | flog -s warn
WARN[0000] A walrus appears                              animal=walrus
ERRO[0000] A walrus appears                              animal=walrus
FATA[0000] A walrus appears                              animal=walrus
```

And yes, this includes multiline logs, such as those pesky .NET logs:

```sh
# Sample .NET app using Microsoft.Extensions.Logging.Console
$ dotnet run
trac: Program[0]
      Sample log
dbug: Program[0]
      Sample log
info: Program[0]
      Sample log
warn: Program[0]
      Sample log
fail: Program[0]
      Sample log

$ dotnet run | flog -s warn
warn: Program[0]
      Sample log
fail: Program[0]
      Sample log
```

## Command-line interface

```sh
~/dev/flog main ❯ flog --help
Usage: flog [<paths> ...]

Arguments:
  [<paths> ...]    File(s) to read logs from. Uses STDIN if unspecified

Flags:
  -h, --help                   Show context-sensitive help.
  -s, --min="info"             Filter out logs below specified severity (exclusive)
  -S, --max="none"             Filter out logs above specified severity (exclusive) [Not yet implemented]
  -t, --since=STRING           Filter out logs timestamped before a specific time (or relative time period ago) [Not yet implemented]
  -t, --before=STRING          Filter out logs timestamped after a specific time (or relative time period ago) [Not yet implemented]
  -e, --exclude=EXCLUDE,...    Filter out logs of specified severity (can be specified multiple times) [Not yet implemented]
```

## Installation

1. Install Go

2. Run the following (outside of a Go project)

   ```sh
   go get github.com/jilleJr/flog
   ```

## Note

This project is under prototype phase.

You are welcome to try it out or participate in the design discussions:
<https://github.com/jilleJr/flog/discussions/categories/ideas>

