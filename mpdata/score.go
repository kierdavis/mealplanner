package mpdata

func CalculateScore(favourite bool, closestServingDistance int) (score float32) {
	score = 1.35 - (2.8 / float32(closestServingDistance + 1))
	
	if favourite {
		score *= 2.0
	}
	
	return score
}