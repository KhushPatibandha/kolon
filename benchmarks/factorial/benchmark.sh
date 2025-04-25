#!/usr/bin/env bash

set -e

javac code/Factorial.java
go build -o code/factorial_go code/factorial.go
go build -o code/kolon ../../cmd/main.go
cd code/

hyperfine --shell=none --warmup 3 --runs 20 './kolon run: factorial.kol' 'python3 factorial.py' './factorial_go' 'java Factorial'

rm -f Factorial.class factorial_go kolon
