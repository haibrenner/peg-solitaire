package bitmap

import (
	"fmt"
	"peg_solitaire/pegsol/position"
)

type Translator struct {
	positions  []position.Position
	toIndex    map[position.Position]int
	BitmapSize int
}

func NewTranslator(positions []position.Position) *Translator {
	t := &Translator{
		positions:  make([]position.Position, len(positions)),
		toIndex:    make(map[position.Position]int, len(positions)),
		BitmapSize: (len(positions) + 63) / 64,
	}
	copy(t.positions, positions)
	for i, p := range t.positions {
		t.toIndex[p] = i
	}
	return t
}

func (t *Translator) ToPosition(index int) (position.Position, error) {
	if index < 0 || index >= len(t.positions) {
		return position.Position{}, fmt.Errorf("index %d out of range [0, %d)", index, len(t.positions))
	}
	return t.positions[index], nil
}

func (t *Translator) ToIndex(p position.Position) (int, error) {
	i, ok := t.toIndex[p]
	if !ok {
		return 0, fmt.Errorf("position %v not found", p)
	}
	return i, nil
}
