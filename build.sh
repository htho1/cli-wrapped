#!/usr/bin/sh

OUTPUT_BIN="./build/cli-wrapped"

echo "> go build -o $OUTPUT_BIN main.go"
echo

go build -o $OUTPUT_BIN main.go