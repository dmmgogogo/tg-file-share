#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o tg-file-share main.go
zip tg-file-share.zip tg-file-share