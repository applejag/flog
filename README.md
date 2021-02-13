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
flog: Omitted 1 Trace, 1 Debug, 1 Information.
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
flog: Omitted 1 Trace, 1 Debug, 1 Information.
warn: Program[0]
      Sample log
fail: Program[0]
      Sample log
```

## Command-line interface

```
sh
$ flog --help
Usage: flog [<paths> ...]

Arguments:
  [<paths> ...]    File(s) to read logs from. Uses STDIN if unspecified

Flags:
  -h, --help                   Show context-sensitive help.
  -s, --min=info               Omit logs below specified severity (exclusive)
  -S, --max=none               Omit logs above specified severity (exclusive)
  -t, --since=STRING           Omit logs timestamped before a specific time (or relative time period ago) [Not yet
                               implemented]
  -t, --before=STRING          Omit logs timestamped after a specific time (or relative time period ago) [Not yet
                               implemented]
  -e, --exclude=EXCLUDE,...    Omit logs of specified severity (can be specified multiple times)
  -i, --include=INCLUDE,...    Omit logs of severity not specified with this flag (can be specified multiple times)
  -q, --quiet                  Omit the 'omitted logs' messages. Shorthand for --verbose=0.
  -v, --verbose=1              Enable verbose output (can be specified up to 2 times, ex: --verbose=2 or -vv)
      --version                Show the version and then exit.

Severities:
  Undefined      0, n, nil, null, none, unde, undefined
  Unknown        1, ?, u, ukwn, unkn, unknown
  Trace          2, t, tra, trac, trce, trace
  Debug          3, d, deb, debu, debg, dbug, debug
  Information    4, i, inf, info, information
  Warning        5, w, wrn, warn, warning
  Error          6, e, err, erro, errr, error
  Critical       7, c, crt, crit, critical
  Fatal          8, f, fata, fatl, fatal
  Panic          9, p, pan, pnc, pani, panic
```

## Installation

1. Head over to the latest release
   <https://github.com/jilleJr/flog/releases/latest>

2. Down in the "Assets", find the one that matches your OS and architecture,
   and download it.

3. Save it to a location you have added to your `$PATH`, for example
   `~/.local/bin`

4. Make sure to just rename it to `flog` (or `flog.exe` for Windows users).

5. Done!


If your OS or architecture is missing from the list, then
[create an issue on it](https://github.com/jilleJr/flog/issues/new/choose) and
I'll make sure to build it to your architecture when I release a new build.

## Building from source

1. Install Go

2. Run the following (outside of a Go project)

   ```sh
   # The -u flag updates the package if you've already the package installed.
   go get -u github.com/jilleJr/flog
   ```

## Main use case

Searching through logs in Kubernetes. For example, to find any error logs in
Kubernetes logs you'd usually have to use `grep` with some extreamly high
context so that you get a glimpse of what the error was about.

With `flog` you can just run:

```sh
kubectl logs name-of-my-pod | flog -s err
```

And there you have it.

## Note

This project is under prototype phase.

You are welcome to try it out or participate in the design discussions:
<https://github.com/jilleJr/flog/discussions/categories/ideas>

