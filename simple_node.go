package tries

type simpleNode[TKey TrieKey, TValue any] struct {
	hasValue bool
	value    TValue
	key      uint8
	next     []simpleNode[TKey, TValue]
}

func (this *simpleNode[TKey, TValue]) Find(key converter[TKey]) (value TValue, ok bool) {
	var k uint8
	k, ok = key.Next()
	if !ok {
		if this.hasValue {
			return this.value, true
		}

		return value, false
	}

	nextNode, found := this.binarySearchNext(k)
	if !found {
		return value, false
	}

	return nextNode.Find(key)
}

func (this *simpleNode[TKey, TValue]) add(key converter[TKey], value TValue) bool {
	k, ok := key.Next()
	if !ok {
		expanded := !this.hasValue
		this.hasValue = true
		this.value = value
		return expanded
	}

	nextNode, found := this.binarySearchNext(k)
	if found {
		return nextNode.add(key, value)
	}

	nextNode = this.insertNewNode(k)
	return nextNode.add(key, value)
}

func (this *simpleNode[TKey, TValue]) insertNewNode(key uint8) *simpleNode[TKey, TValue] {
	if len(this.next) == 0 {
		this.next = append(this.next, simpleNode[TKey, TValue]{key: key})
		return &this.next[0]
	}

	this.next = append(this.next, simpleNode[TKey, TValue]{})
	index := len(this.next) - 2
	for ; index >= 0; index-- {
		if this.next[index].key > key {
			this.next[index+1] = this.next[index]
			continue
		}

		break
	}

	index++
	this.next[index] = simpleNode[TKey, TValue]{key: key}
	return &this.next[index]
}

func (this *simpleNode[TKey, TValue]) binarySearchNext(key uint8) (nextNode *simpleNode[TKey, TValue], found bool) {
	bottom := 0
	top := len(this.next) - 1
	for top >= bottom {
		index := ((top - bottom) / 2) + bottom
		nextNode = &this.next[index]
		test := nextNode.key
		if test == key {
			return nextNode, true
		}

		if test > key {
			top = index - 1
			continue
		}

		bottom = index + 1
	}

	return nextNode, false
}
