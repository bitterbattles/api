package main

type change struct {
	createdOnChanged bool
	newCreatedOn     int64
	votesChanged     bool
	newVotesFor      int
	newVotesAgainst  int
}
