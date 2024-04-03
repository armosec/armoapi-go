package notifications

import (
	"time"

	"github.com/armosec/armoapi-go/armotypes"
)

type ReferenceType string //type of the reference (e.g cve-ticket, slack-message etc)

const (
	//tickets types
	ReferenceTypeClusterControlTicket    ReferenceType = "ticket:cluster:control"
	ReferenceTypeRepositoryControlTicket ReferenceType = "ticket:repository:control"
	ReferenceTypeImageTicket             ReferenceType = "ticket:image"
	ReferenceTypeVulnerabilityTicket     ReferenceType = "ticket:vulnerability"
)

// Referance to external integration (e.g link to jira ticket)
type IntegrationReference struct {
	armotypes.PortalBase `json:",inline" bson:"inline"`
	Provider             ChannelProvider               `json:"provider,omitempty" bson:"provider,omitempty"`             //integration provider (e.g jira, slack, teams)
	ProviderData         interface{}                   `json:"providerData,omitempty" bson:"providerData,omitempty"`     //integration provider data (e.g jira ticket data)
	Type                 ReferenceType                 `json:"type,omitempty" bson:"type,omitempty"`                     //type of the reference (e.g cve-ticket, slack-message etc)
	Owner                *armotypes.EntityIdentifiers  `json:"owner,omitempty" bson:"owner,omitempty"`                   //owner identifiers of this reference (e.g resourceHash, wlid)
	RelatedObjects       []armotypes.EntityIdentifiers `json:"relatedObjects,omitempty" bson:"relatedObjects,omitempty"` //related entities identifiers of this reference (e.g cves, controls)
	CreationTime         time.Time                     `json:"creationTime" bson:"creationTime"`                         //creation time of the reference
}
