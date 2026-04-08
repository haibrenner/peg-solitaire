package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"peg_solitaire/pegsol/board"
	"peg_solitaire/pegsol/dfs"
	"peg_solitaire/pegsol/matrixstate"
)

type args struct {
	inputFile string
	seed      *int
}

func parseArgs() (*args, error) {
	seed := flag.Int("seed", 0, "optional random seed for jump shuffling; if omitted, a random seed is used each run, producing different solutions")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: pegsol [options] <input-file>\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  <input-file>   path to a peg solitaire board file\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		return nil, fmt.Errorf("missing required argument: input file")
	}

	var seedPtr *int
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "seed" {
			seedPtr = seed
		}
	})

	return &args{
		inputFile: flag.Arg(0),
		seed:      seedPtr,
	}, nil
}

func printSolution(b *board.Board, initial board.CompactState, jumps []*board.CompactJump) {
	fmt.Print("\n\n--- Solution ---\n\n")
	state := initial
	for i, jump := range jumps {
		desc, err := b.DescribeJump(jump)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError describing jump %d: %v\n", i+1, err)
			os.Exit(1)
		}
		fmt.Printf("Jump %d: %s\n\n", i+1, desc)
		state = jump.Apply(state)
		ms, err := b.TranslateCompactToMatrixState(state)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError translating state at jump %d: %v\n", i+1, err)
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

	b := board.NewBoard(ms)

	compactJumps, err := b.TranslateAllCoordJumpsToCompact()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError translating jumps: %v\n", err)
		os.Exit(1)
	}

	var seedVal uint64
	if a.seed != nil {
		seedVal = uint64(*a.seed)
		fmt.Printf("\nUsing provided seed: %d\n", seedVal)
	} else {
		seedVal = uint64(time.Now().UnixNano())
		fmt.Printf("\nUsing automatic system seed: %d\n", seedVal)
	}

	pcg := rand.NewPCG(seedVal, seedVal+1) // PCG likes two different values

	r := rand.New(pcg)

	r.Shuffle(len(compactJumps), func(i, j int) {
		compactJumps[i], compactJumps[j] = compactJumps[j], compactJumps[i]
	})

	initialState, err := b.TranslateMatrixToCompactState(ms)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError creating compact state: %v\n\n", err)
		os.Exit(1)
	}

	start := time.Now()
	fmt.Println("Process started at:", start.Format("2006-01-02 15:04:05"))
	solution := dfs.Solve(initialState, compactJumps)
	end := time.Now()
	fmt.Println("Process ended at:", end.Format("2006-01-02 15:04:05"))

	if solution == nil {
		fmt.Println("\nThe puzzle has no solution.")
		os.Exit(0)
	}

	printSolution(b, initialState, solution)

	fmt.Printf("\nSolved in %s\n", end.Sub(start).Round(time.Millisecond))

}
