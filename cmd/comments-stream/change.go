package main

import "github.com/bitterbattles/api/pkg/comments"

type changedComment struct {
	oldComment *comments.Comment
	newComment *comments.Comment
}

type changedBattle struct {
	deltaComments int
}
