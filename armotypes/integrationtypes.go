package armotypes

type TicketManager string

const (
	TicketManagerJira TicketManager = "jira"
)

type Ticket struct {
	GUID          string              `json:"guid,omitempty"`     //ticket guid in armo
	TicketManager TicketManager       `json:"ticketManager"`      //ticket service provider
	OwnerID       string              `json:"ownerID,omitempty"`  //armo entity that owns the ticket
	Subjects      []map[string]string `json:"subjects,omitempty"` //armo entities mentioned in the ticket
	Link          string              `json:"link,omitempty"`     //link to the ticket
	Status        string              `json:"status,omitempty"`   //status of the ticket
	Title         string              `json:"title,omitempty"`    //title of the ticket can be id or other identifier according to the ticket manager
	Severity      string              `json:"severity,omitempty"` //severity of the ticket
	Error         string              `json:"error,omitempty"`    //error message if any
}
