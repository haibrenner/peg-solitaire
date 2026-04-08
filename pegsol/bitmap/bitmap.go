package bitmap

type Bitmap uint64

func FromInts(vals []int) Bitmap {
	var b Bitmap
	for _, v := range vals {
		b |= 1 << uint(v)
	}
	return b
}

func (b Bitmap) ToInts() []int {
	var result []int
	for bit := 0; bit < 64; bit++ {
		if b&(1<<uint(bit)) != 0 {
			result = append(result, bit)
		}
	}
	return result
}
