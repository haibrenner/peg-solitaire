package board

import "peg_solitaire/pegsol/position"

type Move struct {
	JumpingFrom, JumpedOver, JumpTo position.Position
}

func (b *Board) allPossibleMoves() []Move {
	var moves []Move
	rows := len(b.holes)
	for r := range b.holes {
		cols := len(b.holes[r])
		for c := range b.holes[r] {
			if !b.holes[r][c] {
				continue
			}
			// horizontal: right and left
			if c+2 < cols && b.holes[r][c+1] && b.holes[r][c+2] {
				moves = append(moves,
					Move{JumpingFrom: position.Position{Row: r, Col: c}, JumpedOver: position.Position{Row: r, Col: c + 1}, JumpTo: position.Position{Row: r, Col: c + 2}},
					Move{JumpingFrom: position.Position{Row: r, Col: c + 2}, JumpedOver: position.Position{Row: r, Col: c + 1}, JumpTo: position.Position{Row: r, Col: c}},
				)
			}
			// vertical: down and up
			if r+2 < rows && b.holes[r+1][c] && b.holes[r+2][c] {
				moves = append(moves,
					Move{JumpingFrom: position.Position{Row: r, Col: c}, JumpedOver: position.Position{Row: r + 1, Col: c}, JumpTo: position.Position{Row: r + 2, Col: c}},
					Move{JumpingFrom: position.Position{Row: r + 2, Col: c}, JumpedOver: position.Position{Row: r + 1, Col: c}, JumpTo: position.Position{Row: r, Col: c}},
				)
			}
		}
	}
	return moves
}
