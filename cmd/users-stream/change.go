package main

import "github.com/bitterbattles/api/pkg/users"

type change struct {
	oldUser *users.User
	newUser *users.User
}
