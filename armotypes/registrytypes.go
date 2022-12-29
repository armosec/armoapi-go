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
	RegistryName     string     `json:"registryName,omitempty" bson:"registryName"`
	RegistryProvider string     `json:"registryProvider,omitempty" bson:"registryProvider"`
	RegistryToken    string     `json:"registryToken,omitempty" bson:"registryToken"`
	Depth            *int       `json:"depth,omitempty" bson:"depth"`
	Include          []string   `json:"include,omitempty" bson:"include"`
	Exclude          []string   `json:"exclude,omitempty" bson:"exclude"`
	Kind             string     `json:"kind,omitempty" bson:"kind"`
	IsHTTPs          *bool      `json:"isHTTPs,omitempty" bson:"isHTTPs"`
	SkipTLSVerify    *bool      `json:"skipTLSVerify,omitempty" bson:"skipTLSVerify"`
	AuthMethod       AuthMethod `json:"authMethod,omitempty" bson:"authMethod"`
	SecretName       string     `json:"secret,omitempty" bson:"secret"`
}

type AuthMethod struct {
	Username string `json:"username,omitempty" bson:"username"`
	Password string `json:"password,omitempty" bson:"password"`
	Type     string `json:"type,omitempty" bson:"type"`
}

type Repository struct {
	RepositoryName string `json:"repositoryName"`
}
