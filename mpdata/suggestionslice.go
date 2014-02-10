package mpdata

type SuggestionSlice []*Suggestion

func (ss SuggestionSlice) Len() (n int) {
	return len(ss)
}

func (ss SuggestionSlice) Less(i int, j int) (less bool) {
	return ss[i].Score > ss[j].Score
}

func (ss SuggestionSlice) Swap(i int, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

