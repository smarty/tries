package tries

type SimpleTrie[TKey TrieKey, TValue any] struct {
	converter converter[TKey]
	head      simpleNode[TKey, TValue] // head is empty, or the nil key
	length    int
}

func NewTrie[TKey TrieKey, TValue any](transforms ...TransformFunc) (trie Trie[TKey, TValue], err error) {
	var converter converter[TKey]
	converter, err = selectConverter[TKey]()
	if err != nil {
		return nil, err
	}

	if len(transforms) > 0 {
		converter = wrapConverter(converter, transforms)
	}

	return &SimpleTrie[TKey, TValue]{
		converter: converter,
	}, nil
}

func NewTrieFromMap[TKey TrieIntegerString, TValue any](mapped map[TKey]TValue, transforms ...TransformFunc) (trie Trie[TKey, TValue], err error) {
	trie, err = NewTrie[TKey, TValue](transforms...)
	if err != nil {
		return nil, err
	}

	for key, value := range mapped {
		trie.Add(key, value)
	}

	return trie, nil
}

func (this *SimpleTrie[TKey, TValue]) Add(key TKey, value TValue) (expanded bool) {
	this.converter.Load(key)
	expanded = this.head.add(this.converter, value)
	if expanded {
		this.length++
	}

	return expanded
}

func (this *SimpleTrie[TKey, TValue]) Find(key TKey) (value TValue, found bool) {
	this.converter.Load(key)
	return this.head.Find(this.converter)
}

func (this *SimpleTrie[TKey, TValue]) Length() (length int) { // LengthMutexed
	return this.length
}
