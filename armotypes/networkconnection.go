package armotypes

// NetworkConnection network connection
type NetworkConnection struct {
	IPAddress                 string  `json:"ipAddress"`
	Inbound                   bool    `json:"inbound"`
	DNSName                   string  `json:"dnsName"`
	Port                      int32   `json:"port"`
	Protocol                  string  `json:"protocol"`
	EndpointWorkloadName      *string `json:"endpointWorkloadName,omitempty"`
	EndpointWorkloadNamespace *string `json:"endpointWorkloadNamespace,omitempty"`
	EndpointWorkloadKind      *string `json:"endpointWorkloadKind,omitempty"`
}

type AggregatedNetworkConnection struct {
	DNSName     string   `json:"dnsName"`
	Direction   string   `json:"direction"`
	Port        int32    `json:"port"`
	Protocol    string   `json:"protocol"`
	IPAddresses []string `json:"ipAddresses"`
}