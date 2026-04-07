package matrixstate

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	CellFiller = '#'
	CellPeg    = '+'
	CellHole   = '.'
)

type MatrixState struct {
	Cells [][]rune
}

func (ms *MatrixState) String() string {
	var sb strings.Builder
	for _, row := range ms.Cells {
		sb.WriteString(string(row))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func ReadInput(path string) (*MatrixState, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	var rows [][]rune
	dataStarted := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimRight(line, " \t")

		if trimmed == "" {
			if dataStarted {
				rows = append(rows, nil) // sentinel for blank line after data
			}
			continue
		}

		if dataStarted && len(rows) > 0 && rows[len(rows)-1] == nil {
			return nil, fmt.Errorf("malformed input: blank line in the middle of data")
		}

		for _, ch := range trimmed {
			if ch != CellFiller && ch != CellPeg && ch != CellHole {
				return nil, fmt.Errorf("malformed input: data must contain %q, %q, or %q, with allowed trailing whitespaces only; got %q", CellFiller, CellPeg, CellHole, ch)
			}
		}

		if dataStarted && len(trimmed) != len([]rune(string(rows[0]))) {
			return nil, fmt.Errorf("malformed input: inconsistent line length, expected %d got %d", len(rows[0]), len(trimmed))
		}

		dataStarted = true
		rows = append(rows, []rune(trimmed))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// strip trailing blank sentinels
	for len(rows) > 0 && rows[len(rows)-1] == nil {
		rows = rows[:len(rows)-1]
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("malformed input: no data found")
	}

	return &MatrixState{Cells: rows}, nil
}
