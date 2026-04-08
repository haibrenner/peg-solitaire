package bitmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromIntsToInts(t *testing.T) {
	vals := []int{1, 0, 5, 10, 63}
	b := FromInts(vals)
	got := b.ToInts()
	assert.Len(t, got, len(vals))
	assert.ElementsMatch(t, vals, got)
}
