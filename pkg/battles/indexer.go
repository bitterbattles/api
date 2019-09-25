package battles

import (
	"fmt"
	"math"

	"github.com/bitterbattles/api/pkg/index"
	"github.com/bitterbattles/api/pkg/time"
)

const globalKeyPattern = "battleIds:%s"
const authorKeyPattern = "battleIds:forAuthor:%s:%s"

// Indexer is used to index Battles in a variety of ways
type Indexer struct {
	repository index.RepositoryInterface
}

// NewIndexer creates a new indexer instance
func NewIndexer(repository index.RepositoryInterface) *Indexer {
	return &Indexer{repository}
}

// Add adds a new Battle to the relevant indexes
func (indexer *Indexer) Add(battle *Battle) error {
	var key string
	var err error
	member := battle.ID
	score := float64(battle.CreatedOn)
	key = fmt.Sprintf(globalKeyPattern, RecentSort)
	err = indexer.repository.SetScore(key, member, score)
	if err != nil {
		return err
	}
	key = fmt.Sprintf(authorKeyPattern, battle.UserID, RecentSort)
	err = indexer.repository.SetScore(key, member, score)
	if err != nil {
		return err
	}
	err = indexer.UpdateVotes(battle)
	if err != nil {
		return err
	}
	return nil
}

// GetGlobal gets a range of global Battle IDs
func (indexer *Indexer) GetGlobal(sort string, page int, pageSize int) ([]string, error) {
	key := fmt.Sprintf(globalKeyPattern, sort)
	offset, limit := indexer.toRange(page, pageSize)
	return indexer.repository.GetRange(key, offset, limit)
}

// GetByAuthor gets a range of Battle IDs authored by the given user ID
func (indexer *Indexer) GetByAuthor(userID string, sort string, page int, pageSize int) ([]string, error) {
	key := fmt.Sprintf(authorKeyPattern, userID, sort)
	offset, limit := indexer.toRange(page, pageSize)
	return indexer.repository.GetRange(key, offset, limit)
}

// IsBattleAuthor determines whether or not the given user authored the given battle
func (indexer *Indexer) IsBattleAuthor(userID string, battleID string) (bool, error) {
	key := fmt.Sprintf(authorKeyPattern, userID, RecentSort)
	return indexer.repository.HasMember(key, battleID)
}

// UpdateVotes updates indexes related to the Battle's votes
func (indexer *Indexer) UpdateVotes(battle *Battle) error {
	var key string
	var score float64
	var err error
	member := battle.ID
	key = fmt.Sprintf(globalKeyPattern, PopularSort)
	score = indexer.calculatePopularity(battle)
	err = indexer.repository.SetScore(key, member, score)
	if err != nil {
		return err
	}
	key = fmt.Sprintf(authorKeyPattern, battle.UserID, PopularSort)
	err = indexer.repository.SetScore(key, member, score)
	if err != nil {
		return err
	}
	key = fmt.Sprintf(globalKeyPattern, ControversialSort)
	score = indexer.calculateControversy(battle)
	err = indexer.repository.SetScore(key, member, score)
	if err != nil {
		return err
	}
	key = fmt.Sprintf(authorKeyPattern, battle.UserID, ControversialSort)
	err = indexer.repository.SetScore(key, member, score)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes the given Battle from its indexes
func (indexer *Indexer) Delete(battle *Battle) error {
	var key string
	var err error
	member := battle.ID
	author := battle.UserID
	key = fmt.Sprintf(globalKeyPattern, RecentSort)
	err = indexer.repository.Remove(key, member)
	if err != nil {
		return err
	}
	key = fmt.Sprintf(authorKeyPattern, author, RecentSort)
	err = indexer.repository.Remove(key, member)
	if err != nil {
		return err
	}
	key = fmt.Sprintf(globalKeyPattern, PopularSort)
	err = indexer.repository.Remove(key, member)
	if err != nil {
		return err
	}
	key = fmt.Sprintf(authorKeyPattern, author, PopularSort)
	err = indexer.repository.Remove(key, member)
	if err != nil {
		return err
	}
	key = fmt.Sprintf(globalKeyPattern, ControversialSort)
	err = indexer.repository.Remove(key, member)
	if err != nil {
		return err
	}
	key = fmt.Sprintf(authorKeyPattern, author, ControversialSort)
	err = indexer.repository.Remove(key, member)
	if err != nil {
		return err
	}
	return nil
}

func (indexer *Indexer) toRange(page int, pageSize int) (int, int) {
	return (page - 1) * pageSize, pageSize
}

func (indexer *Indexer) calculatePopularity(battle *Battle) float64 {
	totalVotes := float64(battle.VotesFor + battle.VotesAgainst)
	return indexer.getRecencyWeight() + totalVotes
}

func (indexer *Indexer) calculateControversy(battle *Battle) float64 {
	totalVotes := float64(battle.VotesFor + battle.VotesAgainst)
	voteDifference := math.Abs(float64(battle.VotesFor - battle.VotesAgainst))
	return indexer.getRecencyWeight() + totalVotes - voteDifference
}

func (indexer *Indexer) getRecencyWeight() float64 {
	return float64(time.NowUnix() / 86400)
}
