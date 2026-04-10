package board

import (
	"fmt"
	"peg_solitaire/pegsol/bitmap"
)

type CompactJump struct {
	FullMask      bitmap.Bitmap
	OccupiedMask  bitmap.Bitmap
	StartPosition int8
	EndPosition   int8
	Direction     string
}

func (cj *CompactJump) IsValidOn(cs CompactState) bool {
	return cs.Bitmap&cj.FullMask == cj.OccupiedMask
}

func (cj *CompactJump) Apply(cs CompactState) CompactState {
	return CompactState{cs.Bitmap ^ cj.FullMask}
}

func (b *Board) DescribeJump(cj *CompactJump) (string, error) {
	start, err := b.Translator.ToPosition(int(cj.StartPosition))
	if err != nil {
		return "", err
	}
	if cj.Direction == "" {
		return "", fmt.Errorf("cannot determine direction from %v", start)
	}
	return fmt.Sprintf("Peg (%d, %d) %s", start.Row+1, start.Col+1, cj.Direction), nil
}
