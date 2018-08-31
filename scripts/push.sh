#!/bin/sh

echo "Formatting files"
gofmt -w -s $(dirname "$(pwd)") && git push

