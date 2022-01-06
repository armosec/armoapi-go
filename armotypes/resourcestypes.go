package armotypes

import (
	"github.com/armosec/armo-interfaces/interfaces"
)

type Resource struct {
	ResourceID string               `json:"resourceID"`
	Object     interface{}          `json:"object"`
	IMetadata  interfaces.IMetadata `json:"-"`
}
