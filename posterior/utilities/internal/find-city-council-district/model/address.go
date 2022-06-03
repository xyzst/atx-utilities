package model

import "errors"

type Address struct {
	Candidates []Candidate
}

func (receiver Address) GetLikelyCandidate() (*Candidate, error) {
	if len(receiver.Candidates) == 0 {
		return nil, errors.New("no known address matches")
	}

	var highest = 0
	for current := 1; current < len(receiver.Candidates); current++ {
		if receiver.Candidates[current].Score > receiver.Candidates[highest].Score {
			highest = current
		}
	}

	return &receiver.Candidates[highest], nil
}