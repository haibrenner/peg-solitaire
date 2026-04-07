package board

import (
	"fmt"
	"peg_solitaire/pegsol/bitmap"
	"peg_solitaire/pegsol/position"
)

type CompactAtomicMove struct {
	FullMask      bitmap.Bitmap
	OccupiedMask  bitmap.Bitmap
	StartPosition int
	EndPosition   int
}

func (b *Board) TranslateCoordAtomicMoveToCompact(m CoordAtomicMove) (*CompactAtomicMove, error) {
	fullMask, err := b.Translator.PositionsToBitmap([]position.Position{m.JumpFrom, m.JumpOver, m.JumpTo})
	if err != nil {
		return nil, fmt.Errorf("failed to build FullMask: %w", err)
	}
	occupiedMask, err := b.Translator.PositionsToBitmap([]position.Position{m.JumpFrom, m.JumpOver})
	if err != nil {
		return nil, fmt.Errorf("failed to build OccupiedMask: %w", err)
	}
	startPos, err := b.Translator.ToIndex(m.JumpFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to get StartPosition: %w", err)
	}
	endPos, err := b.Translator.ToIndex(m.JumpTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get EndPosition: %w", err)
	}
	return &CompactAtomicMove{
		FullMask:      fullMask,
		OccupiedMask:  occupiedMask,
		StartPosition: startPos,
		EndPosition:   endPos,
	}, nil
}

func (b *Board) TranslateMultipleMovesToCompact() ([]*CompactAtomicMove, error) {
	result := make([]*CompactAtomicMove, len(b.Moves))
	for i, m := range b.Moves {
		cm, err := b.TranslateCoordAtomicMoveToCompact(m)
		if err != nil {
			return nil, err
		}
		result[i] = cm
	}
	return result, nil
}
