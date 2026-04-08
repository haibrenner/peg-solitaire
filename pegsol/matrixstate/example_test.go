package matrixstate_test

import (
	"peg_solitaire/pegsol/matrixstate"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadInputPrintsCorrectly(t *testing.T) {
	ms, err := matrixstate.ReadInput("../inputs/standard_english.txt")
	require.NoError(t, err)
	expected := strings.Join([]string{
		"  +++  ",
		"  +++  ",
		"+++++++",
		"+++.+++",
		"+++++++",
		"  +++  ",
		"  +++  ",
	}, "\n") + "\n"
	assert.Equal(t, expected, ms.String())
}
