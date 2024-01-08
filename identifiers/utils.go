package identifiers

import (
	"fmt"
	"hash/fnv"
	"strings"
)

// CalcHashFNV calculates the hash (FNV) of the string
func CalcHashFNV(id string) string {
	hasher := fnv.New64a()
	hasher.Write([]byte(id))
	return fmt.Sprintf("%v", hasher.Sum64())
}

func CalcResourceHashFNV(customerGUID, cluster, kind, name, namespace, apiVersion string) string {
	strLower := strings.ToLower(fmt.Sprintf("%s/%s/%s/%s/%s/%s", customerGUID, cluster, kind, name, namespace, apiVersion))
	return CalcHashFNV(strLower)

}
