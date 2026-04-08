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
	seed      uint64
}

func parseArgs() (*args, error) {
	seed := flag.Uint64("seed", 0, "optional random seed for jump shuffling; 0 or omitted uses a random seed each run, producing different solutions")
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

	return &args{
		inputFile: flag.Arg(0),
		seed:      *seed,
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

func getSeedValue(seed uint64) uint64 {
	var seedVal uint64
	if seed == 0 {
		seedVal = uint64(time.Now().UnixNano())
		fmt.Printf("A default zero seed. new seed is generated: %d\n", seedVal)
	} else {
		seedVal = seed
		fmt.Printf("Using provided seed: %d\n", seedVal)
	}
	return seedVal
}

func CreateRandFromSeed(seed uint64) *rand.Rand {
	pcg := rand.NewPCG(seed, seed+1) // PCG likes two different values
	return rand.New(pcg)
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

	initialState, err := b.TranslateMatrixToCompactState(ms)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError creating compact state: %v\n\n", err)
		os.Exit(1)
	}

	start := time.Now()
	fmt.Println("Process started at:", start.Format("2006-01-02 15:04:05"))

	seedVal := getSeedValue(a.seed)

	r := CreateRandFromSeed(seedVal)
	r.Shuffle(len(compactJumps), func(i, j int) {
		compactJumps[i], compactJumps[j] = compactJumps[j], compactJumps[i]
	})
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
