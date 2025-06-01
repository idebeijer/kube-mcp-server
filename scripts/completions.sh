#!/bin/sh
# This script generates shell completion files for use in the GoReleaser build process.
set -e
rm -rf completions
mkdir completions

for sh in bash zsh fish; do
	go run main.go completion "$sh" >"completions/kubert.$sh"
done