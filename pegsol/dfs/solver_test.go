package dfs

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

	b := board.NewBoard(ms)

	jumps, err := b.TranslateAllCoordJumpsToCompact()
	require.NoError(t, err)

	initial, err := b.TranslateMatrixToCompactState(ms)
	require.NoError(t, err)

	return initial, jumps
}

func TestSolve_TooSimple(t *testing.T) {
	initial, jumps := prepareForSolve(t, "too_simple.txt")
	solution := Solve(initial, jumps, 0)
	require.NotNil(t, solution)
	assert.Len(t, solution, 3)
}

func TestSolve_TooSimpleUnsolvable(t *testing.T) {
	initial, jumps := prepareForSolve(t, "too_simple_unsolvable.txt")
	solution := Solve(initial, jumps, 0)
	assert.Nil(t, solution)
}
