package board

import (
	"peg_solitaire/pegsol/bitmap"
)

type CompactStep struct {
	FullMask      bitmap.Bitmap
	OccupiedMask  bitmap.Bitmap
	StartPosition int
	EndPosition   int
}

func (m *CompactStep) IsValidOn(cs CompactState) bool {
	return cs.And(m.FullMask).Equal(m.OccupiedMask)
}

func (m *CompactStep) Apply(cs CompactState) CompactState {
	return CompactState{cs.Xor(m.FullMask)}
}
