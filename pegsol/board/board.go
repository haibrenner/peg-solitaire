package board

import (
	"peg_solitaire/pegsol/bitmap"
	"peg_solitaire/pegsol/matrixstate"
	"peg_solitaire/pegsol/position"
)

type Board struct {
	validity   [][]bool
	Moves      []Move
	Translator *bitmap.Translator
	bitmapSize int
}

func NewBoard(ms *matrixstate.MatrixState) *Board {
	validity := buildValidity(ms)
	positions := validPositions(validity)
	translator := bitmap.NewTranslator(positions)
	bitmapSize := (len(positions) + 63) / 64

	b := &Board{
		validity:   validity,
		Translator: translator,
		bitmapSize: bitmapSize,
	}
	b.Moves = b.allPossibleMoves()
	return b
}

func buildValidity(ms *matrixstate.MatrixState) [][]bool {
	validity := make([][]bool, len(ms.Cells))
	for r, row := range ms.Cells {
		validity[r] = make([]bool, len(row))
		for c, cell := range row {
			validity[r][c] = cell != matrixstate.CellFiller
		}
	}
	return validity
}

func validPositions(validity [][]bool) []position.Position {
	var positions []position.Position
	for r, row := range validity {
		for c, valid := range row {
			if valid {
				positions = append(positions, position.Position{Row: r, Col: c})
			}
		}
	}
	return positions
}

func (b *Board) ToBitmap(ms *matrixstate.MatrixState) (bitmap.Bitmap, error) {
	var pegIndices []int
	for r, row := range ms.Cells {
		for c, cell := range row {
			if cell == matrixstate.CellPeg {
				i, err := b.Translator.ToIndex(position.Position{Row: r, Col: c})
				if err != nil {
					return nil, err
				}
				pegIndices = append(pegIndices, i)
			}
		}
	}
	return bitmap.FromInts(pegIndices, b.bitmapSize), nil
}
