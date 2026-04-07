package matrixstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const inputsDir = "../inputs/"

func TestValidInput(t *testing.T) {
	ms, err := ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	require.NotNil(t, ms)
	assert.Equal(t, 7, len(ms.Cells))
	assert.Equal(t, 7, len(ms.Cells[0]))
}

func TestBadlyFormatted1_UnequalLineLength(t *testing.T) {
	_, err := ReadInput(inputsDir + "badly_formatted1.txt")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "inconsistent line length")
}

func TestBadlyFormatted2_BlankLineInMiddle(t *testing.T) {
	_, err := ReadInput(inputsDir + "badly_formatted2.txt")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "blank line in the middle")
}

func TestBadlyFormatted3_DataLineStartsWithWhitespace(t *testing.T) {
	_, err := ReadInput(inputsDir + "badly_formatted3.txt")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "data must contain")
}

func TestBadlyFormatted4_WhitespaceInsideDataLine(t *testing.T) {
	_, err := ReadInput(inputsDir + "badly_formatted4.txt")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "data must contain")
}

func TestBadlyFormatted5_WhitespaceInsideDataLine(t *testing.T) {
	_, err := ReadInput(inputsDir + "badly_formatted5.txt")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "data must contain")
}

func TestAllPossibleMoves(t *testing.T) {
	ms, err := ReadInput(inputsDir + "standard_english.txt")
	require.NoError(t, err)
	assert.Len(t, ms.AllPossibleMoves(), 76)
}
