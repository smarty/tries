# Tries

A fast, generic trie (prefix tree) data structure implementation in Go with support for multiple key types and optional key transformations.

## Features

- **Generic implementation**: Supports any integer type, string type, or slice type as keys
- **Type-safe**: Leverages Go 1.18+ generics for compile-time type safety
- **Flexible**: Add optional transformation functions for key normalization (e.g., case-insensitive matching)
- **Efficient**: Tree-based structure optimized for prefix-based lookups and storage
- **Simple API**: Clean, intuitive interface with just three main operations

## Installation

```bash
go get github.com/smarty/tries
```

## Quick Start

```go
package main

import (
    "github.com/smarty/tries"
)

func main() {
    // Create a trie with string keys and integer values
    trie, err := tries.NewTrie[string, int]()
    if err != nil {
        panic(err)
    }

    // Add some key-value pairs
    isNew := trie.Add("hello", 1)    // isNew = true
    isNew = trie.Add("world", 2)     // isNew = true
    isNew = trie.Add("hello", 10)    // isNew = false (replaced value)

    // Look up values
    value, found := trie.Find("hello") // value = 10, found = true
    value, found = trie.Find("foo")    // value = 0, found = false

    // Get the number of entries
    length := trie.Length() // length = 2
}
```

## Supported Key Types

The trie supports the following key types:

### Integer Types
- Signed: `int`, `int8`, `int16`, `int32`, `int64`
- Unsigned: `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `uintptr`

### String Types
- `string`

### Slice Types
- `[]int`, `[]int8`, `[]int16`, `[]int32`, `[]int64`
- `[]uint`, `[]uint8`, `[]uint16`, `[]uint32`, `[]uint64`, `[]uintptr`

You can also use custom types that are defined as aliases to any of the above types.

## API

### Creating a Trie

```go
// Create a new empty trie
trie, err := tries.NewTrie[TKey, TValue]()

// Create a trie from an existing map
trie, err := tries.NewTrieFromMap(map[string]int{
    "key1": 100,
    "key2": 200,
})
```

### Adding Values

```go
// Add or update a key-value pair
// Returns true if the key was new, false if it replaced an existing value
expanded := trie.Add(key, value)
```

### Finding Values

```go
// Look up a value by key
// Returns the value and a boolean indicating if the key was found
value, found := trie.Find(key)
if found {
    // Use value
}
```

### Getting Size

```go
// Get the number of key-value pairs stored
count := trie.Length()
```

## Advanced Features

### Key Transformation

Apply transformation functions to normalize keys during storage and retrieval operations. This is useful for case-insensitive matching, whitespace normalization, or other preprocessing.

```go
// Create a case-insensitive string trie
lowerTransform := func(in byte) (out byte, use bool) {
    if in >= 'A' && in <= 'Z' {
        return in + 32, true // Convert to lowercase
    }
    return in, true
}

trie, err := tries.NewTrie[string, int](lowerTransform)
if err != nil {
    panic(err)
}

trie.Add("Hello", 1)
value, found := trie.Find("hello") // found = true, value = 1
value, found = trie.Find("HELLO")  // found = true, value = 1
```

Multiple transforms can be chained by passing multiple `TransformFunc` arguments:

```go
trie, err := tries.NewTrie[string, int](transform1, transform2, transform3)
```

### Transform Function Details

Transform functions are applied byte-by-byte with no context about surrounding bytes. Each transform function should:

- Accept a byte as input
- Return the transformed byte and a boolean indicating whether to include it
- Return `use = false` to skip the byte entirely (useful for filtering)

```go
type TransformFunc func(in byte) (out byte, use bool)
```

## Example: URL Path Matching

```go
package main

import "github.com/smarty/tries"

func main() {
    // Create a case-insensitive trie for URL paths
    caseLower := func(in byte) (out byte, use bool) {
        if in >= 'A' && in <= 'Z' {
            return in + 32, true
        }
        return in, true
    }

    trie, _ := tries.NewTrie[string, string](caseLower)

    // Store routes
    trie.Add("api/users", "user_handler")
    trie.Add("api/posts", "post_handler")
    trie.Add("api/comments", "comment_handler")

    // Look up routes
    handler, found := trie.Find("API/Users") // found = true
}
```

## Roadmap

Future enhancements planned for this library:

- **Single-slice trie** - Alternative implementation using a single slice rather than nodes for encoding keys and lookup locations, enabling much faster retrievals and reduced memory overhead
- **Iteration/traversal** - Methods to iterate over all entries or entries matching a common prefix
- **Deletion** - Remove key-value pairs from the trie
- **Prefix matching** - Find all values matching a key prefix
- **Serialization** - Save and load trie state to/from disk for persistence
- **Benchmarking suite** - More comprehensive performance comparisons and optimization

## License

Please see the LICENSE file for licensing information.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
