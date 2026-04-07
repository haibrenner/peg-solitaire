package board

import (
	"testing"

	"peg_solitaire/pegsol/bitmap"
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

func TestCompactAtomicMoveIsValidOn(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)

	b := NewBoard(ms)

	cs, err := b.TranslateMatrixToCompactState(ms)
	require.NoError(t, err)

	compactMoves, err := b.TranslateMultipleCoordMovesToCompact()
	require.NoError(t, err)

	var validMoves []*CompactAtomicMove
	for _, cm := range compactMoves {
		if cm.IsValidOn(cs) {
			validMoves = append(validMoves, cm)
		}
	}

	require.Len(t, validMoves, 4)

	// all end positions must be equal, all other fields must be distinct
	for i := 1; i < len(validMoves); i++ {
		assert.Equal(t, validMoves[0].EndPosition, validMoves[i].EndPosition)
	}

	seenOccupied := map[bitmap.BitmapKey]bool{}
	seenStart := map[int]bool{}
	seenFull := map[bitmap.BitmapKey]bool{}
	for _, cm := range validMoves {
		assert.False(t, seenOccupied[cm.OccupiedMask.Key()], "duplicate OccupiedMask")
		assert.False(t, seenStart[cm.StartPosition], "duplicate StartPosition")
		assert.False(t, seenFull[cm.FullMask.Key()], "duplicate FullMask")
		seenOccupied[cm.OccupiedMask.Key()] = true
		seenStart[cm.StartPosition] = true
		seenFull[cm.FullMask.Key()] = true
	}
}
