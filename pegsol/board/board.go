package board

import (
	"peg_solitaire/pegsol/bitmap"
	"peg_solitaire/pegsol/matrixstate"
	"peg_solitaire/pegsol/position"
)

type Board struct {
	holes      [][]bool
	Moves      []Move
	Translator *bitmap.Translator
}

func NewBoard(ms *matrixstate.MatrixState) *Board {
	holes := getBoardHoleMatrix(ms)
	positions := getValidPositions(holes)
	translator := bitmap.NewTranslator(positions)

	b := &Board{
		holes:      holes,
		Translator: translator,
	}
	b.Moves = b.allPossibleMoves()
	return b
}

func getBoardHoleMatrix(ms *matrixstate.MatrixState) [][]bool {
	holeMat := make([][]bool, len(ms.Cells))
	for r, row := range ms.Cells {
		holeMat[r] = make([]bool, len(row))
		for c, cell := range row {
			holeMat[r][c] = (cell != matrixstate.CellFiller)
		}
	}
	return holeMat
}

func getValidPositions(validity [][]bool) []position.Position {
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
	var pegs []position.Position
	for r, row := range ms.Cells {
		for c, cell := range row {
			if cell == matrixstate.CellPeg {
				pegs = append(pegs, position.Position{Row: r, Col: c})
			}
		}
	}
	return b.Translator.PositionsToBitmap(pegs)
}
