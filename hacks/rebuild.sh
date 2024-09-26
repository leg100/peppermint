#!/usr/bin/env bash

while true; do
    go build -o _build/moto && pkill -f '_build/moto'
    inotifywait -e attrib $(find . -name '*.go') || exit
done
