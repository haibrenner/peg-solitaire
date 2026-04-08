package board

import (
	"fmt"
	"peg_solitaire/pegsol/bitmap"
)

type CompactJump struct {
	FullMask      bitmap.Bitmap
	OccupiedMask  bitmap.Bitmap
	StartPosition int
	EndPosition   int
}

func (cj *CompactJump) IsValidOn(cs CompactState) bool {
	return cs.Bitmap&cj.FullMask == cj.OccupiedMask
}

func (cj *CompactJump) Apply(cs CompactState) CompactState {
	return CompactState{cs.Bitmap ^ cj.FullMask}
}

func (b *Board) DescribeJump(cj *CompactJump) (string, error) {
	start, err := b.Translator.ToPosition(cj.StartPosition)
	if err != nil {
		return "", err
	}
	end, err := b.Translator.ToPosition(cj.EndPosition)
	if err != nil {
		return "", err
	}
	var direction string
	switch {
	case end.Row == start.Row && end.Col == start.Col+2:
		direction = "Right"
	case end.Row == start.Row && end.Col == start.Col-2:
		direction = "Left"
	case end.Col == start.Col && end.Row == start.Row+2:
		direction = "Down"
	case end.Col == start.Col && end.Row == start.Row-2:
		direction = "Up"
	default:
		return "", fmt.Errorf("cannot determine direction from %v to %v", start, end)
	}
	return fmt.Sprintf("Peg (%d, %d) %s", start.Row+1, start.Col+1, direction), nil
}
