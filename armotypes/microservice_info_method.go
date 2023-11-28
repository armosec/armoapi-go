package armotypes

import (
	"fmt"
	"strings"
)

// GetContainerImageDetails extract the docker image details of specific container in list
func (contsImages *ContainersStatusData) GetContainerImageDetails(contName string) (string, string, error) {
	imageID := ""
	imageHash := ""
	if contImagesData, ok := (*contsImages)[contName]; ok {
		if contImageTag, ok := contImagesData["image"]; ok {
			imageID = contImageTag
		} else {
			return "", "", fmt.Errorf("failed to find imageID of container '%s' in contsImages. Full data: %v", contName, contsImages)
		}
		if contImageHash, ok := contImagesData["imageID"]; ok {
			lIdx := strings.LastIndex(contImageHash, ":")
			if lIdx < 0 {
				return "", "", fmt.Errorf("failed to find ':' in imageHash of container '%s' in contsImages. Full data: %v", contName, contsImages)
			}
			imageHash = contImageHash[lIdx+1:]
		} else {
			return "", "", fmt.Errorf("failed to find imageHash of container '%s' in contsImages. Full data: %v", contName, contsImages)
		}
	} else {
		return "", "", fmt.Errorf("failed to find images of container '%s' in contsImages. Full data: %v", contName, contsImages)
	}
	return imageID, imageHash, nil
}

// GetShortName returns the last 2 parts of the microservice
func (msi *MicroserviceInfo) GetShortName() string {
	return fmt.Sprintf("%s/%s", msi.Ancestor.Kind, msi.Ancestor.Name)
}
