#!/usr/bin/env bash

set -e

javac code/Factorial.java
go build -o code/factorial_go code/factorial.go
go build -o code/kolon ../../cmd/main.go
cd code/

hyperfine --export-markdown ../results.md --export-json ../results.json --shell=none --warmup 3 --runs 30 './kolon run: factorial.kol' 'python3 factorial.py' './factorial_go' 'java Factorial' 'node factorial.js'

rm -f Factorial.class factorial_go kolon
