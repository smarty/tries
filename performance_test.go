package tries

import (
	"iter"
	"math/rand"
	"strings"
	"testing"

	"github.com/smarty/benchy"
	"github.com/smarty/benchy/options"
	"github.com/smarty/benchy/providers"
)

func BenchmarkTries_FindTextSimple(b *testing.B) {
	provider := providers.New1(func(string) {})
	simpleTrie, _ := NewTrie[string, bool]()
	flatTrie, _ := NewFlatTrie[string, bool]()
	rawMap := make(map[string]bool)
	for _, v := range englishWords {
		simpleTrie.Add(v, true)
		flatTrie.Add(v, true)
		rawMap[v] = true
		provider.Add(v)
	}

	benchy.New(b, options.Medium).
		RegisterBenchmark("simple trie", provider.WrapBenchmarkFunc(func(s string) {
			value, found := simpleTrie.Find(s)
			if !found || !value {
				b.Fatalf("Value %s not found in simple trie", s)
			}

			if value != true {
				b.Fatalf("Value %s incorrect in simple trie", s)
			}
		})).
		RegisterBenchmark("flat trie", provider.WrapBenchmarkFunc(func(s string) {
			value, found := flatTrie.Find(s)
			if !found || !value {
				b.Fatalf("Value %s not found in flat trie", s)
			}

			if value != true {
				b.Fatalf("Value %s incorrect in flat trie", s)
			}
		})).
		RegisterBenchmark("raw map", provider.WrapBenchmarkFunc(func(s string) {
			value, found := rawMap[s]
			if !found || !value {
				b.Fatalf("Value %s not found in raw map", s)
			}

			if value != true {
				b.Fatalf("Value %s incorrect in raw map", s)
			}
		})).
		Run()
}

func BenchmarkTries_FindTextProcessed(b *testing.B) {
	trieRules := func(in byte) (out byte, use bool) {
		if in >= 'A' && in <= 'Z' {
			return in + 32, true
		}

		if in == '_' {
			return 0, false
		}

		return in, true
	}

	provider := providers.New1(func(string) {})
	simpleTrie, _ := NewTrie[string, bool](trieRules)
	flatTrie, _ := NewFlatTrie[string, bool](trieRules)
	rawMap := make(map[string]bool)
	for _, toStore := range englishWords {
		toSearch := strings.ReplaceAll(strings.ToUpper(toStore), "A", "A_")
		simpleTrie.Add(toStore, true)
		flatTrie.Add(toStore, true)
		rawMap[toStore] = true
		provider.Add(toSearch)
	}

	benchy.New(b, options.Medium).
		RegisterBenchmark("simple trie", provider.WrapBenchmarkFunc(func(s string) {
			value, found := simpleTrie.Find(s)
			if !found || !value {
				b.Fatalf("Value %s not found in simple trie", s)
			}

			if value != true {
				b.Fatalf("Value %s incorrect in simple trie", s)
			}
		})).
		RegisterBenchmark("flat trie", provider.WrapBenchmarkFunc(func(s string) {
			value, found := flatTrie.Find(s)
			if !found || !value {
				b.Fatalf("Value %s not found in flat trie", s)
			}

			if value != true {
				b.Fatalf("Value %s incorrect in flat trie", s)
			}
		})).
		RegisterBenchmark("raw map", provider.WrapBenchmarkFunc(func(s string) {
			value, found := rawMap[strings.ToLower(strings.ReplaceAll(s, "A_", "A"))]
			if !found || !value {
				b.Fatalf("Value %s not found in raw map", s)
			}

			if value != true {
				b.Fatalf("Value %s incorrect in raw map", s)
			}
		})).
		Run()
}

func BenchmarkTries_FindNumber(b *testing.B) {
	const valuesCount = 1_000_000

	provider := providers.New1(func(int64) {})
	simpleTrie, _ := NewTrie[int64, bool]()
	flatTrie, _ := NewFlatTrie[int64, bool]()
	rawMap := make(map[int64]bool)
	for v := range produceValues(valuesCount) {
		simpleTrie.Add(v, true)
		flatTrie.Add(v, true)
		rawMap[v] = true
		provider.Add(v)
	}

	benchy.New(b, options.Medium).
		RegisterBenchmark("simple trie", provider.WrapBenchmarkFunc(func(i int64) {
			value, found := simpleTrie.Find(i)
			if !found || !value {
				b.Fatalf("Value %d not found in simple trie", i)
			}

			if value != true {
				b.Fatalf("Value %d incorrect in simple trie", i)
			}
		})).
		RegisterBenchmark("flat trie", provider.WrapBenchmarkFunc(func(i int64) {
			value, found := flatTrie.Find(i)
			if !found || !value {
				b.Fatalf("Value %d not found in flat trie", i)
			}

			if value != true {
				b.Fatalf("Value %d incorrect in flat trie", i)
			}
		})).
		RegisterBenchmark("raw map", provider.WrapBenchmarkFunc(func(i int64) {
			value, found := rawMap[i]
			if !found || !value {
				b.Fatalf("Value %d not found in raw map", i)
			}

			if value != true {
				b.Fatalf("Value %d incorrect in raw map", i)
			}
		})).
		Run()
}

func BenchmarkAdd(b *testing.B) {
	n := 5_000_000

	benchy.New(b, options.Medium).
		RegisterBenchmark("basic trie", func() {
			trie, _ := NewTrie[int64, int]()
			for value := range produceValues(n) {
				trie.Add(value, int(value))
			}
		}).
		RegisterBenchmark("flat trie", func() {
			trie, _ := NewFlatTrie[int64, int]()
			for value := range produceValues(n) {
				trie.Add(value, int(value))
			}
		}).
		RegisterBenchmark("map", func() {
			m := make(map[int64]int)
			for value := range produceValues(n) {
				m[value] = int(value)
			}
		}).
		Run()
}

func produceValues(count int) iter.Seq[int64] {
	return func(yield func(int64) bool) {
		randomizer := rand.New(rand.NewSource(12345))
		seen := make(map[int64]struct{})
		for i := 0; i < count; i++ {
			n := int64(randomizer.Intn(count * 10))
			for {
				if _, exists := seen[n]; !exists {
					seen[n] = struct{}{}
					break
				}

				n = (n + 1) % int64(count*10)
			}

			if !yield(n) {
				return
			}
		}
	}
}
