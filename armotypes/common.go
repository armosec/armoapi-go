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
