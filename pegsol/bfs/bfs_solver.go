package bfs

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"peg_solitaire/pegsol/board"
	"runtime"
)

type stateInfo struct {
	movesCount  int8
	jumpHistory [64]uint8
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

	current := map[board.CompactStateWithLastPos]stateInfo{
		initialState: {},
	}

	for step := 0; step < numPegs-1; step++ {
		r.Shuffle(len(jumpsWork), func(i, j int) {
			jumpsWork[i], jumpsWork[j] = jumpsWork[j], jumpsWork[i]
		})
		next := make(map[board.CompactStateWithLastPos]stateInfo)

		for state, info := range current {
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
					newInfo := info
					newInfo.movesCount = newMovesCount
					newInfo.jumpHistory[step] = ij.index
					next[newState] = newInfo
				}
			}
		}

		current = nil
		runtime.GC()
		current = next
		runtime.GC()
		slog.Info("BFS step completed", "step", step+1, "states", len(next))
	}

	// find ending state with least moves
	var bestInfo stateInfo
	bestMoves := int8(-1)
	for _, info := range current {
		if bestMoves == -1 || info.movesCount < bestMoves {
			bestInfo = info
			bestMoves = info.movesCount
		}
	}

	if bestMoves == -1 {
		return nil, nil
	}

	return buildSolution(bestInfo.jumpHistory[:numPegs-1], jumps, numPegs-1), nil
}

func buildSolution(history []uint8, jumps []*board.CompactJump, numJumps int) [][]*board.CompactJump {
	var result [][]*board.CompactJump
	var currentMove []*board.CompactJump
	for i := 0; i < numJumps; i++ {
		jump := jumps[history[i]]
		if i > 0 && jump.StartPosition != jumps[history[i-1]].EndPosition {
			result = append(result, currentMove)
			currentMove = nil
		}
		currentMove = append(currentMove, jump)
	}
	result = append(result, currentMove)
	return result
}
