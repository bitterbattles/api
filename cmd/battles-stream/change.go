package main

import "github.com/bitterbattles/api/pkg/battles"

type change struct {
	oldBattle *battles.Battle
	newBattle *battles.Battle
}
