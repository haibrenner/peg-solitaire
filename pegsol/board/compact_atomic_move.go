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
