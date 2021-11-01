package armotypes

import "time"

type ApprovementStatus int

const (
	ApprovementStatusApprove ApprovementStatus = iota + 1
	ApprovementStatusDecline
	ApprovementStatusPending
)

type AssociationStatus int

const (
	AssociationStatusAssigned AssociationStatus = iota + 1
	AssociationStatusShown
	AssociationStatusDeclineByUser
	AssociationStatusHandled // the user took this recommendation into account
	AssociationStatusFixed   // the user fixed the issue in some another way
)

type UpdateAuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	UserName  string    `json:"userName"`
}

type ApprovementState struct {
	UpdateAuditEntry `json:",inline"`
	Status           ApprovementStatus `json:"status"`
}

type RecommendationSkeletonV1 struct {
	PortalBase `json:",inline"`
	// audit for manual changes made in this recommendation
	UpdatesAudit []UpdateAuditEntry `json:"updatesAudit"`
	// the action the user should take
	Action      string `json:"action"`
	Description string `json:"description"`
	// link to some well explained description of this recommendation
	DescriptionLink string `json:"descriptionLink"`
	// the context to show this recommendation in
	Context []ArmoContext `json:"context"`
	// the approvement status. Do we should show this recommendation to users?
	Approvement ApprovementStatus `json:"approvement"`
}

// this structure is dedicated to connect between recommendation and specific resource and trace the user actions taken due to this recommendation
type RecommendationAssociation struct {
	PortalBase `json:",inline"`
	// audit for user actions taken for this recommendation
	UpdatesAudit []UpdateAuditEntry `json:"updatesAudit"`
	// the context to show this recommendation to this customer
	Context []ArmoContext `json:"context"`
	// designator object as we have in current resources represntaion
	// this is about to be useless
	Designators PortalDesignator `json:"designators"`
	// guid of the recommendation in recommendation DB
	RecommendationPrototypeGUID string `json:"recommendationPrototypeGUID"`
	// current status of this recommendation for the given resource
	Status AssociationStatus `json:"status"`
}

// {
//     "customerGUID": "",
//     "designators": { // designator object as we have in current resources represntaion.
//         // limitation of this structure is that we need separate "sub-index" for each attribute and we
//         // don't have any extra information regarding this attributes (like source, equlity and so on)
//         "designatorsType": "attributes",
//         "attributes": {
//             "lbl1": "val11"
//         }
//     },
//     "context": [ // more flexible structure for attributes
//         {
//             "attributeName": "",
//             "attributeValue": "",
//             "attributesource": ""
//         }
//     ],
//     "actionPrototypeID": "<guid>", //guid of the recommendation in recommendation DB
//     "status": 2, // current status of this recommendation for the given resource. (e.g. 1- attached ,2- ignored, 3- applied, 4- another remediation, ...)
//     "updatesAudit": [
//         {
//             "time": "",
//             "userName": "",
//             "status": 3
//         }
//     ]
// }
