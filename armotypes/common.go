package armotypes

import "strings"

// swagger:strfmt uuid4
// Example: 0f42fbe3-d81e-444d-8cc7-bc892c7623e9
type GUID string

type RiskFactor string

type ScanType string

const (
	//risk factors
	RiskFactorExternalFacing RiskFactor = "External facing"
	RiskFactorPrivileged     RiskFactor = "Privileged"
	RiskFactorSecretAccess   RiskFactor = "Secret access"
	RiskFactorDataAccess     RiskFactor = "Data access"
	RiskFactorHostAccess     RiskFactor = "Host access"
	RiskFactorAILLMClient    RiskFactor = "AI/LLM client"
	RiskFactorAILLMServer    RiskFactor = "AI/LLM service"
	RiskFactorInternetFacing RiskFactor = "Internet facing"

	//scan types
	ClusterPosture           ScanType = "cluster"
	RepositoryPosture        ScanType = "repository"
	ContainerVulnerabilities ScanType = "container"
	RegistryVulnerabilities  ScanType = "registry"
)

var RiskFactorMapping = map[string]RiskFactor{
	"C-0256":        RiskFactorExternalFacing,
	"C-0266":        RiskFactorExternalFacing,
	"C-0046":        RiskFactorPrivileged,
	"C-0057":        RiskFactorPrivileged,
	"C-0255":        RiskFactorSecretAccess,
	"C-0257":        RiskFactorDataAccess,
	"C-0038":        RiskFactorHostAccess,
	"C-0041":        RiskFactorHostAccess,
	"C-0044":        RiskFactorHostAccess,
	"C-0048":        RiskFactorHostAccess,
	"C-AILLMClient": RiskFactorAILLMClient,
	"C-AILLMServer": RiskFactorAILLMServer,
}

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

func GetControlIDsByRiskFactors(riskFactorsStr string) []string {
	riskFactors := strings.Split(riskFactorsStr, ",")
	controlIDSet := make(map[string]bool)

	for _, rfStr := range riskFactors {
		rfStr = strings.TrimSpace(rfStr)
		rf := RiskFactor(rfStr)
		for controlID, mappedRF := range RiskFactorMapping {
			if mappedRF == rf {
				controlIDSet[controlID] = true
			}
		}
	}

	// Convert set to slice
	var controlIDs []string
	for controlID := range controlIDSet {
		controlIDs = append(controlIDs, controlID)
	}

	return controlIDs
}
