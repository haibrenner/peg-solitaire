package bitmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromIntsToInts(t *testing.T) {
	vals := []int{0, 1, 63, 64, 127, 128}
	b := FromInts(vals, 3)
	got := b.ToInts()
	assert.Len(t, got, len(vals))
	assert.Equal(t, vals, got)
}

func TestKeyRoundtrip(t *testing.T) {
	b := FromInts([]int{3, 7, 65, 190}, 3)
	got := KeyToBitmap(b.Key(), 3)
	assert.True(t, b.Equal(got))
}

func TestKeyUsableAsMapKey(t *testing.T) {
	m := map[BitmapKey]int{}
	b1 := FromInts([]int{1, 2, 3}, 2)
	b2 := FromInts([]int{4, 5, 6}, 2)
	m[b1.Key()] = 1
	m[b2.Key()] = 2
	assert.Equal(t, 1, m[b1.Key()])
	assert.Equal(t, 2, m[b2.Key()])
}

func TestEqual(t *testing.T) {
	b1 := FromInts([]int{1, 2, 3}, 2)
	b2 := FromInts([]int{1, 2, 3}, 2)
	b3 := FromInts([]int{1, 2, 4}, 2)
	assert.True(t, b1.Equal(b2))
	assert.False(t, b1.Equal(b3))
}

func TestOr(t *testing.T) {
	b1 := FromInts([]int{0, 1}, 1)
	b2 := FromInts([]int{1, 5}, 1)
	assert.Equal(t, []int{0, 1, 5}, b1.Or(b2).ToInts())
}

func TestAnd(t *testing.T) {
	b1 := FromInts([]int{0, 1, 2}, 1)
	b2 := FromInts([]int{1, 2, 3}, 1)
	assert.Equal(t, []int{1, 2}, b1.And(b2).ToInts())
}

func TestXor(t *testing.T) {
	b1 := FromInts([]int{0, 1, 2}, 1)
	b2 := FromInts([]int{1, 2, 3}, 1)
	assert.Equal(t, []int{0, 3}, b1.Xor(b2).ToInts())
}
