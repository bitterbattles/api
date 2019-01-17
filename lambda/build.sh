#!/usr/bin/env bash

export GOOS=linux
export GOARCH=amd64

go build battles_get/battles_get.go
zip ./bin/battles_get.zip battles_get
rm battles_get

go build battles_post/battles_post.go
zip ./bin/battles_post.zip battles_post
rm battles_post

go build votes_post/votes_post.go
zip ./bin/votes_post.zip votes_post
rm votes_post
