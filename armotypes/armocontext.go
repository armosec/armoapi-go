package armotypes

// context attributes based structure to get more flexible and searchable options
type ArmoContext struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
	Source    string `json:"source"`
	// open question: should we suppport in "not exists" or "not equal". E.g. "apply this recommendation only for non GCP hosted K8s clusters"
}
