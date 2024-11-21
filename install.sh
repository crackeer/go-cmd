#!/bin/sh

output="${GOPATH}/bin/got"


go build -o="${output}" main.go

echo "Installed at ${output}"