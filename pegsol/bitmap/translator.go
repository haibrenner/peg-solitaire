package bitmap

import (
	"fmt"
	"peg_solitaire/pegsol/position"
)

type Translator struct {
	positions []position.Position
	toIndex   map[position.Position]int
}

func NewTranslator(positions []position.Position) (*Translator, error) {
	if len(positions) > 64 {
		return nil, fmt.Errorf("number of positions %d exceeds bitmap capacity of 64", len(positions))
	}
	t := &Translator{
		positions: make([]position.Position, len(positions)),
		toIndex:   make(map[position.Position]int, len(positions)),
	}
	copy(t.positions, positions)
	for i, p := range t.positions {
		t.toIndex[p] = i
	}
	return t, nil
}

func (t *Translator) ToPosition(index int) (position.Position, error) {
	if index < 0 || index >= len(t.positions) {
		return position.Position{}, fmt.Errorf("index %d out of range [0, %d)", index, len(t.positions))
	}
	return t.positions[index], nil
}

func (t *Translator) ToPositions(indices []int) ([]position.Position, error) {
	result := make([]position.Position, len(indices))
	for i, idx := range indices {
		p, err := t.ToPosition(idx)
		if err != nil {
			return nil, err
		}
		result[i] = p
	}
	return result, nil
}

func (t *Translator) ToIndex(p position.Position) (int, error) {
	i, ok := t.toIndex[p]
	if !ok {
		return 0, fmt.Errorf("position %v not found", p)
	}
	return i, nil
}

func (t *Translator) ToIndices(positions []position.Position) ([]int, error) {
	result := make([]int, len(positions))
	for i, p := range positions {
		idx, err := t.ToIndex(p)
		if err != nil {
			return nil, err
		}
		result[i] = idx
	}
	return result, nil
}

func (t *Translator) PositionsToBitmap(positions []position.Position) (Bitmap, error) {
	indices, err := t.ToIndices(positions)
	if err != nil {
		return 0, err
	}
	return FromInts(indices), nil
}

func (t *Translator) BitmapToPositions(b Bitmap) ([]position.Position, error) {
	return t.ToPositions(b.ToInts())
}
