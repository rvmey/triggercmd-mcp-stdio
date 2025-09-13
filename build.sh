#!/bin/bash

mkdir -p binaries

oses=("windows" "darwin" "linux")
arches=("amd64" "386" "arm" "arm64")

for os in "${oses[@]}"; do
  for arch in "${arches[@]}"; do
    if [[ "$arch" == "arm" && "$os" != "linux" ]]; then
      continue
    fi
    if [[ "$os" == "windows" ]]; then
      outfile="binaries/triggercmd-mcp-$os-$arch.exe"
    else
      outfile="binaries/triggercmd-mcp-$os-$arch"
    fi
    rm -f "$outfile"
    GOOS="$os" GOARCH="$arch" go build -v -o "$outfile"
  done
done