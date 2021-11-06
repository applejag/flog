# Copyright (C) 2021  Kalle Jillheden
# SPDX-FileCopyrightText: 2021 Kalle Fagerberg
#
# SPDX-License-Identifier: CC0-1.0

build:
	@go build
	# The binary should be available at 'flog' (Unix) or 'flog.exe' (Windows)

test:
	@go test

release: test
	@GOOS=windows GOARCH=amd64 go build -o bin/flog_windows_amd64.exe
	@GOOS=windows GOARCH=386 go build -o bin/flog_windows_386.exe
	@GOOS=windows GOARCH=arm go build -o bin/flog_windows_arm.exe
	@GOOS=linux GOARCH=amd64 go build -o bin/flog_linux_amd64
	@GOOS=linux GOARCH=386 go build -o bin/flog_linux_386
	@GOOS=linux GOARCH=arm go build -o bin/flog_linux_arm
	@GOOS=linux GOARCH=arm64 go build -o bin/flog_linux_arm64

install:
	@go install
	# The binary has been installed to your Go binary path.
	# By default, the binary is found at ~/go/bin/flog
