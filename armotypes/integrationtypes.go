package armotypes

import "time"

type TicketManager string

const (
	TicketManagerJira TicketManager = "jira"
)

type Ticket struct {
	GUID           string              `json:"guid,omitempty"`           //ticket guid in armo
	JiraCollabGUID string              `json:"jiraCollabGUID,omitempty"` //integration guid between jira creator and the ticket
	TicketManager  TicketManager       `json:"ticketManager"`            //ticket service provider
	Owner          map[string]string   `json:"owner,omitempty"`          //armo entity that owns the ticket
	Subjects       []map[string]string `json:"subjects,omitempty"`       //armo entities mentioned in the ticket
	Link           string              `json:"link,omitempty"`           //link to the ticket
	Status         string              `json:"status,omitempty"`         //status of the ticket
	LinkTitle      string              `json:"linkTitle,omitempty"`      //title of the ticket
	Severity       string              `json:"severity,omitempty"`       //severity of the ticket
	Error          string              `json:"error,omitempty"`          //error message if any
	ErrorCode      int                 `json:"errorCode,omitempty"`      //error code if any (e.g. http status code like 401)
	ProviderData   map[string]string   `json:"providerData,omitempty"`   //provider specific data
	CreatedBy      string              `json:"createdBy,omitempty"`      //user that created the ticket

	// metadata for the ticket
	CustomerGUID string     `json:"customerGUID,omitempty"`
	IssueID      string     `json:"issueID,omitempty"`
	SiteID       string     `json:"siteID,omitempty"`
	Timestamp    *time.Time `json:"timestamp,omitempty"`
	ProjectID    string     `json:"projectID,omitempty"`
	IssueTypeID  string     `json:"issueTypeID,omitempty"`
	Provider     string     `json:"provider,omitempty"`
}
