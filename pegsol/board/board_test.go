package board

import (
	"testing"

	"peg_solitaire/pegsol/bitmap"
	"peg_solitaire/pegsol/matrixstate"
	"peg_solitaire/pegsol/position"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const inputsDir = "../inputs/"

func TestAllPossibleJumps(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	b, err := NewBoard(ms)
	require.NoError(t, err)
	assert.Len(t, b.Jumps, 76)
}

func TestTranslateMatrixRoundtrip(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	ms.Cells[0][0] = matrixstate.CellPeg

	b, err := NewBoard(ms)
	require.NoError(t, err)

	cs, err := b.TranslateMatrixToCompactState(ms)
	require.NoError(t, err)

	got, err := b.TranslateCompactToMatrixState(cs)
	require.NoError(t, err)

	assert.Equal(t, ms, got)
}

func TestCompactJumpIsValidOn(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)

	b, err := NewBoard(ms)
	require.NoError(t, err)

	cs, err := b.TranslateMatrixToCompactState(ms)
	require.NoError(t, err)

	compactJumps, err := b.TranslateAllCoordJumpsToCompact()
	require.NoError(t, err)

	var validJumps []*CompactJump
	for _, cm := range compactJumps {
		if cm.IsValidOn(cs) {
			validJumps = append(validJumps, cm)
		}
	}

	require.Len(t, validJumps, 4)

	for i := 1; i < len(validJumps); i++ {
		assert.Equal(t, validJumps[0].EndPosition, validJumps[i].EndPosition)
	}

	seenOccupied := map[bitmap.Bitmap]bool{}
	seenStart := map[int8]bool{}
	seenFull := map[bitmap.Bitmap]bool{}
	for _, cm := range validJumps {
		assert.False(t, seenOccupied[cm.OccupiedMask], "duplicate OccupiedMask")
		assert.False(t, seenStart[cm.StartPosition], "duplicate StartPosition")
		assert.False(t, seenFull[cm.FullMask], "duplicate FullMask")
		seenOccupied[cm.OccupiedMask] = true
		seenStart[cm.StartPosition] = true
		seenFull[cm.FullMask] = true
	}
}

func TestCompactJumpApply(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)

	b, err := NewBoard(ms)
	require.NoError(t, err)

	cs, err := b.TranslateMatrixToCompactState(ms)
	require.NoError(t, err)

	jump := CoordJump{
		JumpFrom: position.Position{Row: 3, Col: 5},
		JumpOver: position.Position{Row: 3, Col: 4},
		JumpTo:   position.Position{Row: 3, Col: 3},
	}
	compactJump, err := b.TranslateCoordJumpToCompact(jump)
	require.NoError(t, err)
	require.True(t, compactJump.IsValidOn(cs))

	resultCs := compactJump.Apply(cs)

	resultMs, err := b.TranslateCompactToMatrixState(resultCs)
	require.NoError(t, err)

	var diffs []position.Position
	for r, row := range ms.Cells {
		for c, cell := range row {
			if cell != resultMs.Cells[r][c] {
				diffs = append(diffs, position.Position{Row: r, Col: c})
			}
		}
	}
	assert.ElementsMatch(t, []position.Position{
		{Row: 3, Col: 5},
		{Row: 3, Col: 4},
		{Row: 3, Col: 3},
	}, diffs)

	assert.Equal(t, matrixstate.CellPeg, resultMs.Cells[3][3])
	assert.Equal(t, matrixstate.CellHole, resultMs.Cells[3][4])
	assert.Equal(t, matrixstate.CellHole, resultMs.Cells[3][5])
}
