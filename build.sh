#!/bin/bash

go get -v

GOARM=7 GOARCH=arm GOOS=linux go build

scp motion-bot main.go pi@pi.local:/home/pi/go/src/github.com/exilesprx/motion-bot