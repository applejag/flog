# flog changelog

This project tries to use [SemVer 2.0.0](https://semver.org)

## v0.3.0 (WIP)

- Added debug logging, enabled by specifying `-vv` or `--verbose=2`. (#13)

## v0.2.0 (2021-02-09)

- Added list of severities in the help text. (#12)

- Fixed `-s Error` being treated as an invalid severity. Severities are now
  case insensitive. (#12)

## v0.1.0 (2021-02-03)

- Added basic filtering on generic log lines, supporting .NET logs.

- Added "Omitted 12 Information." messages. (#6)

- Added severity filtering arguments: `-s`, `-S`, `-i`, `-e` (#5)

- Added argument for "quiet mode" (`-q`) to omit the "Omitted..." messages. (#5)

