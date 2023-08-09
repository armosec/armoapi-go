package armotypes

// PaginationSearchFunc declaring function which returns data ready for pagination
type PaginationSearchFunc func(customerGUID, instacnceID string, wlids []string, paginationObject *V2ListRequest) ([]RawJSONObject, *RespTotal, error)

// PaginationCursorFunc declaring function which returns data ready for pagination by cursor to the next page
type PaginationCursorFunc func(customerGUID, instacnceID string, wlids []string, paginationObject *V2ListRequest) (*V2ListResponse, error)

// PaginationSearchByScopeFiltersScrollFunc declaring function which returns data ready for paginationtype PaginationSearchByScopeFiltersFunc func(customerGUID string, scopeFilters *armotypes.QueryScopeParams, paginationObject *armotypes.V2ListRequest) ([]armotypes.RawJSONObject, *ElasticRespTotal, error)
type PaginationSearchByScopeFiltersScrollFunc func(customerGUID string, scopeFilters *QueryScopeParams, paginationObject *V2ListRequest) (*SearchResponse, error)
type PaginationSearchByScopeFiltersFunc func(customerGUID string, scopeFilters *QueryScopeParams, paginationObject *V2ListRequest) ([]RawJSONObject, *RespTotal, error)
type UniqueValuesSearchByScopeFiltersFunc func(customerGUID string, scopeFilters *QueryScopeParams, reqObj *UniqueValuesRequestV2) (*UniqueValuesResponseV2, error)
type CountFunc func(customerGUID string, scopeFilters *QueryScopeParams, paginationObject *V2ListRequest) (uint64, error)
