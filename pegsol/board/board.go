package board

import (
	"peg_solitaire/pegsol/bitmap"
	"peg_solitaire/pegsol/matrixstate"
	"peg_solitaire/pegsol/position"
)

type Move struct {
	JumpingFrom, JumpedOver, JumpTo position.Position
}

type Board struct {
	holes      [][]bool
	Moves      []Move
	Translator *bitmap.Translator
}

func NewBoard(ms *matrixstate.MatrixState) *Board {
	holes := getBoardHoleMatrix(ms)
	positions := getValidPositions(holes)
	translator := bitmap.NewTranslator(positions)

	return &Board{
		holes:      holes,
		Translator: translator,
		Moves:      allPossibleMoves(holes),
	}
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

func allPossibleMoves(holes [][]bool) []Move {
	var moves []Move
	rows := len(holes)
	for r := range holes {
		cols := len(holes[r])
		for c := range holes[r] {
			if !holes[r][c] {
				continue
			}
			// horizontal: right and left
			if c+2 < cols && holes[r][c+1] && holes[r][c+2] {
				moves = append(moves,
					Move{JumpingFrom: position.Position{Row: r, Col: c}, JumpedOver: position.Position{Row: r, Col: c + 1}, JumpTo: position.Position{Row: r, Col: c + 2}},
					Move{JumpingFrom: position.Position{Row: r, Col: c + 2}, JumpedOver: position.Position{Row: r, Col: c + 1}, JumpTo: position.Position{Row: r, Col: c}},
				)
			}
			// vertical: down and up
			if r+2 < rows && holes[r+1][c] && holes[r+2][c] {
				moves = append(moves,
					Move{JumpingFrom: position.Position{Row: r, Col: c}, JumpedOver: position.Position{Row: r + 1, Col: c}, JumpTo: position.Position{Row: r + 2, Col: c}},
					Move{JumpingFrom: position.Position{Row: r + 2, Col: c}, JumpedOver: position.Position{Row: r + 1, Col: c}, JumpTo: position.Position{Row: r, Col: c}},
				)
			}
		}
	}
	return moves
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
