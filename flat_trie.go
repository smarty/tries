package tries

import "github.com/smarty/tries/internal"

type FlatTrie[TKey TrieKey, TValue any] struct {
	converter converter[TKey]
	length    int

	head      flatNodeHeader
	emptyBody flatNodeBody[TValue] // head is empty, or the nil key

	bodies  []flatNodeBody[TValue]
	headers []flatNodeHeader

	buddyAllocator internal.BuddyAllocator
}

func NewFlatTrie[TKey TrieKey, TValue any](transforms ...TransformFunc) (trie Trie[TKey, TValue], err error) {
	const startingChunkCount = 1 << 10

	var converter converter[TKey]
	converter, err = selectConverter[TKey]()
	if err != nil {
		return nil, err
	}

	if len(transforms) > 0 {
		converter = wrapConverter(converter, transforms)
	}

	return &FlatTrie[TKey, TValue]{
		converter:      converter,
		bodies:         make([]flatNodeBody[TValue], startingChunkCount*flatHeaderChunkSize),
		headers:        make([]flatNodeHeader, startingChunkCount*flatHeaderChunkSize),
		buddyAllocator: *internal.NewBuddyAllocator(startingChunkCount, true),
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
	k, ok := this.converter.Next()
	if !ok {
		expanded = !this.emptyBody.hasValue
		this.emptyBody.hasValue = true
		this.emptyBody.value = value

		if expanded {
			this.length++
		}

		return expanded
	}

	nodeIndex := int64(-1)
	for ok {
		nextIndex := this.nextNode(nodeIndex, k)
		if nextIndex == -1 {
			nextIndex = this.insert(nodeIndex, k)
		}

		k, ok = this.converter.Next()
		nodeIndex = nextIndex
	}

	node := &this.bodies[nodeIndex]
	expanded = !node.hasValue
	node.hasValue = true
	node.value = value
	if expanded {
		this.length++
	}

	return expanded
}

func (this *FlatTrie[TKey, TValue]) Find(key TKey) (value TValue, found bool) {
	this.converter.Load(key)
	k, ok := this.converter.Next()
	if !ok {
		if this.emptyBody.hasValue {
			return this.emptyBody.value, true
		}

		return value, false
	}

	nodeIndex := int64(-1)
	for ok {
		nodeIndex = this.nextNode(nodeIndex, k)
		if nodeIndex == -1 {
			return value, false
		}

		k, ok = this.converter.Next()
	}

	node := this.bodies[nodeIndex]
	if !node.hasValue {
		return value, false
	}

	return node.value, true
}

func (this *FlatTrie[TKey, TValue]) Length() (length int) {
	return this.length
}

// --- private methods --- //

func (this *FlatTrie[TKey, TValue]) expand(node *flatNodeHeader, nodeIndex int64, free bool) {
	currentChunkCount := node.currentChunkCount()
	requiredChunkCount := currentChunkCount + 1

	if node.childCount > 0 && free {
		this.buddyAllocator.Free(int(node.firstChildChunk), int(currentChunkCount))
	}

	newFirstChunk, ok := this.buddyAllocator.Allocate(int(requiredChunkCount))
	if !ok {
		newSize := this.buddyAllocator.Grow()
		this.headers = expandSliceToSize(this.headers, newSize*flatHeaderChunkSize)
		this.bodies = expandSliceToSize(this.bodies, newSize*flatHeaderChunkSize)

		if nodeIndex == -1 {
			node = &this.head
		} else {
			node = &this.headers[nodeIndex]
		}

		this.expand(node, nodeIndex, false)
		return
	}

	oldFirstChunk := node.firstChildChunk
	oldIndex := int64(oldFirstChunk * flatHeaderChunkSize)
	newIndex := int64(newFirstChunk * flatHeaderChunkSize)

	childCount := int64(node.childCount)
	if childCount > 0 {
		copy(this.headers[newIndex:newIndex+childCount], this.headers[oldIndex:oldIndex+childCount])
		copy(this.bodies[newIndex:newIndex+childCount], this.bodies[oldIndex:oldIndex+childCount])
	}

	node.firstChildChunk = uint32(newFirstChunk)
}

func (this *FlatTrie[TKey, TValue]) insert(current int64, key uint8) (newIndex int64) {
	var node *flatNodeHeader
	if current == -1 {
		node = &this.head
	} else {
		node = &this.headers[current]
	}

	if node.readyToExpand() {
		this.expand(node, current, true)
		if current == -1 {
			node = &this.head
		} else {
			node = &this.headers[current]
		}
	}

	newIndex = this.insertNode(node, key)
	return newIndex
}

func (this *FlatTrie[TKey, TValue]) insertNode(node *flatNodeHeader, key uint8) (newIndex int64) {
	newIndex = int64(node.firstChildChunk * flatHeaderChunkSize)
	end := newIndex + int64(node.childCount)
	for ; newIndex < end; newIndex++ {
		if this.headers[newIndex].key > key {
			break
		}
	}

	if newIndex < end {
		copy(this.headers[newIndex+1:end+1], this.headers[newIndex:end])
		copy(this.bodies[newIndex+1:end+1], this.bodies[newIndex:end])
	}

	this.headers[newIndex] = flatNodeHeader{
		key:             key,
		childCount:      0,
		firstChildChunk: 0,
	}

	this.bodies[newIndex] = flatNodeBody[TValue]{}
	node.childCount++

	return newIndex
}

func (this *FlatTrie[TKey, TValue]) nextNode(current int64, key uint8) (next int64) {
	const scanThreshold = 7

	node := this.head
	if current > -1 {
		node = this.headers[current]
	}

	if node.childCount == 0 {
		return -1
	}

	first := int64(node.firstChildChunk * flatHeaderChunkSize)
	end := first + int64(node.childCount)
	_ = this.headers[end-1] // touch last element to avoid bounds check
	if node.childCount < scanThreshold {
		for i := first; i < end; i++ {
			if this.headers[i].key == key {
				return i
			}
		}

		return -1
	}

	low := first
	high := end
	for low < high {
		mid := ((high - low) >> 1) + low
		midNode := this.headers[mid]
		if midNode.key == key {
			return mid
		}

		if key < midNode.key {
			high = mid
			continue
		}

		low = mid + 1
	}

	return -1
}

// --- private functions --- //

func expandSlice[T any](slice []T) (newSlice []T) {
	const factor = 1.2
	newSlice = make([]T, int(float64(len(slice))*factor))
	copy(newSlice, slice)
	return newSlice
}

func expandSliceToSize[T any](slice []T, size int) (newSlice []T) {
	newSlice = make([]T, size)
	copy(newSlice, slice)
	return newSlice
}
