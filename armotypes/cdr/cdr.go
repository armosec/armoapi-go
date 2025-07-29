package cdr

import "github.com/armosec/armoapi-go/armotypes/common"

type CustomerDetails struct {
	// CustomerGUID is the unique identifier of the customer account
	CustomerGUID string `json:"customerGUID"`
	// AccessKey is the access key of the customer account
	AccessKey string `json:"accessKey"`
}

// Cloud services
type CloudService string

const (
	// CloudTrail is the cloudtrail service
	CloudTrail CloudService = "cloudtrail"
	// Add more cloud services here
)

// Cloud providers
type CloudProvider string

const (
	// AWS is the AWS cloud provider
	AWS CloudProvider = "aws"
	// Add more cloud providers here
)

type CloudMetadata struct {
	// Provider is the cloud provider
	Provider CloudProvider `json:"provider,omitempty"`
	// SourceService is the source service (e.g cloudtrail, cloudwatch, etc)
	SourceService CloudService `json:"sourceService,omitempty"`
}

// The types corresponds to the SourceService type
type EventData struct {
	// AWSCloudTrail cloudtrail event
	AWSCloudTrail *CloudTrailEvent `json:"awsCloudTrail,omitempty"`
	// Target resource
	TargetResource string `json:"targetResource,omitempty"`
	// Identifiers of the alert
	Identifiers *common.Identifiers `json:"identifiers,omitempty"`
	// Add more cloud event data here
}

type CdrAlert struct {
	// CloudMetadata is the metadata of the cloud
	CloudMetadata `json:"cloudMetadata,omitempty"`
	// EventData is the event data
	EventData `json:"eventData,omitempty"`
	// RuleName is the name of the rule
	RuleName string `json:"ruleName,omitempty"`
	// RuleID is the unique identifier of the rule
	RuleID string `json:"ruleID,omitempty"`
	// Description is the description of the rule
	Description string `json:"description,omitempty"`
	// Priority is the severity of the rule
	Priority string `json:"priority,omitempty"`
	// Tags is the tags of the rule
	Tags []string `json:"tags,omitempty"`
	// Message is the failure message
	Message string `json:"message,omitempty"`
	// MitreTactic is the MITRE ATT&CK tactic
	MitreTactic string `json:"mitreTactic,omitempty"`
	// MitreTechnique is the MITRE ATT&CK technique
	MitreTechnique string `json:"mitreTechnique,omitempty"`
	// UniqueID is the unique identifier of the alert
	UniqueID string `json:"uniqueID,omitempty"`
}

type CdrAlertBatch struct {
	// CustomerGUID is the unique identifier of the customer
	CustomerGUID string `json:"customerGUID,omitempty"`
	// CloudAccountID is the unique identifier of the cloud account
	CloudAccountID string `json:"cloudAccountID,omitempty"`
	// OrgID is the unique identifier of the organization
	OrgID string `json:"orgID,omitempty"`
	// Provider is the cloud provider
	Provider CloudProvider `json:"provider,omitempty"`
	// RuleFailures is the list of rule failures
	RuleFailures []CdrAlert `json:"ruleFailures,omitempty"`
}
