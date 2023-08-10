package identifiers

import (
	"fmt"
	"hash/fnv"
)

// CalcHashFNV calculates the hash (FNV) of the string
func CalcHashFNV(id string) string {
	hasher := fnv.New64a()
	hasher.Write([]byte(id))
	return fmt.Sprintf("%v", hasher.Sum64())
}
