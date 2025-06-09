#!/bin/bash
# This script generates shell completion files for use in the GoReleaser build process.
set -euo pipefail
rm -rf completions
mkdir completions

for shell in bash zsh fish; do
	go run main.go completion "$shell" >"completions/kube-mcp-server.$shell"
done