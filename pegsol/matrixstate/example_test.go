package matrixstate_test

import (
	"fmt"
	"peg_solitaire/pegsol/matrixstate"
)

func ExampleReadInput() {
	ms, err := matrixstate.ReadInput("../inputs/standard_english.txt")
	if err != nil {
		panic(err)
	}
	fmt.Print(ms.String())
	// Output:
	// ##+++##
	// ##+++##
	// +++++++
	// +++.+++
	// +++++++
	// ##+++##
	// ##+++##
}
