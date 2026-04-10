package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"peg_solitaire/pegsol/bfs"
	"peg_solitaire/pegsol/board"
	"peg_solitaire/pegsol/matrixstate"
)

type args struct {
	inputFile string
	seed      uint64
}

func parseArgs() (*args, error) {
	seed := flag.Uint64("seed", 0, "optional random seed for jump shuffling; 0 or omitted uses a random seed each run, producing different solutions")
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

	return &args{
		inputFile: flag.Arg(0),
		seed:      *seed,
	}, nil
}

func getSeedValue(seed uint64) uint64 {
	if seed == 0 {
		seedVal := uint64(time.Now().UnixNano())
		fmt.Printf("A default zero seed. new seed is generated: %d\n", seedVal)
		return seedVal
	}
	fmt.Printf("Using provided seed: %d\n", seed)
	return seed
}

func printSolution(b *board.Board, solution [][]*board.CompactJump) {
	fmt.Print("\n\n--- Solution ---\n\n")
	for i, move := range solution {
		fmt.Printf("Move %d:\n", i+1)
		for _, jump := range move {
			desc, err := b.DescribeJump(jump)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\nError describing jump in move %d: %v\n", i+1, err)
				os.Exit(1)
			}
			fmt.Println(desc)
		}
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

	start := time.Now()
	fmt.Println("Process started at:", start.Format("2006/01/02 15:04:05"))

	seedVal := getSeedValue(a.seed)

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

	printSolution(b, solution)

	fmt.Printf("\nSolved in %s\n", end.Sub(start).Round(time.Millisecond))
}
