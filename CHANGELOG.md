# flog changelog

This project tries to use [SemVer 2.0.0](https://semver.org)

# v0.2.0 (2021-02-09)

- Added list of severities in the help text.

- Fixed `-s Error` being treated as an invalid severity. Severities are now
  case insensitive.

# v0.1.0 (2021-02-03)

- Added basic filtering on generic log lines, supporting .NET logs.

- Added "Omitted 12 Information." messages.

- Added severity filtering arguments: `-s`, `-S`, `-i`, `-e`

- Added argument for "quiet mode" (`-q`) to omit the "Omitted..." messages.

