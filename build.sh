#!/usr/bin/env bash
if [ ! -d ./bin ]; then
  mkdir -p ./bin;
fi
if [ ! -d ./packages ]; then
  mkdir -p ./packages;
fi
read -p i'nput version name:' version_name
#read version_name

GOOS=windows GOARCH=amd64 go build -o bin/game_windows.exe .
GOOS=linux GOARCH=amd64 go build -o bin/game_linux.exe .