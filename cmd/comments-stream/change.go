package main

import "github.com/bitterbattles/api/pkg/comments"

type change struct {
	oldComment *comments.Comment
	newComment *comments.Comment
}
