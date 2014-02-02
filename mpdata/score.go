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
	
	score *= 0.1 + float32(math.Tanh(float64(dist) * 0.2))
	s.tagScores[tag] = score
}

func (s *Scorer) ScoreMeal(favourite bool, csd int, tags []string) (score float32) {
	score = 1.35 - (2.8 / float32(csd + 1))
	
	for _, tag := range tags {
		tagScore, ok := s.tagScores[tag]
		if ok {
			score *= tagScore
		}
	}
	
	if favourite {
		score *= 2.0
	}
	
	return score
}
