package battles

import "math"

// Scorer scores Battles
type Scorer struct {
}

// NewScorer creates a new Scorer instance
func NewScorer() *Scorer {
	return &Scorer{}
}

// ScorePopularity scores a Battle's popularity
func (scorer *Scorer) ScorePopularity(battle *Battle) int64 {
	totalVotes := int64(battle.VotesFor + battle.VotesAgainst)
	totalComments := int64(battle.Comments)
	return scorer.getRecencyWeight(battle) + totalVotes + totalComments
}

// ScoreControversy scores a Battle's controversy
func (scorer *Scorer) ScoreControversy(battle *Battle) int64 {
	totalVotes := int64(battle.VotesFor + battle.VotesAgainst)
	totalComments := int64(battle.Comments)
	voteDifference := int64(math.Abs(float64(battle.VotesFor - battle.VotesAgainst)))
	return scorer.getRecencyWeight(battle) + totalComments + totalVotes - voteDifference
}

func (scorer *Scorer) getRecencyWeight(battle *Battle) int64 {
	return battle.CreatedOn / 86400
}
