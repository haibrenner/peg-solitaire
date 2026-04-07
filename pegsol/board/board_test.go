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

func TestAllPossibleSteps(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	b := NewBoard(ms)
	assert.Len(t, b.Steps, 76)
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

func TestCompactStepIsValidOn(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)

	b := NewBoard(ms)

	cs, err := b.TranslateMatrixToCompactState(ms)
	require.NoError(t, err)

	compactSteps, err := b.TranslateAllCoordStepsToCompact()
	require.NoError(t, err)

	var validSteps []*CompactStep
	for _, cm := range compactSteps {
		if cm.IsValidOn(cs) {
			validSteps = append(validSteps, cm)
		}
	}

	require.Len(t, validSteps, 4)

	for i := 1; i < len(validSteps); i++ {
		assert.Equal(t, validSteps[0].EndPosition, validSteps[i].EndPosition)
	}

	seenOccupied := map[bitmap.BitmapKey]bool{}
	seenStart := map[int]bool{}
	seenFull := map[bitmap.BitmapKey]bool{}
	for _, cm := range validSteps {
		assert.False(t, seenOccupied[cm.OccupiedMask.Key()], "duplicate OccupiedMask")
		assert.False(t, seenStart[cm.StartPosition], "duplicate StartPosition")
		assert.False(t, seenFull[cm.FullMask.Key()], "duplicate FullMask")
		seenOccupied[cm.OccupiedMask.Key()] = true
		seenStart[cm.StartPosition] = true
		seenFull[cm.FullMask.Key()] = true
	}
}

func TestCompactStepApply(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)

	b := NewBoard(ms)

	cs, err := b.TranslateMatrixToCompactState(ms)
	require.NoError(t, err)

	step := CoordStep{
		JumpFrom: position.Position{Row: 3, Col: 5},
		JumpOver: position.Position{Row: 3, Col: 4},
		JumpTo:   position.Position{Row: 3, Col: 3},
	}
	compactStep, err := b.TranslateCoordStepToCompact(step)
	require.NoError(t, err)
	require.True(t, compactStep.IsValidOn(cs))

	resultCs := compactStep.Apply(cs)

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
