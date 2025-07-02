package common

import "strconv"

type ProcessEntity struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	CommandLine string `json:"commandLine,omitempty" bson:"commandLine,omitempty"`
}

type FileEntity struct {
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	Directory string `json:"directory,omitempty" bson:"directory,omitempty"`
}

type DnsEntity struct {
	Domain string `json:"domain,omitempty" bson:"domain,omitempty"`
}

type NetworkEntity struct {
	DstIP    string `json:"dstIP,omitempty" bson:"dstIP,omitempty"`
	DstPort  int    `json:"dstPort,omitempty" bson:"dstPort,omitempty"`
	Protocol string `json:"protocol,omitempty" bson:"protocol,omitempty"`
}

type HttpEntity struct {
	Method    string `json:"method,omitempty" bson:"method,omitempty"`
	Domain    string `json:"domain,omitempty" bson:"domain,omitempty"`
	UserAgent string `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
	Endpoint  string `json:"endpoint,omitempty" bson:"endpoint,omitempty"`
	Payload   string `json:"payload,omitempty" bson:"payload,omitempty"`
}

type CloudAPIEntity struct {
	Service  string `json:"service,omitempty" bson:"service,omitempty"`
	APICall  string `json:"apiCall,omitempty" bson:"apiCall,omitempty"`
	Resource string `json:"resource,omitempty" bson:"resource,omitempty"`
	User     string `json:"user,omitempty" bson:"user,omitempty"`
}

type Identifiers struct {
	Process  *ProcessEntity  `json:"process,omitempty" bson:"process,omitempty"`
	File     *FileEntity     `json:"file,omitempty" bson:"file,omitempty"`
	Dns      *DnsEntity      `json:"dns,omitempty" bson:"dns,omitempty"`
	Network  *NetworkEntity  `json:"network,omitempty" bson:"network,omitempty"`
	Http     *HttpEntity     `json:"http,omitempty" bson:"http,omitempty"`
	CloudAPI *CloudAPIEntity `json:"cloud,omitempty" bson:"cloud,omitempty"`
}

func (identifiers *Identifiers) Flatten() map[string]string {
	identifiers_map := make(map[string]string)

	if identifiers.Process != nil {
		if identifiers.Process.Name != "" {
			identifiers_map["process.name"] = identifiers.Process.Name
		}
		if identifiers.Process.CommandLine != "" {
			identifiers_map["process.commandLine"] = identifiers.Process.CommandLine
		}
	}
	if identifiers.File != nil {
		if identifiers.File.Name != "" {
			identifiers_map["file.name"] = identifiers.File.Name
		}
		if identifiers.File.Directory != "" {
			identifiers_map["file.directory"] = identifiers.File.Directory
		}
	}
	if identifiers.Dns != nil {
		if identifiers.Dns.Domain != "" {
			identifiers_map["dns.domain"] = identifiers.Dns.Domain
		}
	}
	if identifiers.Network != nil {
		if identifiers.Network.DstIP != "" {
			identifiers_map["network.dstIP"] = identifiers.Network.DstIP
		}
		if identifiers.Network.DstPort != 0 {
			identifiers_map["network.dstPort"] = strconv.Itoa(identifiers.Network.DstPort)
		}
		if identifiers.Network.Protocol != "" {
			identifiers_map["network.protocol"] = identifiers.Network.Protocol
		}
	}
	if identifiers.Http != nil {
		if identifiers.Http.Method != "" {
			identifiers_map["http.method"] = identifiers.Http.Method
		}
		if identifiers.Http.Domain != "" {
			identifiers_map["http.domain"] = identifiers.Http.Domain
		}
		if identifiers.Http.UserAgent != "" {
			identifiers_map["http.userAgent"] = identifiers.Http.UserAgent
		}
		if identifiers.Http.Endpoint != "" {
			identifiers_map["http.endpoint"] = identifiers.Http.Endpoint
		}
		if identifiers.Http.Payload != "" {
			identifiers_map["http.payload"] = identifiers.Http.Payload
		}
	}
	if identifiers.CloudAPI != nil {
		if identifiers.CloudAPI.Service != "" {
			identifiers_map["cloud.service"] = identifiers.CloudAPI.Service
		}
		if identifiers.CloudAPI.APICall != "" {
			identifiers_map["cloud.apiCall"] = identifiers.CloudAPI.APICall
		}
		if identifiers.CloudAPI.Resource != "" {
			identifiers_map["cloud.resource"] = identifiers.CloudAPI.Resource
		}
		if identifiers.CloudAPI.User != "" {
			identifiers_map["cloud.user"] = identifiers.CloudAPI.User
		}
	}
	return identifiers_map
}
