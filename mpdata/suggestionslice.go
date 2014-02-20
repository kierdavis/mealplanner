package mpdata

// SuggestionSlice is an implementation of 'sort.Interface', allowing a
// list of suggestions to be sorted by their score.
type SuggestionSlice []*Suggestion

// Len returns the number of suggestions in the list.
func (ss SuggestionSlice) Len() (n int) {
	return len(ss)
}

// Less returns whether the suggestion indexed by 'i' has a higher score than
// the suggestion indexed by 'j' and so should be placed closer to the start of
// the list.
func (ss SuggestionSlice) Less(i int, j int) (less bool) {
	return ss[i].Score > ss[j].Score
}

// Swap exchanges the suggestions indexed by 'i' and 'j'.
func (ss SuggestionSlice) Swap(i int, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}
