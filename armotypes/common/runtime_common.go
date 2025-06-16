package common

type ProcessEntity struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	CommandLine string `json:"commandLine,omitempty" bson:"commandLine,omitempty"`
}

type FileEntity struct {
	Path      string `json:"path,omitempty" bson:"path,omitempty"`
	Directory string `json:"directory,omitempty" bson:"directory,omitempty"`
}

type NetworkEntity struct {
	Domain   string `json:"domain,omitempty" bson:"domain,omitempty"`
	Address  string `json:"address,omitempty" bson:"address,omitempty"`
	DstPort  string `json:"dstPort,omitempty" bson:"dstPort,omitempty"`
	Protocol string `json:"protocol,omitempty" bson:"protocol,omitempty"`
}

type HttpEntity struct {
	Method   string            `json:"method,omitempty" bson:"method,omitempty"`
	Domain   string            `json:"domain,omitempty" bson:"domain,omitempty"`
	Endpoint string            `json:"endpoint,omitempty" bson:"endpoint,omitempty"`
	Header   map[string]string `json:"header,omitempty" bson:"header,omitempty"`
	Proto    string            `json:"proto,omitempty" bson:"proto,omitempty"`
	Payload  string            `json:"payload,omitempty" bson:"payload,omitempty"`
}

type CloudAPIEntity struct {
	Service  string `json:"service,omitempty" bson:"service,omitempty"`
	Action   string `json:"action,omitempty" bson:"action,omitempty"`
	Resource string `json:"resource,omitempty" bson:"resource,omitempty"`
	User     string `json:"user,omitempty" bson:"user,omitempty"`
}

type Identifiers struct {
	File     *FileEntity     `json:"file,omitempty" bson:"file,omitempty"`
	Network  *NetworkEntity  `json:"network,omitempty" bson:"network,omitempty"`
	Http     *HttpEntity     `json:"http,omitempty" bson:"http,omitempty"`
	CloudAPI *CloudAPIEntity `json:"cloudAPI,omitempty" bson:"cloudAPI,omitempty"`
}
