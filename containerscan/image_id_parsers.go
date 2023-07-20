package containerscan

import (
	"strings"
)

func getRegistryFromImageID(imageID string) string {
	var registryAndRepo string
	atParts := strings.Split(imageID, "@")

	if len(atParts) > 1 {
		registryAndRepo = atParts[0]
	} else {
		colonParts := strings.Split(imageID, ":")
		if len(colonParts) == 0 {
			return ""
		}
		registryAndRepo = colonParts[0]
	}

	registryAndRepoParts := strings.SplitN(registryAndRepo, "/", 2)

	if len(registryAndRepoParts) < 2 {
		return ""
	}
	if strings.Contains(registryAndRepoParts[0], ".") {
		return registryAndRepoParts[0]
	}

	return ""
}

func getRepositoryFromImageID(imageID string) string {
	var registryAndRepo string
	atParts := strings.Split(imageID, "@")

	if len(atParts) > 1 {
		registryAndRepo = atParts[0]
	} else {
		colonParts := strings.Split(imageID, ":")
		if len(colonParts) == 0 {
			return ""
		}
		registryAndRepo = colonParts[0]
	}
	if !strings.Contains(registryAndRepo, ".") {
		return registryAndRepo
	}

	registryAndRepoParts := strings.SplitN(registryAndRepo, "/", 2)

	if len(registryAndRepoParts) < 2 {
		return ""
	}

	return registryAndRepoParts[1]
}

func getImageTagFromImageID(imageID string) string {
	parts := strings.Split(imageID, ":")
	isHash := strings.Contains(imageID, "@")
	if isHash {
		parts = strings.Split(imageID, "@")
	}

	// The tag is expected to be after the last colon
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}

	return "" // return an empty string if imageID doesn't contain a tag
}
