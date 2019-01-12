#!/usr/bin/env bash

export GOOS=linux
export GOARCH=amd64


go build battles_get.go
zip ./bin/battles_get.zip battles_get
rm battles_get

go build battles_post.go
zip ./bin/battles_post.zip battles_post
rm battles_post