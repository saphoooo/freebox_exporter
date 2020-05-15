#!/bin/bash

set -eu

ARCH=${1:-arm}
GOARCH=$ARCH go build -ldflags "-s -w"
ls -lh freebox_exporter
file freebox_exporter
