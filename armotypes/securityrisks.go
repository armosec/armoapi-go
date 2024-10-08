package armotypes

import (
	"encoding/json"
	"errors"
)

type SecurityIssueStatus string
type RiskType string
type SecurityIssueSeverity string
type ResolvedReason string

const (
	SecurityIssueStatusDetected  SecurityIssueStatus = "Detected"
	SecurityIssueStatusResolved  SecurityIssueStatus = "Resolved"
	SecurityIssueStatusException SecurityIssueStatus = "Exception"

	RiskTypeControl                  RiskType = "Control"
	RiskTypeControlWithNetworkPolicy RiskType = "ControlWithNetworkPolicy"
	RiskTypeAttackPath               RiskType = "AttackPath"
	RiskTypeVulnerability            RiskType = "Vulnerability"

	SecurityIssueSeverityCritical SecurityIssueSeverity = "Critical"
	SecurityIssueSeverityHigh     SecurityIssueSeverity = "High"
	SecurityIssueSeverityMedium   SecurityIssueSeverity = "Medium"
	SecurityIssueSeverityLow      SecurityIssueSeverity = "Low"

	ResolvedReasonResourceDeleted ResolvedReason = "ResourceDeleted"
	ResolvedReasonClusterDeleted  ResolvedReason = "ClusterDeleted"
	ResolvedReasonRiskResolved    ResolvedReason = "RiskResolved"
)

// Risk represents an individual risk with an ID and type
type Risk struct {
	ID   string   `json:"ID"`
	Type RiskType `json:"type"`

	// field to be enriched by the backend, relevant only to type Control
	FixByNetworkPolicy bool `json:"fixByNetworkPolicy"`
}

// UnmarshalJSON is a custom unmarshaler for RiskType that validates its value
func (rt *RiskType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case string(RiskTypeControl), string(RiskTypeAttackPath), string(RiskTypeVulnerability):
		*rt = RiskType(s)
		return nil
	default:
		return errors.New("invalid RiskType value")
	}
}

// SecurityRisk represents the main object with various fields and an array of Risks
type SecurityRisk struct {
	ID               string           `json:"ID"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	WhatIs           string           `json:"whatIs"`
	Severity         string           `json:"severity"`
	Category         string           `json:"category"`
	Remediation      string           `json:"remediation"`
	Risks            []Risk           `json:"risks"`
	SecurityIssues   []ISecurityIssue `json:"securityIssues,omitempty"`
	SmartRemediation bool             `json:"smartRemediation"`
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
	SecurityRiskID                   string `json:"securityRiskID"`
	SecurityRiskName                 string `json:"securityRiskName"`
	Category                         string `json:"category"`
	Severity                         string `json:"severity"`
	LastUpdated                      string `json:"lastUpdated"`
	AffectedClustersCount            int    `json:"affectedClustersCount"`
	AffectedNamespacesCount          int    `json:"affectedNamespacesCount"`
	AffectedResourcesCount           int    `json:"affectedResourcesCount"`
	ResourcesDetectedLastUpdateCount int    `json:"resourcesDetectedLastUpdateCount"`
	ResourcesResolvedLastUpdateCount int    `json:"resourcesResolvedLastUpdateCount"`

	ResourcesDetectedLastChangeCount int        `json:"resourcesDetectedLastChangeCount"`
	ResourcesDetectedLastChange      []Resource `json:"resourcesDetectedLastChange"`

	// resources that are resolved excluding deleted
	ResourcesResolvedLastChangeCount int        `json:"resourcesResolvedLastChangeCount"`
	ResourcesResolvedLastChange      []Resource `json:"resourcesResolvedLastChange"`

	// resources that are resolved because of a kubernetes resource deletion or cluster deletion
	ResourcesDeletedLastChangeCount int        `json:"resourcesDeletedLastChangeCount"`
	ResourcesDeletedLastChange      []Resource `json:"resourcesDeletedLastChange"`

	AffectedResourcesChange int `json:"affectedResourcesChange"`

	// if True, control supports smart remediation
	// swagger:ignore
	SupportsSmartRemediation bool `json:"supportsSmartRemediation"` // DEPRECATED
	SmartRemediation         bool `json:"smartRemediation"`

	Tickets []Ticket `json:"tickets,omitempty"`
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
	GetClusterName() string
	GetShortClusterName() string
	SetClusterName(string)
	SetShortClusterName(string)
	GetK8sResourceHash() string
	SetTickets([]Ticket)
}

type SecurityIssue struct {
	ISecurityIssue   `json:",inline,omitempty"`
	Cluster          string   `json:"cluster"`
	ClusterShortName string   `json:"clusterShortName"`
	Namespace        string   `json:"namespace"`
	ResourceName     string   `json:"resourceName"`
	Kind             string   `json:"kind"`
	ResourceID       string   `json:"resourceID"`
	K8sResourceHash  string   `json:"k8sResourceHash"`
	RiskID           string   `json:"riskID"` // controlID/attackTrackID
	RiskType         RiskType `json:"riskType,omitempty"`

	SecurityRiskID string `json:"securityRiskID"`

	Status SecurityIssueStatus `json:"status"`

	IsNew bool `json:"isNew"`

	LastTimeDetected    string `json:"lastTimeDetected,omitempty"`
	LastTimeResolved    string `json:"lastTimeResolved,omitempty"`
	ExceptionApplied    bool   `json:"exceptionApplied"`
	ExceptionPolicyGUID string `json:"exceptionPolicyGUID"`

	Tickets []Ticket `json:"tickets,omitempty"`
}

func (si *SecurityIssue) GetClusterName() string {
	return si.Cluster
}

func (si *SecurityIssue) GetShortClusterName() string {
	return si.ClusterShortName
}

func (si *SecurityIssue) SetClusterName(clusterName string) {
	si.Cluster = clusterName
}

func (si *SecurityIssue) SetShortClusterName(clusterShortName string) {
	si.ClusterShortName = clusterShortName
}

func (si *SecurityIssue) GetK8sResourceHash() string {
	return si.K8sResourceHash
}

func (si *SecurityIssue) SetTickets(tickets []Ticket) {
	si.Tickets = tickets
}

type SecurityIssueControl struct {
	SecurityIssue `json:",inline"`
	ControlID     string `json:"controlID"`
	ReportGUID    string `json:"reportGUID"`
	FrameworkName string `json:"frameworkName"`

	// relevant for controls with network policy fix
	AppliedNetworkPolicyType string              `json:"appliedNetworkPolicyType,omitempty"`
	NetworkPolicyStatus      NetworkPolicyStatus `json:"networkPolicyStatus,omitempty"`

	MissingRuntimeInfoReason MissingRuntimeInfoReason `json:"missingRuntimeInfoReason,omitempty"`
}

type SecurityIssueAttackPath struct {
	SecurityIssue `json:",inline"`
	AttackChainID string `json:"attackChainID"`
	FirstSeen     string `json:"firstSeen"`
}

type SecurityIssueVulnerability struct {
	SecurityIssue    `json:",inline"`
	CriticalCount    int                 `json:"criticalCount"`
	HighCount        int                 `json:"highCount"`
	MediumCount      int                 `json:"mediumCount"`
	LowCount         int                 `json:"lowCount"`
	SeverityStats    map[string][]string `json:"severityStats"`
	RiskFactorsCount int                 `json:"riskFactorsCount"`
	RiskFactors      []RiskFactor        `json:"riskFactors"`
}

type SecurityRiskExceptionPolicy struct {
	BaseExceptionPolicy `json:",inline"`
	Name                string `json:"name"`
	Category            string `json:"category"`
	Severity            string `json:"severity"`
	SecurityRiskID      string `json:"securityRiskID"`
	Risks               []Risk `json:"risks"`
}

type SecurityIssuesTrends struct {

	// date in format yyyy-mm-dd
	Date string `json:"date"`

	// new detected issues within the day
	NewDetected int `json:"newDetected"`

	// new resolved issues within the day
	NewResolved int `json:"newResolved"`

	TotalNewDetectedUpToDate int `json:"totalNewDetectedUpToDate"`

	TotalNewResolvedUpToDate int `json:"totalNewResolvedUpToDate"`

	// new detected issues at the end of the day
	NewDetectedEndOfDay int `json:"newDetectedEndOfDay"`

	// new resolved issues at the end of the day
	NewResolvedEndOfDay int `json:"newResolvedEndOfDay"`

	// total detected from the beginning of the period until current date
	TotalDetectedUpToDate int `json:"totalDetectedUpToDate"`
}

type SecurityIssuesTrendsSummary struct {
	SecurityIssuesTrends []SecurityIssuesTrends `json:"securityIssuesTrends"`

	// total issues detected for the period
	TotalDetectedForPeriod int `json:"totalDetectedForPeriod"`

	// total issues resolved for the period
	TotalResolvedForPeriod int `json:"totalResolvedForPeriod"`

	// current detected issues
	CurrentDetected int `json:"currentDetected"`

	// CurrentDetected - TotalDetectedUpToDate of first date of period.
	ChangeFromBeginningOfPeriod int `json:"changeFromBeginningOfPeriod"`
}
