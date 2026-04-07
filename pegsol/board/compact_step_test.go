package board

import (
	"testing"

	"peg_solitaire/pegsol/matrixstate"
	"peg_solitaire/pegsol/position"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescribeStep_Valid(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	b := NewBoard(ms)

	coordStep := CoordStep{
		JumpFrom: position.Position{Row: 3, Col: 5},
		JumpOver: position.Position{Row: 3, Col: 4},
		JumpTo:   position.Position{Row: 3, Col: 3},
	}
	cs, err := b.TranslateCoordStepToCompact(coordStep)
	require.NoError(t, err)

	desc, err := b.DescribeStep(cs)
	require.NoError(t, err)
	assert.Equal(t, "Peg (4, 6) Left", desc)
}

func TestDescribeStep_InvalidDirection(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	b := NewBoard(ms)

	coordStep := CoordStep{
		JumpFrom: position.Position{Row: 3, Col: 5},
		JumpOver: position.Position{Row: 3, Col: 4},
		JumpTo:   position.Position{Row: 3, Col: 2},
	}
	cs, err := b.TranslateCoordStepToCompact(coordStep)
	require.NoError(t, err)

	_, err = b.DescribeStep(cs)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot determine direction")
}
