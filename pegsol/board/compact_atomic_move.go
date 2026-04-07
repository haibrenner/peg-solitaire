package board

import (
	"peg_solitaire/pegsol/bitmap"
)

type CompactAtomicStep struct {
	FullMask      bitmap.Bitmap
	OccupiedMask  bitmap.Bitmap
	StartPosition int
	EndPosition   int
}

func (m *CompactAtomicStep) IsValidOn(cs CompactState) bool {
	return cs.And(m.FullMask).Equal(m.OccupiedMask)
}

func (m *CompactAtomicStep) Apply(cs CompactState) CompactState {
	return CompactState{cs.Xor(m.FullMask)}
}
