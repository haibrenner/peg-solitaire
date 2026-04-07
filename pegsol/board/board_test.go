package board

import (
	"testing"

	"peg_solitaire/pegsol/matrixstate"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const inputsDir = "../inputs/"

func TestAllPossibleMoves(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	b := NewBoard(ms)
	assert.Len(t, b.Moves, 76)
}

func TestTranslateMatrixRoundtrip(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	ms.Cells[0][0] = matrixstate.CellPeg

	b := NewBoard(ms)

	cs, err := b.TranslateMatrixToCompactState(ms)
	require.NoError(t, err)

	got, err := b.TranslateCompactToMatrixState(cs)
	require.NoError(t, err)

	assert.Equal(t, ms, got)
}
