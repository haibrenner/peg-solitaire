package matrixstate

type Position struct {
	Row, Col int
}

type Move struct {
	JumpingFrom, JumpedOver, JumpTo Position
}

func (ms *MatrixState) AllPossibleMoves() []Move {
	var moves []Move
	rows := len(ms.Cells)
	for r := range ms.Cells {
		cols := len(ms.Cells[r])
		for c := range ms.Cells[r] {
			if ms.Cells[r][c] == CellFiller {
				continue
			}
			// horizontal: right and left
			if c+2 < cols && ms.Cells[r][c+1] != CellFiller && ms.Cells[r][c+2] != CellFiller {
				moves = append(moves,
					Move{JumpingFrom: Position{r, c}, JumpedOver: Position{r, c + 1}, JumpTo: Position{r, c + 2}},
					Move{JumpingFrom: Position{r, c + 2}, JumpedOver: Position{r, c + 1}, JumpTo: Position{r, c}},
				)
			}
			// vertical: down and up
			if r+2 < rows && ms.Cells[r+1][c] != CellFiller && ms.Cells[r+2][c] != CellFiller {
				moves = append(moves,
					Move{JumpingFrom: Position{r, c}, JumpedOver: Position{r + 1, c}, JumpTo: Position{r + 2, c}},
					Move{JumpingFrom: Position{r + 2, c}, JumpedOver: Position{r + 1, c}, JumpTo: Position{r, c}},
				)
			}
		}
	}
	return moves
}
