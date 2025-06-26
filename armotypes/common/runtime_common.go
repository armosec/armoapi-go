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
	Method   string `json:"method,omitempty" bson:"method,omitempty"`
	Domain   string `json:"domain,omitempty" bson:"domain,omitempty"`
	Endpoint string `json:"endpoint,omitempty" bson:"endpoint,omitempty"`
	Proto    string `json:"proto,omitempty" bson:"proto,omitempty"`
	Payload  string `json:"payload,omitempty" bson:"payload,omitempty"`
}

type CloudAPIEntity struct {
	Service  string `json:"service,omitempty" bson:"service,omitempty"`
	Action   string `json:"action,omitempty" bson:"action,omitempty"`
	Resource string `json:"resource,omitempty" bson:"resource,omitempty"`
	User     string `json:"user,omitempty" bson:"user,omitempty"`
}

type Identifiers struct {
	Process  *ProcessEntity  `json:"process,omitempty" bson:"process,omitempty"`
	File     *FileEntity     `json:"file,omitempty" bson:"file,omitempty"`
	Network  *NetworkEntity  `json:"network,omitempty" bson:"network,omitempty"`
	Http     *HttpEntity     `json:"http,omitempty" bson:"http,omitempty"`
	CloudAPI *CloudAPIEntity `json:"cloudAPI,omitempty" bson:"cloudAPI,omitempty"`
}

func (identifiers *Identifiers) Flatten() map[string]string {
	identifiers_map := make(map[string]string)

	if identifiers.Process != nil {
		identifiers_map["process.name"] = identifiers.Process.Name
		identifiers_map["process.commandLine"] = identifiers.Process.CommandLine
	}
	if identifiers.File != nil {
		identifiers_map["file.path"] = identifiers.File.Path
		identifiers_map["file.directory"] = identifiers.File.Directory
	}
	if identifiers.Network != nil {
		identifiers_map["network.domain"] = identifiers.Network.Domain
		identifiers_map["network.address"] = identifiers.Network.Address
		identifiers_map["network.dstPort"] = identifiers.Network.DstPort
		identifiers_map["network.protocol"] = identifiers.Network.Protocol
	}
	if identifiers.Http != nil {
		identifiers_map["http.method"] = identifiers.Http.Method
		identifiers_map["http.domain"] = identifiers.Http.Domain
		identifiers_map["http.endpoint"] = identifiers.Http.Endpoint
	}
	if identifiers.CloudAPI != nil {
		identifiers_map["cloudAPI.service"] = identifiers.CloudAPI.Service
		identifiers_map["cloudAPI.action"] = identifiers.CloudAPI.Action
		identifiers_map["cloudAPI.resource"] = identifiers.CloudAPI.Resource
		identifiers_map["cloudAPI.user"] = identifiers.CloudAPI.User
	}
	return identifiers_map
}
