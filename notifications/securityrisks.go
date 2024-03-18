package notifications

import "github.com/armosec/armoapi-go/identifiers"

type SecurityIssuePushNotification struct {
	NewSecurityIssues NewSecurityIssues
}

type NewSecurityIssues []NewSecurityIssue

type NewSecurityIssue struct {
	CustomerGUID         string `json:"customerGUID"`
	SecurityRiskID       string `json:"securityRiskID"`
	SecurityRiskName     string `json:"securityRiskName"`
	SecurityRiskSeverity string `json:"securityRiskSeverity"`
	SecurityRiskCategory string `json:"securityRiskCategory"`
	Designators          []identifiers.PortalDesignator
}
