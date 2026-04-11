package bfs

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"peg_solitaire/pegsol/board"
)

type stateInfo struct {
	movesCount int8
	prev       board.CompactStateWithLastPos
}

type indexedJump struct {
	index uint8
	jump  board.CompactJump
}

func Solve(initial board.CompactState, jumps []*board.CompactJump, seedVal uint64) ([][]*board.CompactJump, error) {
	if len(jumps) > 256 {
		return nil, fmt.Errorf("number of jumps %d exceeds maximum of 256", len(jumps))
	}

	pcg := rand.NewPCG(seedVal, seedVal+1)
	r := rand.New(pcg)

	jumpsWork := make([]indexedJump, len(jumps))
	for i, j := range jumps {
		jumpsWork[i] = indexedJump{index: uint8(i), jump: *j}
	}
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
			for _, ij := range jumpsWork {
				if !ij.jump.IsValidOn(state.CompactState) {
					continue
				}
				newCompact := ij.jump.Apply(state.CompactState)
				newMovesCount := info.movesCount
				if state.LastPegPos != ij.jump.StartPosition {
					newMovesCount++
				}
				newState := board.CompactStateWithLastPos{
					CompactState: newCompact,
					LastPegPos:   ij.jump.EndPosition,
				}
				if existing, found := next[newState]; !found || newMovesCount < existing.movesCount {
					next[newState] = stateInfo{movesCount: newMovesCount, prev: state}
				}
			}
		}

		levels[step+1] = next
		slog.Info("BFS step completed", "step", step+1, "states", len(next))
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
		return nil, nil
	}

	// reconstruct flat path by following prev values through levels
	flatPath := make([]board.CompactStateWithLastPos, numPegs)
	flatPath[numPegs-1] = bestState
	prev := lastLevel[bestState].prev
	for i := numPegs - 2; i >= 0; i-- {
		flatPath[i] = prev
		if i > 0 {
			prev = levels[i][prev].prev
		}
	}

	grouped := groupByMoves(flatPath, levels)
	return statesGroupToJumpsGroup(grouped, jumps)
}

func statesGroupToJumpsGroup(grouped [][]board.CompactStateWithLastPos, jumps []*board.CompactJump) ([][]*board.CompactJump, error) {
	result := make([][]*board.CompactJump, len(grouped))
	for i, move := range grouped {
		moveJumps := make([]*board.CompactJump, len(move)-1)
		for j := 0; j < len(move)-1; j++ {
			jump, err := findJump(move[j].CompactState, move[j+1].CompactState, jumps)
			if err != nil {
				return nil, err
			}
			moveJumps[j] = jump
		}
		result[i] = moveJumps
	}
	return result, nil
}

func groupByMoves(flatPath []board.CompactStateWithLastPos, levels []map[board.CompactStateWithLastPos]stateInfo) [][]board.CompactStateWithLastPos {
	var result [][]board.CompactStateWithLastPos
	var currentMove []board.CompactStateWithLastPos
	currentMove = append(currentMove, flatPath[0])
	currentMoveCount := int8(1)
	for i := 1; i < len(flatPath); i++ {
		info := levels[i][flatPath[i]]
		if info.movesCount > currentMoveCount {
			result = append(result, currentMove)
			currentMove = []board.CompactStateWithLastPos{flatPath[i-1]}
			currentMoveCount = info.movesCount
		}
		currentMove = append(currentMove, flatPath[i])
	}
	result = append(result, currentMove)
	return result
}

func findJump(from, to board.CompactState, jumps []*board.CompactJump) (*board.CompactJump, error) {
	for _, jump := range jumps {
		if jump.IsValidOn(from) && jump.Apply(from) == to {
			return jump, nil
		}
	}
	return nil, fmt.Errorf("no jump found between states %v and %v", from, to)
}
