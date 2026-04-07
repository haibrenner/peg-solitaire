package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"peg_solitaire/pegsol/matrixstate"
)

const inputsDir = "../inputs/"

func TestAllPossibleMoves(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	b := NewBoard(ms)
	assert.Len(t, b.Moves, 76)
}
