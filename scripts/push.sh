#!/bin/sh

echo "Formatting files then pushing content to Git repository..."
gofmt -w -s $(dirname "$(pwd)") && git push
