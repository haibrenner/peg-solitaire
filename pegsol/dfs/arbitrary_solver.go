package dfs

import "peg_solitaire/pegsol/board"

func Solve(initial board.CompactState, jumpsPool []*board.CompactJump) []*board.CompactJump {
	numPegs := len(initial.ToInts())
	states := make([]board.CompactState, numPegs)
	jumpsApplied := make([]*board.CompactJump, numPegs-1)
	nextJump := make([]int, numPegs)

	states[0] = initial
	depth := 0

	for depth >= 0 {
		if depth == numPegs-1 {
			return jumpsApplied
		}

		found := false
		for i := nextJump[depth]; i < len(jumpsPool); i++ {
			s := jumpsPool[i]
			if s.IsValidOn(states[depth]) {
				nextJump[depth] = i + 1
				states[depth+1] = s.Apply(states[depth])
				jumpsApplied[depth] = s
				nextJump[depth+1] = 0
				depth++
				found = true
				break
			}
		}

		if !found {
			nextJump[depth] = 0
			depth--
		}
	}

	return nil
}
