package armotypes

const (
	LowestHelmVersionSupportedRegistryScan = "v1.7.9"
	RegistriInfoArgKey                     = "registryInfo-v1"
)

type RegistryJobParams struct {
	Name            string `json:"name,omitempty"`
	ID              string `json:"id,omitempty"`
	ClusterName     string `json:"clusterName"`
	RegistryName    string `json:"registryName"`
	CronTabSchedule string `json:"cronTabSchedule,omitempty"`
	JobID           string `json:"jobID,omitempty"`
}

type RegistriInfoArg struct {
	RegistryName string `json:"registryName"`
}
