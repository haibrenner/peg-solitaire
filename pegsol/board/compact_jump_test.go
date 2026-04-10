package board

import (
	"testing"

	"peg_solitaire/pegsol/matrixstate"
	"peg_solitaire/pegsol/position"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescribeJump_Valid(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	b, err := NewBoard(ms)
	require.NoError(t, err)

	coordJump := CoordJump{
		JumpFrom:  position.Position{Row: 3, Col: 5},
		JumpOver:  position.Position{Row: 3, Col: 4},
		JumpTo:    position.Position{Row: 3, Col: 3},
		Direction: "Left",
	}
	cs, err := b.TranslateCoordJumpToCompact(coordJump)
	require.NoError(t, err)

	desc, err := b.DescribeJump(cs)
	require.NoError(t, err)
	assert.Equal(t, "Peg (4, 6) Left", desc)
}

func TestDescribeJump_InvalidDirection(t *testing.T) {
	ms, err := matrixstate.ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	b, err := NewBoard(ms)
	require.NoError(t, err)

	coordJump := CoordJump{
		JumpFrom: position.Position{Row: 3, Col: 5},
		JumpOver: position.Position{Row: 3, Col: 4},
		JumpTo:   position.Position{Row: 3, Col: 2},
	}
	cs, err := b.TranslateCoordJumpToCompact(coordJump)
	require.NoError(t, err)

	_, err = b.DescribeJump(cs)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot determine direction")
}
