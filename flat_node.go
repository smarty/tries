package tries

// flatNode represents a node in the FlatTrie array.
type flatNode[TValue any] struct {
	keyPart    byte
	hasValue   bool
	value      TValue
	firstChild int
	childCount int
}

func (this flatNode[TValue]) calculateArena() int {
	switch {
	case this.childCount <= 4:
		return 4
	case this.childCount <= 16:
		return 16
	case this.childCount <= 64:
		return 64
	default:
		return 256
	}
}
