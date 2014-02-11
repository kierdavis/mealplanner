package mpdata

import (
	"math"
)

type Scorer struct {
	tagScores map[string]float32
}

func NewScorer() (s *Scorer) {
	return &Scorer{
		tagScores: make(map[string]float32),
	}
}

func (s *Scorer) AddTagScore(tag string, dist int) {
	score, ok := s.tagScores[tag]
	if !ok {
		score = 1.0 // the default
	}

	score *= 0.1 + float32(math.Tanh(float64(dist)*0.2))
	s.tagScores[tag] = score
}

func (s *Scorer) ScoreSuggestion(sugg *Suggestion) {
	score := float32(1)

	if sugg.CSD < 0 {
		score *= 1.6
	} else {
		score *= 1.45 - (2.8 / float32(sugg.CSD+1))
	}

	for _, tag := range sugg.MT.Tags {
		tagScore, ok := s.tagScores[tag]
		if ok {
			score *= tagScore
		}
	}

	if sugg.MT.Meal.Favourite {
		score *= 2.0
	}

	sugg.Score = score
}
