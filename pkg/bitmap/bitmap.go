package bitmap

type Bitmap struct {
	bits []byte
	size int
}

func NewBitmap(size int) *Bitmap {
	if size <= 0 {
		size = 250
	}
	return &Bitmap{
		bits: make([]byte, size),
		size: size * 8,
	}
}

func (b *Bitmap) Set(id string) {
	idx := hash(id) % b.size
	byteIdx := idx / 8
	bitIdx := idx % 8
	b.bits[byteIdx] |= 1 << bitIdx
}

func (b *Bitmap) IsSet(id string) bool {
	idx := hash(id) % b.size
	byteIdx := idx / 8
	bitIdx := idx % 8

	return (b.bits[byteIdx] & 1 << bitIdx) != 0
}

func (b *Bitmap) Export() []byte {
	return b.bits
}

func Load(bits []byte) *Bitmap {
	if len(bits) == 0 {
		return nil
	}
	return &Bitmap{
		bits: bits,
		size: len(bits) * 8,
	}
}

func hash(id string) int {
	seed := 131313
	h := 0
	for _, c := range id {
		h = h*seed + int(c)
	}
	return h & 0x7FFFFFFF
}
