package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"peg_solitaire/pegsol/bfs"
	"peg_solitaire/pegsol/board"
	"peg_solitaire/pegsol/matrixstate"
)

type args struct {
	inputFile string
}

func parseArgs() (*args, error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: pegsol-bfs <input-file>\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  <input-file>   path to a peg solitaire board file\n\n")
	}
	flag.Parse()

	if flag.NArg() < 1 {
		return nil, fmt.Errorf("missing required argument: input file")
	}

	return &args{inputFile: flag.Arg(0)}, nil
}

func printSolution(b *board.Board, initial board.CompactState, solution [][]*board.CompactJump) {
	fmt.Print("\n\n--- Solution ---\n\n")
	state := initial
	for i, move := range solution {
		if len(move) == 0 {
			continue
		}
		startPos, err := b.Translator.ToPosition(int(move[0].StartPosition))
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError getting start position for move %d: %v\n", i+1, err)
			os.Exit(1)
		}
		directions := make([]string, len(move))
		for j, jump := range move {
			directions[j] = jump.Direction
			state = jump.Apply(state)
		}
		fmt.Printf("Move %d: Peg in position (%d, %d): %s\n", i+1, startPos.Row+1, startPos.Col+1, strings.Join(directions, ", "))
		ms, err := b.TranslateCompactToMatrixState(state)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError translating state after move %d: %v\n", i+1, err)
			os.Exit(1)
		}
		fmt.Println(ms.String())
	}
}

func main() {
	a, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	ms, err := matrixstate.ReadInput(a.inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nInitial board state:")
	fmt.Println(ms.String())

	if ms.IsAlgebraicallyInfeasible() {
		fmt.Println("\nThe puzzle is algebraically infeasible and cannot be solved.")
		os.Exit(0)
	}

	b, err := board.NewBoard(ms)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError creating board: %v\n", err)
		os.Exit(1)
	}

	compactJumps, err := b.TranslateAllCoordJumpsToCompact()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError translating jumps: %v\n", err)
		os.Exit(1)
	}

	initialState, err := b.TranslateMatrixToCompactState(ms)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError creating compact state: %v\n\n", err)
		os.Exit(1)
	}

	seedVal := uint64(time.Now().UnixNano())
	fmt.Printf("Using seed: %d\n", seedVal)

	start := time.Now()
	fmt.Println("Process started at:", start.Format("2006/01/02 15:04:05"))

	solution, err := bfs.Solve(initialState, compactJumps, seedVal)
	end := time.Now()
	fmt.Println("Process ended at:", end.Format("2006/01/02 15:04:05"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError during solving: %v\n", err)
		os.Exit(1)
	}

	if solution == nil {
		fmt.Println("\nThe puzzle has no solution.")
		os.Exit(0)
	}

	printSolution(b, initialState, solution)

	fmt.Printf("\nSolved in %s\n", end.Sub(start).Round(time.Millisecond))
}
