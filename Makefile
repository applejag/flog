# Filter multiline logs based on the log's severity
# Copyright (C) 2021  Kalle Jillheden
#
# flog is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# flog is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

build:
	go build
	# The binary should be available at 'flog' (Unix) or 'flog.exe' (Windows)

test:
	go test

release: test
	GOOS=windows GOARCH=amd64 go build -o bin/flog_windows_amd64.exe
	GOOS=windows GOARCH=386 go build -o bin/flog_windows_386.exe
	GOOS=windows GOARCH=arm go build -o bin/flog_windows_arm.exe
	GOOS=linux GOARCH=amd64 go build -o bin/flog_linux_amd64
	GOOS=linux GOARCH=386 go build -o bin/flog_linux_386
	GOOS=linux GOARCH=arm go build -o bin/flog_linux_arm
	GOOS=linux GOARCH=arm64 go build -o bin/flog_linux_arm64

install:
	go install
