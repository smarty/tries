package tries

import "errors"

// ! important ! !!! Move this to contracts.go when done editing and before commit !!!

var (
	ErrorBadTrieKey = errors.New("unable to create Trie with bad key type")
)
