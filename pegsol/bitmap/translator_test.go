package bitmap

import (
	"testing"

	"peg_solitaire/pegsol/position"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTranslator_TooManyPositions(t *testing.T) {
	positions := make([]position.Position, 65)
	for i := range positions {
		positions[i] = position.Position{Row: i, Col: 0}
	}
	_, err := NewTranslator(positions)
	require.Error(t, err)
}

var testPositions = []position.Position{
	{Row: 0, Col: 1},
	{Row: 0, Col: 2},
	{Row: 1, Col: 0},
	{Row: 1, Col: 1},
}

func newTestTranslator() *Translator {
	t, err := NewTranslator(testPositions)
	if err != nil {
		panic(err)
	}
	return t
}

func TestToPosition(t *testing.T) {
	tr := newTestTranslator()
	p, err := tr.ToPosition(2)
	require.NoError(t, err)
	assert.Equal(t, position.Position{Row: 1, Col: 0}, p)
}

func TestToPositionOutOfRange(t *testing.T) {
	tr := newTestTranslator()
	_, err := tr.ToPosition(10)
	require.Error(t, err)
}

func TestToIndex(t *testing.T) {
	tr := newTestTranslator()
	i, err := tr.ToIndex(position.Position{Row: 1, Col: 1})
	require.NoError(t, err)
	assert.Equal(t, 3, i)
}

func TestToIndexNotFound(t *testing.T) {
	tr := newTestTranslator()
	_, err := tr.ToIndex(position.Position{Row: 9, Col: 9})
	require.Error(t, err)
}

func TestToPositions(t *testing.T) {
	tr := newTestTranslator()
	got, err := tr.ToPositions([]int{0, 2})
	require.NoError(t, err)
	assert.ElementsMatch(t, []position.Position{{Row: 0, Col: 1}, {Row: 1, Col: 0}}, got)
}

func TestToIndices(t *testing.T) {
	tr := newTestTranslator()
	got, err := tr.ToIndices([]position.Position{{Row: 0, Col: 2}, {Row: 1, Col: 1}})
	require.NoError(t, err)
	assert.ElementsMatch(t, []int{1, 3}, got)
}

func TestPositionsToBitmap(t *testing.T) {
	tr := newTestTranslator()
	bm, err := tr.PositionsToBitmap([]position.Position{{Row: 0, Col: 1}, {Row: 1, Col: 1}})
	require.NoError(t, err)
	assert.ElementsMatch(t, []int{0, 3}, bm.ToInts())
}

func TestBitmapToPositions(t *testing.T) {
	tr := newTestTranslator()
	bm := FromInts([]int{0, 3})
	got, err := tr.BitmapToPositions(bm)
	require.NoError(t, err)
	assert.ElementsMatch(t, []position.Position{{Row: 0, Col: 1}, {Row: 1, Col: 1}}, got)
}

func TestPositionsBitmapRoundtrip(t *testing.T) {
	tr := newTestTranslator()
	original := []position.Position{{Row: 1, Col: 0}, {Row: 0, Col: 2}}
	bm, err := tr.PositionsToBitmap(original)
	require.NoError(t, err)
	got, err := tr.BitmapToPositions(bm)
	require.NoError(t, err)
	assert.ElementsMatch(t, original, got)
}
