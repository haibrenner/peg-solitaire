package bfs

import (
	"testing"

	"peg_solitaire/pegsol/board"
	"peg_solitaire/pegsol/matrixstate"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const inputsDir = "../inputs/"

func prepareForSolve(t *testing.T, filename string) (board.CompactState, []*board.CompactJump) {
	t.Helper()
	ms, err := matrixstate.ReadInput(inputsDir + filename)
	require.NoError(t, err)

	b, err := board.NewBoard(ms)
	require.NoError(t, err)

	jumps, err := b.TranslateAllCoordJumpsToCompact()
	require.NoError(t, err)

	initial, err := b.TranslateMatrixToCompactState(ms)
	require.NoError(t, err)

	return initial, jumps
}

func TestSolve_TooSimple(t *testing.T) {
	initial, jumps := prepareForSolve(t, "test_too_simple.txt")
	solution, wasPruned, err := Solve(initial, jumps, 0, 0)
	require.NoError(t, err)
	require.NotNil(t, solution)
	assert.False(t, wasPruned)
	assert.Len(t, solution, 2)
	totalJumps := 0
	for _, move := range solution {
		totalJumps += len(move)
	}
	assert.Equal(t, 3, totalJumps)
}

func TestSolve_TooSimpleUnsolvable(t *testing.T) {
	initial, jumps := prepareForSolve(t, "test_too_simple_unsolvable.txt")
	solution, _, err := Solve(initial, jumps, 0, 0)
	require.NoError(t, err)
	assert.Nil(t, solution)
}
