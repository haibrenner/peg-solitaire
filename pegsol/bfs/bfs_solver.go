package bfs

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"runtime"

	"peg_solitaire/pegsol/board"
)

type stateInfo struct {
	movesCount  int8
	jumpHistory [64]uint8
}

type indexedJump struct {
	index uint8
	jump  board.CompactJump
}

func Solve(initial board.CompactState, jumps []*board.CompactJump, seedVal uint64, maxStates int) ([][]*board.CompactJump, bool, error) {
	if len(jumps) > 256 {
		return nil, false, fmt.Errorf("number of jumps %d exceeds maximum of 256", len(jumps))
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

	wasPruned := false
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

		// save on memory by pruning the next map if it exceeds the max states limit
		var nextBelowBound map[board.CompactStateWithLastPos]stateInfo
		if maxStates > 0 && len(next) > maxStates {
			slog.Info("BFS pruning next map", "step", step+1, "before", len(next), "limit", maxStates)
			nextBelowBound = pruneMap(next, maxStates)
			slog.Info("BFS pruning done", "step", step+1, "after", len(nextBelowBound))
			wasPruned = true
		} else {
			nextBelowBound = next
		}

		// allow GC to reclaim memory from the current map before the next iteration
		next = nil
		current = nil
		runtime.GC()

		current = nextBelowBound
		nextBelowBound = nil
		runtime.GC()

		slog.Info("BFS step completed", "step", step+1, "states", len(current))
	}

	var bestInfo stateInfo
	bestMoves := int8(-1)
	for _, info := range current {
		if bestMoves == -1 || info.movesCount < bestMoves {
			bestInfo = info
			bestMoves = info.movesCount
		}
	}

	if bestMoves == -1 {
		return nil, wasPruned, nil
	}

	return buildSolution(bestInfo.jumpHistory[:numPegs-1], jumps, numPegs-1), wasPruned, nil
}

func pruneMap(m map[board.CompactStateWithLastPos]stateInfo, maxStates int) map[board.CompactStateWithLastPos]stateInfo {
	pruned := make(map[board.CompactStateWithLastPos]stateInfo, maxStates)
	for k, v := range m {
		pruned[k] = v
		if len(pruned) == maxStates {
			break
		}
	}
	return pruned
}

func buildSolution(history []uint8, jumps []*board.CompactJump, numJumps int) [][]*board.CompactJump {
	var result [][]*board.CompactJump
	var currentMove []*board.CompactJump
	for i := range numJumps {
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
