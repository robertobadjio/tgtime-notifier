#!/bin/bash

# Optionally, set default values
# removeServer="default value for 127.0.0.1"

. deploy/deploy.config

env GOOS=linux GOARCH=amd64 go build -o build/officetime-tracker -v cmd/officetime/main.go
ssh root@$removeServer 'systemctl stop officetime-tracker.service && rm /var/officetime/officetime-tracker'
scp build/officetime-tracker root@$removeServer:/var/officetime/
ssh root@$removeServer 'systemctl start officetime-tracker.service'

exit 0