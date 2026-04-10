package bfs

import (
	"math/rand/v2"
	"peg_solitaire/pegsol/board"
)

type stateInfo struct {
	movesCount int8
	prev       board.CompactStateWithLastPos
}

func Solve(initial board.CompactState, jumps []*board.CompactJump, seedVal uint64) []board.CompactStateWithLastPos {
	pcg := rand.NewPCG(seedVal, seedVal+1)
	r := rand.New(pcg)

	jumpsWork := make([]*board.CompactJump, len(jumps))
	copy(jumpsWork, jumps)
	numPegs := len(initial.ToInts())

	initialState := board.CompactStateWithLastPos{
		CompactState: initial,
		LastPegPos:   -1,
	}

	levels := make([]map[board.CompactStateWithLastPos]stateInfo, numPegs)
	levels[0] = map[board.CompactStateWithLastPos]stateInfo{
		initialState: {movesCount: 0, prev: board.CompactStateWithLastPos{}},
	}

	for step := 0; step < numPegs-1; step++ {
		r.Shuffle(len(jumpsWork), func(i, j int) {
			jumpsWork[i], jumpsWork[j] = jumpsWork[j], jumpsWork[i]
		})
		next := make(map[board.CompactStateWithLastPos]stateInfo)

		currLevel := levels[step]
		for state, info := range currLevel {
			for _, jump := range jumpsWork {
				if !jump.IsValidOn(state.CompactState) {
					continue
				}
				newCompact := jump.Apply(state.CompactState)
				newMovesCount := info.movesCount
				if state.LastPegPos != jump.StartPosition {
					newMovesCount++
				}
				newState := board.CompactStateWithLastPos{
					CompactState: newCompact,
					LastPegPos:   jump.EndPosition,
				}
				if existing, found := next[newState]; !found || newMovesCount < existing.movesCount {
					next[newState] = stateInfo{movesCount: newMovesCount, prev: state}
				}
			}
		}

		levels[step+1] = next
	}

	// find ending state with least moves
	lastLevel := levels[numPegs-1]
	var bestState board.CompactStateWithLastPos
	bestMoves := int8(-1)
	for state, info := range lastLevel {
		if bestMoves == -1 || info.movesCount < bestMoves {
			bestState = state
			bestMoves = info.movesCount
		}
	}

	if bestMoves == -1 {
		return nil
	}

	// reconstruct path by following prev values through levels
	path := make([]board.CompactStateWithLastPos, numPegs)
	path[numPegs-1] = bestState
	prev := lastLevel[bestState].prev
	for i := numPegs - 2; i >= 0; i-- {
		path[i] = prev
		if i > 0 {
			prev = levels[i][prev].prev
		}
	}

	return path
}
