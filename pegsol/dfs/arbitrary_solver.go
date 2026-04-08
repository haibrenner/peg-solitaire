package dfs

import (
	"math/rand/v2"
	"peg_solitaire/pegsol/board"
)

func Solve(initial board.CompactState, jumpsPool []*board.CompactJump, seedVal uint64) []*board.CompactJump {
	pcg := rand.NewPCG(seedVal, seedVal+1)
	r := rand.New(pcg)

	numPegs := len(initial.ToInts())
	numJumps := len(jumpsPool)

	// Pre-shuffle the jumps pool for each depth to ensure different solution paths across runs with the same seed.
	// Each depth has its own ordering, giving more varied exploration. Arguably, this adds more randomness to the search, and may result in faster solutions on average.
	shuffledPools := make([][]*board.CompactJump, numPegs-1)
	for d := range shuffledPools {
		copy_ := make([]*board.CompactJump, numJumps)
		copy(copy_, jumpsPool)
		r.Shuffle(len(copy_), func(i, j int) {
			copy_[i], copy_[j] = copy_[j], copy_[i]
		})
		shuffledPools[d] = copy_
	}

	states := make([]board.CompactState, numPegs)
	jumpsApplied := make([]*board.CompactJump, numPegs-1)
	nextJump := make([]int, numPegs)
	memoizationExhaustedStates := make(map[board.CompactState]bool)

	states[0] = initial
	depth := 0

	for depth >= 0 {
		if depth == numPegs-1 {
			return jumpsApplied
		}

		found := false
		for i := nextJump[depth]; i < numJumps; i++ {
			s := shuffledPools[depth][i]
			if s.IsValidOn(states[depth]) {
				states[depth+1] = s.Apply(states[depth])
				if memoizationExhaustedStates[states[depth+1]] {
					continue
				}
				nextJump[depth] = i + 1
				jumpsApplied[depth] = s
				nextJump[depth+1] = 0
				depth++
				found = true
				break
			}
		}

		if !found {
			memoizationExhaustedStates[states[depth]] = true
			nextJump[depth] = 0
			depth--
		}
	}

	return nil
}
