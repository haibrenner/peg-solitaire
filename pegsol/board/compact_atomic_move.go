package board

import (
	"peg_solitaire/pegsol/bitmap"
)

type CompactAtomicMove struct {
	FullMask      bitmap.Bitmap
	OccupiedMask  bitmap.Bitmap
	StartPosition int
	EndPosition   int
}

func (m *CompactAtomicMove) IsValidOn(cs CompactState) bool {
	return cs.And(m.FullMask).Equal(m.OccupiedMask)
}

func (m *CompactAtomicMove) Apply(cs CompactState) CompactState {
	return CompactState{cs.Xor(m.FullMask)}
}
