package cdr

import "time"

type CloudTrailEvent struct {
	EventVersion        string                 `json:"eventVersion"`
	UserIdentity        UserIdentity           `json:"userIdentity"`
	EventTime           time.Time              `json:"eventTime"`
	EventSource         string                 `json:"eventSource"`
	EventName           string                 `json:"eventName"`
	AWSRegion           string                 `json:"awsRegion"`
	SourceIPAddress     string                 `json:"sourceIPAddress"`
	UserAgent           string                 `json:"userAgent"`
	ErrorCode           string                 `json:"errorCode,omitempty"`
	ErrorMessage        string                 `json:"errorMessage,omitempty"`
	RequestParameters   map[string]interface{} `json:"requestParameters,omitempty"`
	ResponseElements    map[string]interface{} `json:"responseElements,omitempty"`
	AdditionalEventData map[string]interface{} `json:"additionalEventData,omitempty"`
	RequestID           string                 `json:"requestId"`
	EventID             string                 `json:"eventId"`
	EventType           string                 `json:"eventType"`
	APIVersion          string                 `json:"apiVersion,omitempty"`
	ReadOnly            bool                   `json:"readOnly"`
	ManagementEvent     bool                   `json:"managementEvent"`
	Resources           []Resource             `json:"resources,omitempty"`
	RecipientAccountId  string                 `json:"recipientAccountId,omitempty"`
	SharedEventID       string                 `json:"sharedEventId,omitempty"`
	VpcEndpointId       string                 `json:"vpcEndpointId,omitempty"`
	TLSDetails          *TLSDetails            `json:"tlsDetails,omitempty"`
	ServiceEventDetails map[string]interface{} `json:"serviceEventDetails,omitempty"`
}

type UserIdentity struct {
	Type           string          `json:"type"`
	PrincipalID    string          `json:"principalId"`
	ARN            string          `json:"arn,omitempty"`
	AccountID      string          `json:"accountId"`
	OrgID          string          `json:"orgId,omitempty"`
	AccessKeyID    string          `json:"accessKeyId,omitempty"`
	UserName       string          `json:"userName,omitempty"`
	InvokedBy      string          `json:"invokedBy,omitempty"`
	SessionContext *SessionContext `json:"sessionContext,omitempty"`
	OnBehalfOf     *OnBehalfOf     `json:"onBehalfOf,omitempty"`
	CredentialId   string          `json:"credentialId,omitempty"`
}

type OnBehalfOf struct {
	UserId           string `json:"userId"`
	IdentityStoreArn string `json:"identityStoreArn"`
}

type SessionContext struct {
	SessionIssuer *SessionIssuer `json:"sessionIssuer,omitempty"`
	Attributes    *Attributes    `json:"attributes,omitempty"`
}

type SessionIssuer struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
	ARN         string `json:"arn"`
	AccountID   string `json:"accountId"`
	UserName    string `json:"userName"`
}

type Attributes struct {
	MfaAuthenticated string `json:"mfaAuthenticated,omitempty"`
	CreationDate     string `json:"creationDate,omitempty"`
}

type Resource struct {
	ResourceType string `json:"resourceType"`
	ResourceName string `json:"resourceName,omitempty"`
	ResourceARN  string `json:"ARN,omitempty"`
}

type TLSDetails struct {
	TLSVersion               string `json:"tlsVersion,omitempty"`
	CipherSuite              string `json:"cipherSuite,omitempty"`
	ClientProvidedHostHeader string `json:"clientProvidedHostHeader,omitempty"`
}
