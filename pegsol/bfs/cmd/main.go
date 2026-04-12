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
	maxStates int
}

func parseArgs() (*args, error) {
	maxStates := flag.Int("max-states", 0, "maximum number of states per BFS level (0 = unlimited); must be 0 or >= 1,000,000")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: pegsol-bfs [options] <input-file>\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  <input-file>   path to a peg solitaire board file\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		return nil, fmt.Errorf("missing required argument: input file")
	}

	if *maxStates > 0 && *maxStates < 1_000_000 {
		return nil, fmt.Errorf("max-states value %d is too small; must be 0 (unlimited) or >= 1,000,000 to detect approximately optimal solutions", *maxStates)
	}

	return &args{
		inputFile: flag.Arg(0),
		maxStates: *maxStates,
	}, nil
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

	solution, wasPruned, err := bfs.Solve(initialState, compactJumps, seedVal, a.maxStates)
	end := time.Now()
	fmt.Println("Process ended at:", end.Format("2006/01/02 15:04:05"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError during solving: %v\n", err)
		os.Exit(1)
	}

	if wasPruned {
		fmt.Println("\nNote: solution may not be optimal due to map pruning.")
	} else {
		fmt.Println("\nNo pruning was done - an optimal solution was found.")
	}

	if solution == nil {
		fmt.Println("\nNo solution found. If pruning was applied and the puzzle is known to be solvable, consider increasing the max-states.")
		os.Exit(0)
	}

	printSolution(b, initialState, solution)

	fmt.Printf("\nSolved in %s\n", end.Sub(start).Round(time.Millisecond))
}
