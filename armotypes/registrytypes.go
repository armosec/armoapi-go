package armotypes

type RegistryJobParams struct {
	Name            string `json:"name,omitempty"`
	ID              string `json:"id,omitempty"`
	ClusterName     string `json:"clusterName"`
	RegistryName    string `json:"registryName"`
	CronTabSchedule string `json:"cronTabSchedule,omitempty"`
	JobID           string `json:"jobID,omitempty"`
}

type RegistryInfo struct {
	RegistryName     string     `json:"registryName"`
	RegistryProvider string     `json:"registryProvider"`
	RegistryToken    string     `json:"registryToken"`
	Depth            int        `json:"depth"`
	Include          []string   `json:"include,omitempty"`
	Exclude          []string   `json:"exclude,omitempty"`
	Kind             string     `json:"kind,omitempty"`
	IsHTTPs          bool       `json:"isHTTPs"`
	SkipTLS          bool       `json:"skipTLS"`
	AuthMethod       AuthMethod `json:"authMethod"`
}

type AuthMethod struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}
