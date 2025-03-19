package identifiers

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/google/uuid"
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

func CalcContainerHashFNV(customerGUID, cluster, podName, containerName, namespace string) string {
	strLower := strings.ToLower(fmt.Sprintf("%s/%s/%s/%s/%s", customerGUID, cluster, podName, containerName, namespace))
	return CalcHashFNV(strLower)
}

func GenerateExceptionUID() (string, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return newUUID.String(), nil
}

// ConvertResourceIDToResourceHashFNV expects to get resourceID in the format of `apiVersion/namespace/kind/name`
// for e.g `apps/v1/default/Deployment/deploymenttest1`
func ConvertResourceIDToResourceHashFNV(customerGUID, clusterName, resourceID string) string {
	parts := strings.Split(resourceID, "/")
	if len(parts) < 4 {
		fmt.Println("Invalid resourceID format. Expected format: apiVersion/namespace/kind/name")
		return ""
	}
	// Adjust the apiVersion to remove the leading '/' if present
	apiVersion := strings.TrimPrefix(strings.Join(parts[:len(parts)-3], "/"), "/")
	namespace := parts[len(parts)-3]
	kind := parts[len(parts)-2]
	name := parts[len(parts)-1]

	return CalcResourceHashFNV(customerGUID, clusterName, kind, name, namespace, apiVersion)
}
