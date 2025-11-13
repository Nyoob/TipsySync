#!/bin/bash


# do: sudo apt-get install mingw-w64
#
export CC=x86_64-w64-mingw32-gcc
export CGO_ENABLED=1
export GOOS=windows
export GOARCH=amd64


if [ "$1" = "debug" ]; then
  echo "Compiling debug version"
  wails build --platform windows/amd64 -debug -windowsconsole -ldflags "-X 'main.BuildNumber=debug_$(date +%s)'"
elif [ "$1" = "prod" ]; then
  echo "Compiling prod version"
  wails build --platform windows/amd64 -ldflags "-X 'main.BuildNumber=$(date +%s)'"
else
  echo "Usage: $0 [debug|prod]"
fi
