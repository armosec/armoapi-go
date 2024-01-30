package armotypes

import (
	"time"

	"github.com/armosec/gojay"
)

type RespTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

const (
	V2ListExistsOperator   string = "exists"
	V2ListEqualOperator    string = "equal"
	V2ListMissingOperator  string = "missing"
	V2ListMatchOperator    string = "match"
	V2ListGreaterOperator  string = "greater"
	V2ListLowerOperator    string = "lower"
	V2ListRegexOperator    string = "regex"
	V2ListLikeOperator     string = "like"
	V2ListRangeOperator    string = "range"
	V2ListIgnoreCaseOption string = "ignorecase"
	V2ListArrayOperator    string = "arraymatch"

	V2ListAscendingSort  string = "asc"
	V2ListDescendingSort string = "desc"

	V2ListValueSeparator    = ","
	V2ListOperatorSeparator = "|"
	V2ListSubQuerySeparator = "&"
	V2ListSortTypeSeparator = ":"
	V2ListEscapeChar        = "\\"
)

type QueryScopeParams struct {
	InstanceID string
	Cluster    []string
	Namespace  []string
	WLIDs      []string
	Kind       []string
	Name       []string
	Repository []string
	Registry   []string
	Tag        []string
	Custom     map[string][]string
}

// payload for querying/filtering a list, key: <fieldname> and value is the string value
type RetrieveObjectsByRequestPayload struct {
	MultipleItems map[string][]string
	SingleItems   map[string]string
	Exists        []string
	MustNot       []map[string]interface{}
	ExcludeFields []string
}

// RawJSONObject holds bytes of JSON object
type RawJSONObject gojay.EmbeddedJSON

// MarshalJSON implements the json.marshaler interface
func (rjo *RawJSONObject) MarshalJSON() ([]byte, error) {
	return *rjo, nil
}

type SearchResponse struct {
	Result []RawJSONObject
	Total  *RespTotal
	Cursor *Cursor
	Sort   *SearchAfterResp
}

type Filters struct {
	InstanceIDField  string
	ClusterNameField string
	NamespaceField   string
	WlidField        string
	KindField        string
	NameField        string
	RegistryField    string
	RepositoryField  string
	TagField         string
}

type RespTotal64 struct {
	Value    uint64 `json:"value"`
	Relation string `json:"relation"`
}

type V2ListResponse V2ListResponseGeneric[interface{}]

// V2ListResponse holds the response of some list request with some metadata
type V2ListResponseGeneric[T any] struct {
	Total    RespTotal `json:"total"`
	Response T         `json:"response"`
	// Cursor for quick access to the next page. Not supported yet
	Cursor string `json:"cursor"`
}

// TODO use armotypes.V2ListRequest
// V2ListRequest descripts what portion of the list the client is requesting
// swagger:model PaginationRequest
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
	// Cursor to the next page of former requset.
	// Cursor cannot be used with another parameters of this struct
	Cursor           *Cursor `json:"cursorV1,omitempty"`
	CursorDepracated string  `json:"cursor"`
	// FieldsList allow us to return only subset of the source document fields
	// Don't expose FieldsList outside without well designed decision
	// swagger:ignore
	FieldsList              []string          `json:"includeFields"`
	FieldsReverseKeywordMap map[string]string `json:"-"`
	// TODO: reuse cursor struct (few line above)
	SearchAfter *SearchAfterResp `json:"searchAfter"`
	// For PUT request, can be used to update only specific fields with specific values
	// map of field name to new value
	FieldsToUpdate map[string]string `json:"fieldsToUpdate"`
	//internal flag to indicate if the request is validated (avoid fixing pagination twice in the same request)
	FixedPageNum bool `json:"_FixedPageNum"`
}

type Cursor struct {
	Id        string    `json:"id,omitempty"`
	KeepAlive *Duration `json:"keepAlive,omitempty"`
}

type SearchAfterResp struct {
	Sort interface{} `json:"sort"`
}

// UniqueValuesRequestV2 holds data to return unique values to
type UniqueValuesRequestV2 struct {
	Fields map[string]string `json:"fields"`
	// Which elements of the list to return, each field can hold multiple values separated by comma
	// Example: ": {"severity": "High,Medium",		"type": "61539,30303"}
	// An empty map means "return the complete list"
	InnerFilters []map[string]string `json:"innerFilters"`
	PageSize     int                 `json:"pageSize,omitempty"`
	//for apis that support pagination
	PageNum                 *int              `json:"pageNum,omitempty"`
	FieldsReverseKeywordMap map[string]string `json:"-"`
	Cursor                  string            `json:"-"`
	// The time window to search (Default: since - beginning of the time, until - now)
	Since          *time.Time `json:"since,omitempty"`
	Until          *time.Time `json:"until,omitempty"`
	TimestampField string     `json:"-"`
}

// UniqueCardinalityResponseV2 holds response data of cardinality request
type UniqueCardinalityResponseV2 struct {
	Fields map[string]uint64 `json:"fields"`
	// UniqueTotalCount map[string]ElasticRespTotal `json:"totalCount"`
}

// UniqueValuesResponseFieldsCount  holds response data of UniqueValuesResponseV2 request
type UniqueValuesResponseFieldsCount struct {
	Field string `json:"key"`
	Count int64  `json:"count"`
}

// UniqueValuesResponseV2 holds response data of unique values
type UniqueValuesResponseV2 struct {
	Fields      map[string][]string                          `json:"fields"`
	FieldsCount map[string][]UniqueValuesResponseFieldsCount `json:"fieldsCount"`
	// UniqueTotalCount map[string]ElasticRespTotal `json:"totalCount"`
}

type Duration time.Duration

var defaultDuration = Duration(5 * time.Minute)
