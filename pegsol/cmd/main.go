package main

import (
	"flag"
	"fmt"
	"os"

	"peg_solitaire/pegsol/matrixstate"
)

type args struct {
	inputFile string
	seed      int
}

func parseArgs() (*args, error) {
	seed := flag.Int("seed", 0, "optional seed value")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: pegsol -input <file> [-seed <int>]\n\n")
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
	_ = a.seed
}
