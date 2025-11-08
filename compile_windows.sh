#!/bin/bash


# do: sudo apt-get install mingw-w64
#
export CC=x86_64-w64-mingw32-gcc
export CGO_ENABLED=1
export GOOS=windows
export GOARCH=amd64


if [ "$1" = "dev" ]; then
  echo "Compiling dev version"
  wails build --platform windows/amd64 -debug -windowsconsole
elif [ "$1" = "prod" ]; then
  echo "Compiling prod version"
  wails build --platform windows/amd64
else
  echo "Usage: $0 [dev|prod]"
fi
