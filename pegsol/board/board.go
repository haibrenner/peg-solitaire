package board

import (
	"fmt"
	"peg_solitaire/pegsol/bitmap"
	"peg_solitaire/pegsol/matrixstate"
	"peg_solitaire/pegsol/position"
)

type Move struct {
	JumpingFrom, JumpedOver, JumpTo position.Position
}

type Board struct {
	validCells [][]bool
	Moves      []Move
	Translator *bitmap.Translator
}

func NewBoard(ms *matrixstate.MatrixState) *Board {
	validCells := buildValidCells(ms)
	positions := getValidPositions(validCells)
	translator := bitmap.NewTranslator(positions)

	return &Board{
		validCells: validCells,
		Translator: translator,
		Moves:      allPossibleMoves(validCells),
	}
}

func buildValidCells(ms *matrixstate.MatrixState) [][]bool {
	validCellsMat := make([][]bool, len(ms.Cells))
	for r, row := range ms.Cells {
		validCellsMat[r] = make([]bool, len(row))
		for c, cell := range row {
			validCellsMat[r][c] = (cell != matrixstate.CellFiller)
		}
	}
	return validCellsMat
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

func allPossibleMoves(validCells [][]bool) []Move {
	var moves []Move
	rows := len(validCells)
	for r := range validCells {
		cols := len(validCells[r])
		for c := range validCells[r] {
			if !validCells[r][c] {
				continue
			}
			// horizontal: right and left
			if c+2 < cols && validCells[r][c+1] && validCells[r][c+2] {
				moves = append(moves,
					Move{JumpingFrom: position.Position{Row: r, Col: c}, JumpedOver: position.Position{Row: r, Col: c + 1}, JumpTo: position.Position{Row: r, Col: c + 2}},
					Move{JumpingFrom: position.Position{Row: r, Col: c + 2}, JumpedOver: position.Position{Row: r, Col: c + 1}, JumpTo: position.Position{Row: r, Col: c}},
				)
			}
			// vertical: down and up
			if r+2 < rows && validCells[r+1][c] && validCells[r+2][c] {
				moves = append(moves,
					Move{JumpingFrom: position.Position{Row: r, Col: c}, JumpedOver: position.Position{Row: r + 1, Col: c}, JumpTo: position.Position{Row: r + 2, Col: c}},
					Move{JumpingFrom: position.Position{Row: r + 2, Col: c}, JumpedOver: position.Position{Row: r + 1, Col: c}, JumpTo: position.Position{Row: r, Col: c}},
				)
			}
		}
	}
	return moves
}

func (b *Board) TranslateMatrixToCompactState(ms *matrixstate.MatrixState) (*CompactState, error) {
	pegs := ms.OccupiedCells()
	bm, err := b.Translator.PositionsToBitmap(pegs)
	if err != nil {
		return nil, err
	}
	return &CompactState{bm}, nil
}

func (b *Board) TranslateCompactToMatrixState(cs *CompactState) (*matrixstate.MatrixState, error) {
	cells := make([][]rune, len(b.validCells))
	for r, row := range b.validCells {
		cells[r] = make([]rune, len(row))
		for c, isValidCell := range row {
			if isValidCell {
				cells[r][c] = matrixstate.CellHole
			} else {
				cells[r][c] = matrixstate.CellFiller
			}
		}
	}
	pegs, err := b.Translator.BitmapToPositions(cs.Bitmap)
	if err != nil {
		return nil, err
	}
	for _, p := range pegs {
		if cells[p.Row][p.Col] == matrixstate.CellFiller {
			return nil, fmt.Errorf("invalid position in bitmap: %v", p)
		}
		cells[p.Row][p.Col] = matrixstate.CellPeg
	}
	return &matrixstate.MatrixState{Cells: cells}, nil
}
