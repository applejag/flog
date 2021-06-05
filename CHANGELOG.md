# flog changelog

This project tries to use [SemVer 2.0.0](https://semver.org)

<!--
    Remember to update the version.go file whenever you add a new version.
-->

## v0.4.0 (WIP)

- Added license information embedded into the program. The conditions can be
  retrieved by running `flog --license-c` and the warranty can be solely
  retrieved by running `flog --license-w`. (#15)

- Added shell completions, installed via `--completions-install`. (#17)

- Added JSON log parsing. Reads message, timestamp, and level from JSON
  properties, given that each line is a separate JSON object. (#18)

## v0.3.0 (2020-02-13)

- Added debug logging, enabled by specifying `-vv` or `--verbose=2`. (#13)

- Fixed default maximum severity being set to "Unknown", i.e. the lowest severity
  level, causing all logs except the unknown ones to be omitted. This was caused
  by the value "none" being invalidly interpreted as "Unknown", and is now
  correctly interpreted as "Undefined". (#14)

## v0.2.0 (2021-02-09)

- Added list of severities in the help text. (#12)

- Fixed `-s Error` being treated as an invalid severity. Severities are now
  case insensitive. (#12)

## v0.1.0 (2021-02-03)

- Added basic filtering on generic log lines, supporting .NET logs.

- Added "Omitted 12 Information." messages. (#6)

- Added severity filtering arguments: `-s`, `-S`, `-i`, `-e` (#5)

- Added argument for "quiet mode" (`-q`) to omit the "Omitted..." messages. (#5)

