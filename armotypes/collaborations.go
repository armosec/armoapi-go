package armotypes

type CollabAssignee struct {

	//example: can be channelID(slack) "C02HD5MU9G8" and etc.
	AssgineeID string `json:"assigneeID"`

	//example: #abuse(slack)
	AssigneeName string `json:"assigneeName"`

	//put here properties of the assignee, ad
	AdditionalInfo []ArmoContext `json:"additionalInfo"`
}
