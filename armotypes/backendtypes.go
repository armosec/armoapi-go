package armotypes

import "time"

type V2ListRequest struct {
	// properties of the requested next page
	// Use ValidatePageProperties to set PageSize field
	PageSize *int `json:"pageSize"`
	// One can leave it empty for 0, then call ValidatePageProperties
	PageNum *int `json:"pageNum"`
	// The time window of the list to return. Default: since - begining og the time, until - now.
	Since *time.Time `json:"since"`
	Until *time.Time `json:"until"`
	// Which elements of the list to return, each field can hold multiple values separated by comma
	// An empty map means "return the complete list"
	// Example: [{"severity": "High,Medium",		"type": "61539,30303"}]
	InnerFilters []map[string]string `json:"innerFilters"`
	// How to order (sort) the list, field name + sort order (asc/desc), like https://www.w3schools.com/sql/sql_orderby.asp
	// When empty, the default sort order is used. To disable the default sort order, set IgnoreDefaultSort to true
	// Example: timestamp:asc,severity:desc
	OrderBy string `json:"orderBy"`
	// When true, the default sort order is ignored
	// TODO: take it off, and use the default sort order when OrderBy is empty
	IgnoreDefaultSort bool `json:"ignoreDefaultOrderBy,omitempty"`
	// FieldsList allow us to return only subset of the source document fields
	// Don't expose FieldsList outside without well designed decision
	FieldsList              []string          `json:"includeFields"`
	FieldsReverseKeywordMap map[string]string `json:"-"`
}
