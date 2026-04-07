package bitmap

import "unsafe"

type Bitmap []uint64
type BitmapKey string

func NewBitmap(size int) Bitmap {
	return make(Bitmap, size)
}

// FromInts sets bits corresponding to each value in the input slice.
func FromInts(vals []int, size int) Bitmap {
	b := make(Bitmap, size)
	for _, v := range vals {
		b[v/64] |= 1 << uint(v%64)
	}
	return b
}

// ToInts returns the positions of all set bits.
func (b Bitmap) ToInts() []int {
	var result []int
	for i, word := range b {
		for bit := 0; bit < 64; bit++ {
			if word&(1<<uint(bit)) != 0 {
				result = append(result, i*64+bit)
			}
		}
	}
	return result
}

func (b Bitmap) Key() BitmapKey {
	return BitmapKey(unsafe.Slice((*byte)(unsafe.Pointer(&b[0])), len(b)*8))
}

func KeyToBitmap(k BitmapKey, size int) Bitmap {
	b := make(Bitmap, size)
	copy(unsafe.Slice((*byte)(unsafe.Pointer(&b[0])), size*8), k)
	return b
}

func (b Bitmap) Equal(other Bitmap) bool {
	for i := range b {
		if b[i] != other[i] {
			return false
		}
	}
	return true
}

func (b Bitmap) Or(other Bitmap) Bitmap {
	result := make(Bitmap, len(b))
	for i := range b {
		result[i] = b[i] | other[i]
	}
	return result
}

func (b Bitmap) And(other Bitmap) Bitmap {
	result := make(Bitmap, len(b))
	for i := range b {
		result[i] = b[i] & other[i]
	}
	return result
}

func (b Bitmap) Xor(other Bitmap) Bitmap {
	result := make(Bitmap, len(b))
	for i := range b {
		result[i] = b[i] ^ other[i]
	}
	return result
}
