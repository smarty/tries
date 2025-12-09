package tries

type (
	// TrieInteger defines any integer types that can be used as a key type for
	// a [Trie].
	TrieInteger interface {
		~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~uintptr | ~int8 | ~int16 | ~int32 | ~int64 | ~int
	}

	// TrieString defines any string types that can be used as a key type for a
	// [Trie].
	TrieString interface {
		~string
	}

	// TrieSlice defines any slice types that can be used as a key type for a
	// [Trie].
	TrieSlice interface {
		~[]uint8 | ~[]uint16 | ~[]uint32 | ~[]uint64 | ~[]uint | ~[]uintptr | ~[]int8 | ~[]int16 | ~[]int32 | ~[]int64 | ~[]int
	}

	// TrieIntegerString marries the [TrieInteger] and [TrieString] together.
	TrieIntegerString interface {
		TrieInteger | TrieString
	}

	// TrieKey defines all types that can be used as a key type for a [Trie].
	TrieKey interface {
		TrieIntegerString | TrieSlice
	}

	Trie[TKey TrieKey, TValue any] interface {
		// Add inserts a new key-value pair, overwriting any extant value if the
		// key is already present.
		//
		// Parameters:
		//   - key is the key to associate the new value with.
		//   - value is the new value to be stored alongside the key.
		//
		// Returns:
		//   - expanded is `true` if this operation created a new entry or
		//     `false` if it replaced a value.
		Add(key TKey, value TValue) (expanded bool)

		// Find looks for the provided key, and if found, returns the associated
		// value.
		//
		// Parameters:
		//   - key is the lookup key.
		//
		// Returns:
		//   - value is the found value, or the zero value if not found.
		//   - found is `true` if the key was found or `false` otherwise.
		Find(key TKey) (value TValue, found bool)

		// Length returns the current number of key-value pairs stored in this
		// [Trie].
		Length() (length int)
	}

	// TransformFunc is used to transform a [TrieKey] for any normalization
	// processes when performing a store or retrieval operation. Normalization
	// is performed using a keyhole approach (one byte at a time with no context).
	//
	// Parameters:
	//   - in is the current byte to analyze.
	//
	// Returns:
	//   - out is the transformed byte value.
	//   - use is set to `true` if the output should be used in composing a key
	//     value or comparing key values, or `false` to indicate that this byte
	//     should be ignored entirely.
	TransformFunc func(in byte) (out byte, use bool)
)
