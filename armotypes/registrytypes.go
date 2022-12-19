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
	RegistryName     string     `json:"registryName,omitempty"`
	RegistryProvider string     `json:"registryProvider,omitempty"`
	RegistryToken    string     `json:"registryToken,omitempty"`
	Depth            int        `json:"depth,omitempty"`
	Include          []string   `json:"include,omitempty"`
	Exclude          []string   `json:"exclude,omitempty"`
	Kind             string     `json:"kind,omitempty"`
	IsHTTPs          bool       `json:"isHTTPs,omitempty"`
	SkipTLSVerify    bool       `json:"skipTLSVerify,omitempty"`
	AuthMethod       AuthMethod `json:"authMethod,omitempty"`
}

type AuthMethod struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Type     string `json:"type,omitempty"`
}
