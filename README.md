# flog

Filter logs on their serverity, with automatic detection of log formats

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


This project is under prototype phase.

You are welcome to try it out or participate in the design discussions:
<https://github.com/jilleJr/flog/discussions/categories/ideas>

