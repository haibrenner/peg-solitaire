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
	seed := flag.Int("seed", 0, "optional random seed for step shuffling; if omitted, a random seed is used each run, producing different solutions")
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

func printSolution(b *board.Board, initial board.CompactState, steps []*board.CompactStep) {
	fmt.Println("\n--- Solution ---")
	state := initial
	for i, step := range steps {
		desc, err := b.DescribeStep(step)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error describing step %d: %v\n", i+1, err)
			os.Exit(1)
		}
		fmt.Printf("Step %d: %s\n", i+1, desc)
		state = step.Apply(state)
		ms, err := b.TranslateCompactToMatrixState(state)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error translating state at step %d: %v\n", i+1, err)
			os.Exit(1)
		}
		fmt.Print(ms.String())
	}
}

func main() {
	a, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	ms, err := matrixstate.ReadInput(a.inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(ms.String())

	if ms.IsAlgebraicallyInfeasible() {
		fmt.Println("The puzzle is algebraically infeasible and cannot be solved.")
		os.Exit(0)
	}

	b := board.NewBoard(ms)

	compactSteps, err := b.TranslateAllCoordStepsToCompact()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error translating steps: %v\n", err)
		os.Exit(1)
	}

	var seedVal uint64
	if a.seed != nil {
		seedVal = uint64(*a.seed)
		fmt.Printf("Using provided seed: %d\n", seedVal)
	} else {
		seedVal = uint64(time.Now().UnixNano())
		fmt.Printf("Using automatic system seed: %d\n", seedVal)
	}

	pcg := rand.NewPCG(seedVal, seedVal+1) // PCG likes two different values

	r := rand.New(pcg)

	r.Shuffle(len(compactSteps), func(i, j int) {
		compactSteps[i], compactSteps[j] = compactSteps[j], compactSteps[i]
	})

	initialState, err := b.TranslateMatrixToCompactState(ms)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating compact state: %v\n", err)
		os.Exit(1)
	}

	start := time.Now()
	fmt.Println("Solving started at:", start.Format("2006-01-02 15:04:05"))
	solution := dfs.Solve(initialState, compactSteps)
	end := time.Now()
	fmt.Println("Solving ended at:", end.Format("2006-01-02 15:04:05"))

	if solution == nil {
		fmt.Println("The puzzle has no solution.")
		os.Exit(0)
	}

	printSolution(b, initialState, solution)

	fmt.Printf("\nSolved in %s\n", end.Sub(start).Round(time.Millisecond))

}
