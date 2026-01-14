package tries

const (
	flatHeaderChunkSize = 8
	flatHeaderChunkMask = flatHeaderChunkSize - 1
)

type flatNodeHeader struct { // 8 bytes, 8 per scan-line
	key             byte
	childCount      uint16
	firstChildChunk uint32
}

type flatNodeBody[T any] struct {
	hasValue bool
	value    T
}

func (this *flatNodeHeader) currentChunkCount() (count byte) {
	return byte((uint16(this.childCount) + flatHeaderChunkMask) / flatHeaderChunkSize)
}

func (this *flatNodeHeader) readyToExpand() bool {
	return (this.childCount & flatHeaderChunkMask) == 0
}
