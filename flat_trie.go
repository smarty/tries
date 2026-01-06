package tries

const (
	arenaNumber4   = 4
	arenaNumber16  = 16
	arenaNumber64  = 64
	arenaNumber256 = 256
)

type nodeLoc struct {
	arena int // 0 root, else 4, 16, 64, 256
	index int // index in arena slice
}

type FlatTrie[TKey TrieKey, TValue any] struct {
	// root node is the nil key
	root flatNode[TValue]

	converter converter[TKey]
	length    int

	// arenas contain chunks of contiguous flatNodes based on sibling count
	arena4   []flatNode[TValue]
	arena16  []flatNode[TValue]
	arena64  []flatNode[TValue]
	arena256 []flatNode[TValue]

	// useTables indicate which arena slots are used and which are free. Tables
	// are used to find free slots when adding new nodes and scale with the
	// arenas.
	useTable4   []bool
	useTable16  []bool
	useTable64  []bool
	useTable256 []bool
}

func NewFlatTrie[TKey TrieKey, TValue any](transforms ...TransformFunc) (trie Trie[TKey, TValue], err error) {
	var converter converter[TKey]
	converter, err = selectConverter[TKey]()
	if err != nil {
		return nil, err
	}

	if len(transforms) > 0 {
		converter = wrapConverter(converter, transforms)
	}

	const (
		// these represent the number of sibling chunks, not nodes
		a4   = 1 << 10 // Most common arena size for most data types, so more chunks.
		a16  = 1 << 4
		a64  = 1 << 4
		a256 = 1 << 2 // really only ever used for high volume data with even distribution
	)

	return &FlatTrie[TKey, TValue]{
		converter:   converter,
		arena4:      make([]flatNode[TValue], a4*arenaNumber4),
		useTable4:   make([]bool, a4),
		arena16:     make([]flatNode[TValue], a16*arenaNumber16),
		useTable16:  make([]bool, a16),
		arena64:     make([]flatNode[TValue], a64*arenaNumber64),
		useTable64:  make([]bool, a64),
		arena256:    make([]flatNode[TValue], a256*arenaNumber256),
		useTable256: make([]bool, a256),
	}, nil
}

func NewFlatTrieFromMap[TKey TrieIntegerString, TValue any](mapped map[TKey]TValue, transforms ...TransformFunc) (trie Trie[TKey, TValue], err error) {
	trie, err = NewFlatTrie[TKey, TValue](transforms...)
	if err != nil {
		return nil, err
	}

	for key, value := range mapped {
		trie.Add(key, value)
	}

	return trie, nil
}

func (this *FlatTrie[TKey, TValue]) Add(key TKey, value TValue) (expanded bool) {
	this.converter.Load(key)

	cur := nodeLoc{arena: 0, index: 0} // 0 == root

	for {
		k, ok := this.converter.Next()
		if !ok {
			return this.setTerminalValue(cur, value)
		}

		if child, found := this.findChild(cur, k); found {
			cur = child
			continue
		}

		cur = this.addChild(cur, k)
	}
}

func (this *FlatTrie[TKey, TValue]) Find(key TKey) (value TValue, found bool) {
	this.converter.Load(key)
	node := &this.root

	for {
		k, ok := this.converter.Next()
		if !ok {
			if node.hasValue {
				return node.value, true
			}

			return value, false
		}

		childIdx, childCount := node.firstChild, node.childCount
		if childCount == 0 {
			return value, false
		}

		arena, _ := this.selectArena(childCount)

		foundChild := false
		for i := 0; i < childCount; i++ {
			child := &arena[childIdx+i]
			if child.keyPart == k {
				node = child
				foundChild = true
				break
			}
		}

		if !foundChild {
			return value, false
		}
	}
}

func (this *FlatTrie[TKey, TValue]) Length() int {
	return this.length
}

// private methods

func (this *FlatTrie[TKey, TValue]) addChild(parent nodeLoc, k byte) nodeLoc {
	// Never keep a pointer across alloc/promote (those can expand arenas).
	node := this.nodePtr(parent)

	// 1) First child: allocate a 4-chunk, insert at [start], set parent.
	if node.childCount == 0 {
		start := this.allocateChunk(arenaNumber4)

		node = this.nodePtr(parent) // reacquire after alloc (arena growth may have occurred)
		node.firstChild = start
		node.childCount = 1

		this.arena4[start] = flatNode[TValue]{keyPart: k}
		return nodeLoc{arena: arenaNumber4, index: start}
	}

	// 2) Not first child: either insert into current chunk, or promote.
	arenaSize := node.calculateArena()

	// 2a) If chunk has capacity, insert (append for <256, sorted insert for 256).
	if node.childCount < arenaSize {
		start := node.firstChild
		count := node.childCount

		if arenaSize == arenaNumber256 {
			pos := this.insertSorted256(start, count, k)
			node = this.nodePtr(parent)
			node.childCount = count + 1
			return nodeLoc{arena: arenaNumber256, index: pos}
		}

		// append for 4/16/64
		pos := start + count
		node.childCount = count + 1

		switch arenaSize {
		case arenaNumber4:
			this.arena4[pos] = flatNode[TValue]{keyPart: k}
		case arenaNumber16:
			this.arena16[pos] = flatNode[TValue]{keyPart: k}
		default:
			this.arena64[pos] = flatNode[TValue]{keyPart: k}
		}

		return nodeLoc{arena: arenaSize, index: pos}
	}

	// 2b) Chunk is full: promote to next arena size and then insert.
	return this.promoteAndInsert(parent, arenaSize, k)
}

func (this *FlatTrie[TKey, TValue]) allocateChunk(arenaSize int) int {
	for {
		switch arenaSize {
		case arenaNumber4:
			for i := 0; i < len(this.useTable4); i++ {
				if !this.useTable4[i] {
					this.useTable4[i] = true
					start := i * arenaNumber4
					this.zeroChunk(arenaNumber4, start)
					return start
				}
			}
		case arenaNumber16:
			for i := 0; i < len(this.useTable16); i++ {
				if !this.useTable16[i] {
					this.useTable16[i] = true
					start := i * arenaNumber16
					this.zeroChunk(arenaNumber16, start)
					return start
				}
			}
		case arenaNumber64:
			for i := 0; i < len(this.useTable64); i++ {
				if !this.useTable64[i] {
					this.useTable64[i] = true
					start := i * arenaNumber64
					this.zeroChunk(arenaNumber64, start)
					return start
				}
			}
		default:
			for i := 0; i < len(this.useTable256); i++ {
				if !this.useTable256[i] {
					this.useTable256[i] = true
					start := i * arenaNumber256
					this.zeroChunk(arenaNumber256, start)
					return start
				}
			}
		}

		// No free chunk: grow this arena and retry.
		this.doubleArenaSize(arenaSize)
	}
}

func (this *FlatTrie[TKey, TValue]) arenaSlice(arenaSize int) []flatNode[TValue] {
	switch arenaSize {
	case arenaNumber4:
		return this.arena4
	case arenaNumber16:
		return this.arena16
	case arenaNumber64:
		return this.arena64
	default:
		return this.arena256
	}
}

func (this *FlatTrie[TKey, TValue]) binarySearch256(start int, count int, k byte) (idx int, found bool) {
	lo, hi := 0, count
	base := start

	for lo < hi {
		mid := (lo + hi) >> 1
		v := this.arena256[base+mid].keyPart
		if v < k {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	if lo < count && this.arena256[base+lo].keyPart == k {
		return base + lo, true
	}

	return 0, false
}

func (this *FlatTrie[TKey, TValue]) copyChunk(oldSize int, oldStart int, newSize int, newStart int, count int) {
	switch oldSize {
	case arenaNumber4:
		src := this.arena4[oldStart : oldStart+count]
		switch newSize {
		case arenaNumber16:
			copy(this.arena16[newStart:newStart+count], src)
		case arenaNumber64:
			copy(this.arena64[newStart:newStart+count], src)
		default:
			copy(this.arena256[newStart:newStart+count], src)
		}
	case arenaNumber16:
		src := this.arena16[oldStart : oldStart+count]
		switch newSize {
		case arenaNumber64:
			copy(this.arena64[newStart:newStart+count], src)
		default:
			copy(this.arena256[newStart:newStart+count], src)
		}
	default: // 64 -> 256 only
		src := this.arena64[oldStart : oldStart+count]
		copy(this.arena256[newStart:newStart+count], src)
	}
}

func (this *FlatTrie[TKey, TValue]) doubleArenaSize(arena int) {
	switch arena {
	case arenaNumber4:
		this.arena4, this.useTable4 = expandArena(this.arena4, this.useTable4)
	case arenaNumber16:
		this.arena16, this.useTable16 = expandArena(this.arena16, this.useTable16)
	case arenaNumber64:
		this.arena64, this.useTable64 = expandArena(this.arena64, this.useTable64)
	case arenaNumber256:
		this.arena256, this.useTable256 = expandArena(this.arena256, this.useTable256)
	}
}

func (this *FlatTrie[TKey, TValue]) findChild(parent nodeLoc, k byte) (child nodeLoc, found bool) {
	node := this.nodePtr(parent)
	if node.childCount == 0 {
		return nodeLoc{}, false
	}

	arenaSize := node.calculateArena()
	start := node.firstChild
	count := node.childCount

	if arenaSize == arenaNumber256 {
		idx, ok := this.binarySearch256(start, count, k)
		if !ok {
			return nodeLoc{}, false
		}
		return nodeLoc{arena: arenaNumber256, index: idx}, true
	}

	arena := this.arenaSlice(arenaSize)
	limit := start + count
	for i := start; i < limit; i++ {
		if arena[i].keyPart == k {
			return nodeLoc{arena: arenaSize, index: i}, true
		}
	}

	return nodeLoc{}, false
}

func (this *FlatTrie[TKey, TValue]) freeChunk(arenaSize int, start int) {
	chunkIdx := start / arenaSize
	this.zeroChunk(arenaSize, start)

	switch arenaSize {
	case arenaNumber4:
		this.useTable4[chunkIdx] = false
	case arenaNumber16:
		this.useTable16[chunkIdx] = false
	case arenaNumber64:
		this.useTable64[chunkIdx] = false
	default:
		this.useTable256[chunkIdx] = false
	}
}

func (this *FlatTrie[TKey, TValue]) insertSorted256(start int, count int, k byte) int {
	lo, hi := 0, count
	base := start
	for lo < hi {
		mid := (lo + hi) >> 1
		if this.arena256[base+mid].keyPart < k {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	// Shift nodes after insert right by 1.
	for i := count; i > lo; i-- {
		this.arena256[base+i] = this.arena256[base+i-1]
	}

	this.arena256[base+lo] = flatNode[TValue]{keyPart: k}
	return base + lo
}

func nextArenaSize(current int) int {
	switch current {
	case arenaNumber4:
		return arenaNumber16
	case arenaNumber16:
		return arenaNumber64
	default:
		return arenaNumber256
	}
}

func (this *FlatTrie[TKey, TValue]) nodePtr(l nodeLoc) *flatNode[TValue] {
	if l.arena == 0 {
		return &this.root
	}
	switch l.arena {
	case arenaNumber4:
		return &this.arena4[l.index]
	case arenaNumber16:
		return &this.arena16[l.index]
	case arenaNumber64:
		return &this.arena64[l.index]
	default:
		return &this.arena256[l.index]
	}
}

func (this *FlatTrie[TKey, TValue]) promoteAndInsert(parent nodeLoc, oldSize int, k byte) nodeLoc {
	node := this.nodePtr(parent)
	oldStart := node.firstChild
	oldCount := node.childCount

	newSize := nextArenaSize(oldSize)
	newStart := this.allocateChunk(newSize)

	// Copy siblings into new chunk.
	this.copyChunk(oldSize, oldStart, newSize, newStart, oldCount)

	// Free old chunk (and zero to drop references).
	this.freeChunk(oldSize, oldStart)

	// If we promoted into 256, sort the copied portion once.
	if newSize == arenaNumber256 {
		this.sortChunk256(newStart, oldCount)
	}

	// Update parent (reacquire in case alloc expanded).
	node = this.nodePtr(parent)
	node.firstChild = newStart
	node.childCount = oldCount // will increment after insertion below

	// Insert the new child.
	if newSize == arenaNumber256 {
		pos := this.insertSorted256(newStart, oldCount, k)
		node = this.nodePtr(parent)
		node.childCount = oldCount + 1
		return nodeLoc{arena: arenaNumber256, index: pos}
	}

	// append for 16/64 (promotion target can only be 16 or 64 here)
	pos := newStart + oldCount
	switch newSize {
	case arenaNumber16:
		this.arena16[pos] = flatNode[TValue]{keyPart: k}
	case arenaNumber64:
		this.arena64[pos] = flatNode[TValue]{keyPart: k}
	}

	node = this.nodePtr(parent)
	node.childCount = oldCount + 1
	return nodeLoc{arena: newSize, index: pos}
}

func (this *FlatTrie[TKey, TValue]) selectArena(childCount int) (arena []flatNode[TValue], useTable []bool) {
	switch {
	case childCount <= arenaNumber4:
		return this.arena4, this.useTable4
	case childCount <= arenaNumber16:
		return this.arena16, this.useTable16
	case childCount <= arenaNumber64:
		return this.arena64, this.useTable64
	default:
		return this.arena256, this.useTable256
	}
}

func (this *FlatTrie[TKey, TValue]) setTerminalValue(at nodeLoc, value TValue) (expanded bool) {
	node := this.nodePtr(at)
	if !node.hasValue {
		this.length++
		node.hasValue = true
		node.value = value
		return true
	}
	node.value = value
	return false
}

func (this *FlatTrie[TKey, TValue]) sortChunk256(start int, count int) {
	// Insertion sort is fine here: count will be <= 64 when promoting into 256.
	chunkBase := start
	for unsortedIndex := 1; unsortedIndex < count; unsortedIndex++ {
		unsortedNode := this.arena256[chunkBase+unsortedIndex]
		sortedIndex := unsortedIndex - 1
		for sortedIndex >= 0 && this.arena256[chunkBase+sortedIndex].keyPart > unsortedNode.keyPart {
			this.arena256[chunkBase+sortedIndex+1] = this.arena256[chunkBase+sortedIndex]
			sortedIndex--
		}

		this.arena256[chunkBase+sortedIndex+1] = unsortedNode
	}
}

func (this *FlatTrie[TKey, TValue]) zeroChunk(arenaSize int, start int) {
	switch arenaSize {
	case arenaNumber4:
		for i := 0; i < arenaNumber4; i++ {
			this.arena4[start+i] = flatNode[TValue]{}
		}
	case arenaNumber16:
		for i := 0; i < arenaNumber16; i++ {
			this.arena16[start+i] = flatNode[TValue]{}
		}
	case arenaNumber64:
		for i := 0; i < arenaNumber64; i++ {
			this.arena64[start+i] = flatNode[TValue]{}
		}
	default:
		for i := 0; i < arenaNumber256; i++ {
			this.arena256[start+i] = flatNode[TValue]{}
		}
	}
}

// private functions

func expandArena[TValue any](arena []flatNode[TValue], useTable []bool) ([]flatNode[TValue], []bool) {
	const doubleFactor = 2

	oldSize := len(arena)
	newSize := oldSize * doubleFactor
	newArena := make([]flatNode[TValue], newSize)
	copy(newArena, arena)

	oldSize = len(useTable)
	newSize = oldSize * doubleFactor
	newUseTable := make([]bool, newSize)
	copy(newUseTable, useTable)
	return newArena, newUseTable
}
