#!/bin/bash

# temporary bash script. Will be converted to make file soon

PLATFORM=$1

if [ "$PLATFORM" == "darwin" ]
then
    go build -o main .
else
    GOOS=linux go build -o main .
fi
