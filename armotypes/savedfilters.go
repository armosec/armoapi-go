package armotypes

type SavedFilter struct {
	PortalBase `json:",inline" bson:"inline"`
	Subject    string              `json:"subject" bson:"subject"`
	View       string              `json:"view" bson:"view"`
	Filters    []map[string]string `json:"filters" bson:"filters"`
	IsDefault  bool                `json:"isDefault" bson:"isDefault"`
}
