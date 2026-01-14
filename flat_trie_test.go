package tries

import (
	"math/rand"
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/assertions/should"
)

func TestCreate(t *testing.T) {
	and := assertions.New(t)
	trie, err := NewFlatTrie[string, int]()
	and.So(err, should.BeNil)
	and.So(trie, should.NotBeNil)
	and.So(trie.Length(), should.Equal, 0)
}

func TestAddFind(t *testing.T) {
	trie, _ := NewFlatTrie[string, int]()

	t.Run("add hello", func(t *testing.T) {
		and := assertions.New(t)
		and.So(trie.Add("hello", 1), should.BeTrue)
	})

	t.Run("add world", func(t *testing.T) {
		and := assertions.New(t)
		and.So(trie.Add("world", 2), should.BeTrue)
	})

	t.Run("add empty", func(t *testing.T) {
		and := assertions.New(t)
		and.So(trie.Add("", 3), should.BeTrue)
	})

	t.Run("add duplicate hello", func(t *testing.T) {
		and := assertions.New(t)
		and.So(trie.Add("hello", 4), should.BeFalse)
	})

	t.Run("check length", func(t *testing.T) {
		and := assertions.New(t)
		and.So(trie.Length(), should.Equal, 3)
	})

	t.Run("find hello", func(t *testing.T) {
		and := assertions.New(t)
		value, found := trie.Find("hello")
		and.So(found, should.BeTrue)
		and.So(value, should.Equal, 4)
	})

	t.Run("find world", func(t *testing.T) {
		and := assertions.New(t)
		value, found := trie.Find("world")
		and.So(found, should.BeTrue)
		and.So(value, should.Equal, 2)
	})

	t.Run("find empty", func(t *testing.T) {
		and := assertions.New(t)
		value, found := trie.Find("")
		and.So(found, should.BeTrue)
		and.So(value, should.Equal, 3)
	})

	t.Run("find notfound", func(t *testing.T) {
		and := assertions.New(t)
		_, found := trie.Find("notfound")
		and.So(found, should.BeFalse)
	})
}

func TestGrow(t *testing.T) {
	and := assertions.New(t)
	trie, _ := NewFlatTrie[uint16, uint16]()

	const itemsToAdd = 10_000
	randomizer := rand.New(rand.NewSource(12345))
	var randomValues []uint16
	for i := 0; i < itemsToAdd; i++ {
		key := uint16(randomizer.Intn(65536))
		randomValues = append(randomValues, key)
	}

	for _, i := range randomValues {
		trie.Add(i, i)
	}

	for _, i := range randomValues {
		value, found := trie.Find(i)
		and.So(found, should.BeTrue)
		and.So(value, should.Equal, i)
	}
}
