package dfs

import "peg_solitaire/pegsol/board"

func Solve(initial board.CompactState, stepsPool []*board.CompactStep) []*board.CompactStep {
	numPegs := len(initial.ToInts())
	states := make([]board.CompactState, numPegs)
	stepsApplied := make([]*board.CompactStep, numPegs-1)
	nextStep := make([]int, numPegs)

	states[0] = initial
	depth := 0

	for depth >= 0 {
		if depth == numPegs-1 {
			return stepsApplied
		}

		found := false
		for i := nextStep[depth]; i < len(stepsPool); i++ {
			s := stepsPool[i]
			if s.IsValidOn(states[depth]) {
				nextStep[depth] = i + 1
				states[depth+1] = s.Apply(states[depth])
				stepsApplied[depth] = s
				nextStep[depth+1] = 0
				depth++
				found = true
				break
			}
		}

		if !found {
			nextStep[depth] = 0
			depth--
		}
	}

	return nil
}
