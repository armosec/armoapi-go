package notifications

import "github.com/armosec/gojay"

// UnmarshalJSONObject --
func (ert *TopCtrlItem) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {

	switch key {
	case "id":
		return dec.String(&ert.ControlID)
	case "guid":
		return dec.String(&ert.ControlGUID)
	case "name":
		return dec.String(&ert.Name)
	case "remediation":
		return dec.String(&ert.Remediation)
	case "description":
		return dec.String(&ert.Description)
	case "clustersCount":
		return dec.Int64(&ert.ClustersCount)
	case "severityOverall":
		return dec.Int64(&ert.SeverityOverall)
	case "baseScore":
		return dec.Int64(&ert.BaseScore)
	}
	return nil
}

// NKeys --
func (ert *TopCtrlItem) NKeys() int {
	return 0
}
