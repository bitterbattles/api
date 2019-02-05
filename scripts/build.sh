#!/usr/bin/env bash

rm -f battles-get.zip
GOARCH=amd64 GOOS=linux go build -o battles-get ../cmd/battles-get
zip battles-get.zip battles-get
rm battles-get

rm -f battles-post.zip
GOARCH=amd64 GOOS=linux go build -o battles-post ../cmd/battles-post
zip battles-post.zip battles-post
rm battles-post

rm -f battles-stream.zip
GOARCH=amd64 GOOS=linux go build -o battles-stream ../cmd/battles-stream
zip battles-stream.zip battles-stream
rm battles-stream

rm -f votes-post.zip
GOARCH=amd64 GOOS=linux go build -o votes-post ../cmd/votes-post
zip votes-post.zip votes-post
rm votes-post

rm -f votes-stream.zip
GOARCH=amd64 GOOS=linux go build -o votes-stream ../cmd/votes-stream
zip votes-stream votes-stream
rm votes-stream