package armotypes

type UserApplicationProfileRequest struct {
	IncidentGUID string `json:"incidentGUID"`
	ApplicationProfileIdentifiers 
}


type ApplicationProfileIdentifiers struct {
	WorkloadKind string `json:"workloadKind"`
	WorkloadName string `json:"workloadName"`
	Namespace    string `json:"namespace"`
	Cluster 	string `json:"cluster"`
}