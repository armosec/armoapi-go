package notifications

import "github.com/armosec/armoapi-go/identifiers"

type RuntimeIncidentPushNotification struct {
	NewRuntimeIncident NewRuntimeIncident
}

type NewRuntimeIncident struct {
	CustomerGUID        string                       `json:"customerGUID"`
	IncidentPolicyGUIDs []string                     `json:"incidentPolicyGUIDs"`
	IncidentGUID        string                       `json:"incidentGUID"`
	IncidentName        string                       `json:"incidentName"` // incidentType.Name - threatName
	Severity            string                       `json:"severity"`
	Resource            identifiers.PortalDesignator `json:"resource"` // Pod, Node, Workload, Namespace, Cluster, etc.
	Link                string                       `json:"link"`
	Response            *RuntimeIncidentResponse     `json:"response,omitempty"`
}

type RuntimeIncidentResponse struct {
	Action string `json:"action"`
}
