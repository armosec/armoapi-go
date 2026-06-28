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
	// RiskFactorNone is a filter sentinel representing "no risk factors".
	// It is NOT a mappable risk factor and has no entry in RiskFactorMapping.
	// Do not pass it to GetControlIDsByRiskFactors (it will return an empty slice).
	RiskFactorNone RiskFactor = "None"

	//scan types
	ClusterPosture           ScanType = "cluster"
	RepositoryPosture        ScanType = "repository"
	ContainerVulnerabilities ScanType = "container"
	RegistryVulnerabilities  ScanType = "registry"

	// Agentic entity-type values — the SINGLE source of truth for the AI-Sandbox
	// agentic classification, shared by the inventory/dashboard and the
	// AI-Sandbox serving so both classify identical inputs identically (no
	// duplicated deriveEntityType). The exact strings MUST stay in sync with the
	// values served by postgres-connector services/aisandbox/view.go.
	//
	// EntityTypeMCPServer marks a subject that exposes/serves AI capability
	// (ai_server_providers present) — it acts as an MCP server / model host.
	EntityTypeMCPServer = "MCP Server"
	// EntityTypeAIAgent marks a subject that consumes AI (ai_client_providers
	// present, no server providers) — an AI agent / client workload.
	EntityTypeAIAgent = "AI Agent"
)

// AgenticEntityType is the SINGLE agentic classification rule shared by the
// inventory/dashboard and the AI-Sandbox serving. Precedence is server-wins: a
// subject that serves AI (serverProviders non-empty) is an MCP Server;
// otherwise one that consumes AI (clientProviders non-empty) is an AI Agent;
// neither ⇒ "" (not agentic / unknown).
func AgenticEntityType(clientProviders, serverProviders []string) string {
	if len(serverProviders) > 0 {
		return EntityTypeMCPServer
	}
	if len(clientProviders) > 0 {
		return EntityTypeAIAgent
	}
	return ""
}

// IsAgentic is the binary (yes/no) agentic verdict used by the inventory badge.
// A subject is agentic when it has any AI client OR any AI server provider.
func IsAgentic(clientProviders, serverProviders []string) bool {
	return len(clientProviders) > 0 || len(serverProviders) > 0
}

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
