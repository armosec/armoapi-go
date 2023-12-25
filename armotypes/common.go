package armotypes

// swagger:strfmt uuid4
// Example: 0f42fbe3-d81e-444d-8cc7-bc892c7623e9
type GUID string

type RiskFactor string

const (
	RiskFactorInternetFacing RiskFactor = "Internet facing"
	RiskFactorPrivileged     RiskFactor = "Privileged"
	RiskFactorSecretAccess   RiskFactor = "Secret access"
	RiskFactorDataAccess     RiskFactor = "Data access"
	RiskFactorHostAccess     RiskFactor = "Host access"
)

var RiskFactorMapping = map[string]RiskFactor{
	"C-0256": RiskFactorInternetFacing,
	"C-0046": RiskFactorPrivileged,
	"C-0057": RiskFactorPrivileged,
	"C-0255": RiskFactorSecretAccess,
	"C-0257": RiskFactorDataAccess,
	"C-0038": RiskFactorHostAccess,
	"C-0041": RiskFactorHostAccess,
	"C-0044": RiskFactorHostAccess,
	"C-0048": RiskFactorHostAccess,
}

// todo : verify where to add it

// GetRiskFactors returns a list of unique risk factors for given control IDs.
func GetRiskFactors(controlIDs []string) []RiskFactor {
	riskFactorSet := make(map[RiskFactor]bool)
	for _, id := range controlIDs {
		if riskFactor, exists := RiskFactorMapping[id]; exists {
			riskFactorSet[riskFactor] = true
		}
	}

	var riskFactors []RiskFactor
	for riskFactor := range riskFactorSet {
		riskFactors = append(riskFactors, riskFactor)
	}
	return riskFactors
}
