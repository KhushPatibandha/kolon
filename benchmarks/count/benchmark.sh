#!/usr/bin/env bash

set -e

javac code/Count.java
go build -o code/count_go code/count.go
go build -o code/kolon ../../cmd/main.go
cd code/

hyperfine --export-markdown ../results.md --export-json ../results.json --shell=none --warmup 3 --runs 30 './kolon run: count.kol' 'python3 count.py' './count_go' 'java Count' 'node count.js'

rm -f Count.class count_go kolon
