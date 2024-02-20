package armotypes

import (
	"encoding/json"
	"errors"
)

type SecurityIssueStatus string
type RiskType string
type SecurityIssueSeverity string

const (
	SecurityIssueStatusDetected SecurityIssueStatus = "Detected"
	SecurityIssueStatusResolved SecurityIssueStatus = "Resolved"

	RiskTypeControl    RiskType = "Control"
	RiskTypeAttackPath RiskType = "AttackPath"

	SecurityIssueSeverityCritical SecurityIssueSeverity = "Critical"
	SecurityIssueSeverityHigh     SecurityIssueSeverity = "High"
	SecurityIssueSeverityMedium   SecurityIssueSeverity = "Medium"
	SecurityIssueSeverityLow      SecurityIssueSeverity = "Low"
)

// Risk represents an individual risk with an ID and type
type Risk struct {
	ID   string   `json:"ID"`
	Type RiskType `json:"type"`
}

// UnmarshalJSON is a custom unmarshaler for RiskType that validates its value
func (rt *RiskType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case string(RiskTypeControl), string(RiskTypeAttackPath):
		*rt = RiskType(s)
		return nil
	default:
		return errors.New("invalid RiskType value")
	}
}

// SecurityRisk represents the main object with various fields and an array of Risks
type SecurityRisk struct {
	ID             string           `json:"ID"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	WhatIs         string           `json:"whatIs"`
	Severity       string           `json:"severity"`
	Category       string           `json:"category"`
	Remediation    string           `json:"remediation"`
	Risks          []Risk           `json:"risks"`
	SecurityIssues []ISecurityIssue `json:"securityIssues,omitempty"`
}

func (sr *SecurityRisk) GetRisks() []Risk {
	return sr.Risks
}

func (sr *SecurityRisk) GetRisksIDsByType(riskType RiskType) []string {
	var risksIDs []string
	for _, risk := range sr.Risks {
		if risk.Type == riskType {
			risksIDs = append(risksIDs, risk.ID)
		}
	}
	return risksIDs
}

func (sr *SecurityRisk) GetRiskTypes() []RiskType {

	riskTypes := make(map[RiskType]interface{})

	for _, risk := range sr.Risks {
		riskTypes[risk.Type] = nil
	}

	keys := make([]RiskType, 0, len(riskTypes))
	for key := range riskTypes {
		keys = append(keys, key)
	}

	return keys

}

type SecurityIssuesSummary struct {
	CustomerGUID            string `json:"customerGUID"`
	SecurityRiskID          string `json:"securityRiskID"`
	SecurityRiskName        string `json:"securityRiskName"`
	Category                string `json:"category"`
	Severity                string `json:"severity"`
	LastUpdated             string `json:"lastUpdated"`
	AffectedClustersCount   int    `json:"affectedClustersCount"`
	AffectedNamespacesCount int    `json:"affectedNamespacesCount"`
	AffectedResourcesCount  int    `json:"affectedResourcesCount"`
	AffectedResourcesChange int    `json:"affectedResourcesChange"`
}

type SecurityIssuesCategories struct {
	CategoryResourceCounters map[string]int `json:"categoryResourceCounter"`
	TotalResources           int            `json:"totalResources"`
}

func (sic *SecurityIssuesCategories) SetCategoryTotal(category string, total int) {
	sic.CategoryResourceCounters[category] = total
}

type SecurityIssuesSeverities struct {
	SeverityResourceCounters map[SecurityIssueSeverity]int `json:"severityResourceCounter"`
	TotalResources           int                           `json:"totalResources"`
}

func NewSecurityIssuesCategories() SecurityIssuesCategories {
	return SecurityIssuesCategories{
		CategoryResourceCounters: map[string]int{},
		TotalResources:           0,
	}

}

func (sis *SecurityIssuesSeverities) SetSeverityTotal(severity SecurityIssueSeverity, total int) {
	sis.SeverityResourceCounters[severity] = total
}

func NewSecurityIssuesSeverities() SecurityIssuesSeverities {
	return SecurityIssuesSeverities{
		SeverityResourceCounters: map[SecurityIssueSeverity]int{
			SecurityIssueSeverityCritical: 0,
			SecurityIssueSeverityHigh:     0,
			SecurityIssueSeverityMedium:   0,
			SecurityIssueSeverityLow:      0,
		},
		TotalResources: 0,
	}

}

type ISecurityIssue interface {
}

type SecurityIssue struct {
	ISecurityIssue  `json:",inline,omitempty"`
	Cluster         string   `json:"cluster"`
	Namespace       string   `json:"namespace"`
	ResourceName    string   `json:"resourceName"`
	ResourceID      string   `json:"resourceID"`
	K8sResourceHash string   `json:"k8sResourceHash"`
	RiskID          string   `json:"riskID"` // controlID/attackTrackID
	RiskType        RiskType `json:"riskType,omitempty"`

	SecurityRiskID string `json:"securityRiskID"`

	Status SecurityIssueStatus `json:"status"`

	IsNew bool `json:"isNew"`

	LastTimeDetected string `json:"lastTimeDetected,omitempty"`
	LastTimeResolved string `json:"lastTimeResolved,omitempty"`
}

type SecurityIssueControl struct {
	SecurityIssue `json:",inline"`
	ControlID     string `json:"controlID"`
	ReportGUID    string `json:"reportGUID"`
	FrameworkName string `json:"frameworkName"`
}

type SecurityIssueAttackPath struct {
	SecurityIssue `json:",inline"`
	AttackChainID string `json:"attackChainID"`
}
