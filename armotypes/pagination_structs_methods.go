package armotypes

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// FixOrAddAsInnerFilters adds the query scope params as inner filters of the request to adapt the right field names
func (qsp *QueryScopeParams) FixOrAddAsInnerFilters(paginationReq *V2ListRequest, instanceIDField, clusterNameField, namespaceField,
	wlidField, kindField, nameField string) {
	filter := Filters{
		InstanceIDField:  instanceIDField,
		ClusterNameField: clusterNameField,
		NamespaceField:   namespaceField,
		WlidField:        wlidField,
		KindField:        kindField,
		NameField:        nameField,
	}
	qsp.FixOrAddAsInnerFiltersMap(paginationReq, filter)
}

func (qsp *QueryScopeParams) FixOrAddAsInnerFiltersMap(paginationReq *V2ListRequest, filters Filters) {
	if len(paginationReq.InnerFilters) == 0 {
		paginationReq.InnerFilters = []map[string]string{{}}
	}
	for filterIdx := range paginationReq.InnerFilters {
		if qsp.InstanceID != "" && filters.InstanceIDField != "" {
			paginationReq.InnerFilters[filterIdx][filters.InstanceIDField] = qsp.InstanceID
		}
		if filters.ClusterNameField != "" && len(qsp.Cluster) > 0 {
			paginationReq.InnerFilters[filterIdx][filters.ClusterNameField] = strings.Join(qsp.Cluster, ",")
		}
		if filters.NamespaceField != "" && len(qsp.Namespace) > 0 {
			paginationReq.InnerFilters[filterIdx][filters.NamespaceField] = strings.Join(qsp.Namespace, ",")
		}
		if filters.WlidField != "" && len(qsp.WLIDs) > 0 {
			paginationReq.InnerFilters[filterIdx][filters.WlidField] = strings.Join(qsp.WLIDs, ",")
		}
		if filters.KindField != "" && len(qsp.Kind) > 0 {
			paginationReq.InnerFilters[filterIdx][filters.KindField] = strings.Join(qsp.Kind, ",")
		}
		if filters.NameField != "" && len(qsp.Name) > 0 {
			paginationReq.InnerFilters[filterIdx][filters.NameField] = strings.Join(qsp.Name, ",")
		}
		if filters.RepositoryField != "" && len(qsp.Repository) > 0 {
			paginationReq.InnerFilters[filterIdx][filters.RepositoryField] = strings.Join(qsp.Repository, ",")
		}
		if filters.RegistryField != "" && len(qsp.Registry) > 0 {
			paginationReq.InnerFilters[filterIdx][filters.RegistryField] = strings.Join(qsp.Registry, ",")
		}
		if filters.TagField != "" && len(qsp.Tag) > 0 {
			paginationReq.InnerFilters[filterIdx][filters.TagField] = strings.Join(qsp.Tag, ",")
		}
	}
}

// ValidatePageProperties validate page size and page number to be valid
func (u *UniqueValuesRequestV2) ValidatePageProperties(maxPageSize int) {
	if maxPageSize < 1 {
		return
	}
	if u.PageSize > maxPageSize || u.PageSize <= 0 {
		u.PageSize = maxPageSize
	}
}

func (d *Duration) SetDuration(duration time.Duration) {
	*d = Duration(duration)
}

func (d Duration) String() string {
	dur := time.Duration(d).String()
	// If the duration ends with 0s, remove the 0s. It causes an elastic error
	Idx := strings.Index(dur, "m0s")
	if Idx != -1 {
		return dur[:Idx+1]
	}
	return dur
}

func (d Duration) IsValid() bool {
	return d > 0 && d < defaultDuration // elasticsearch documentation: Control how long to keep the search context alive Default: 5m
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

func zeroIntPtr() *int {
	zeroPtr := 0
	return &zeroPtr
}

// ValidatePageProperties validate page size and page number to be valid
func (lr *V2ListRequest) ValidatePageProperties(maxPageSize int) {
	// we are counting from 0 while the user counts from 1... so:
	if !lr.fixedPageNum {
		lr.fixedPageNum = true
		if lr.PageNum == nil {
			lr.PageNum = zeroIntPtr()
		} else {
			if *lr.PageNum > 0 {
				*lr.PageNum--
			} else {
				*lr.PageNum = 0
			}
		}
	}
	if maxPageSize < 1 {
		return
	}
	if lr.PageSize == nil || *lr.PageSize > maxPageSize || *lr.PageSize < 0 {
		lr.PageSize = &maxPageSize
	}
}

// GetFieldsNames retunrs slice of Fields names
func (lr *V2ListRequest) GetFieldsNames() []string {
	res := make([]string, 0, len(lr.InnerFilters))
	for filterObjIdx := range lr.InnerFilters {
		for fieldName := range lr.InnerFilters[filterObjIdx] {
			res = append(res, fieldName)
		}
	}
	if lr.OrderBy != "" {
		sortFields := strings.Split(lr.OrderBy, ",")
		for fieldIdx := range sortFields {
			fieldSlice := strings.Split(sortFields[fieldIdx], ":")
			if len(fieldSlice) > 0 {
				res = append(res, fieldSlice[0])
			}
		}
	}
	return res
}
