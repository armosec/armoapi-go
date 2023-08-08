package armotypes

import (
	"strings"

	"github.com/armosec/gojay"
)

// ListFields list all UniqueValuesResponseV2 fields
func (uvr *UniqueValuesResponseV2) ListFields(key string) []string {
	fields := []string{}
	for i := range uvr.FieldsCount[key] {
		fields = append(fields, uvr.FieldsCount[key][i].Field)
	}
	return fields
}

// ReplaceFieldsFromKeywords restores the original fields names from the .keyword if necessary
func (uvr *UniqueValuesResponseV2) ReplaceFieldsFromKeywords(keywordMap map[string]string) {
	for fullFieldName, fieldName := range keywordMap {
		if _, ok := uvr.Fields[fullFieldName]; fieldName != fullFieldName && ok {
			uvr.Fields[fieldName] = uvr.Fields[fullFieldName]
			delete(uvr.Fields, fullFieldName)
		}
	}
	for fullFieldName, fieldName := range keywordMap {
		if _, ok := uvr.FieldsCount[fullFieldName]; fieldName != fullFieldName && ok {
			uvr.FieldsCount[fieldName] = uvr.FieldsCount[fullFieldName]
			delete(uvr.FieldsCount, fullFieldName)
		}
	}
}

// ValidateOrderBy vlidate that the order-by field is well configured to the desired state
func (lr *V2ListRequest) ValidateOrderBy(defaultDescOrder string) {
	if defaultDescOrder == "" {
		return
	}
	if lr.OrderBy == "" {
		lr.OrderBy = defaultDescOrder
	}
	if !strings.HasSuffix(lr.OrderBy, ":desc") && !strings.HasSuffix(lr.OrderBy, ":asc") {
		lr.OrderBy = lr.OrderBy + ":desc"
	}
}

// ReplaceFieldsToKeywords replaces the original fields names to the .keyword if necessary
func (lr *V2ListRequest) ReplaceFieldsToKeywords(keywordMap map[string]string) {
	if lr.FieldsReverseKeywordMap == nil {
		lr.FieldsReverseKeywordMap = make(map[string]string, len(keywordMap))
	}
	for fieldName, fullFieldName := range keywordMap {
		for filterObjIdx := range lr.InnerFilters {
			if fieldName != fullFieldName && lr.InnerFilters[filterObjIdx][fieldName] != "" {
				lr.InnerFilters[filterObjIdx][fullFieldName] = lr.InnerFilters[filterObjIdx][fieldName]
				delete(lr.InnerFilters[filterObjIdx], fieldName)
			}
		}
		lr.FieldsReverseKeywordMap[fullFieldName] = fieldName
	}
	sortFields := strings.Split(lr.OrderBy, ",")
	for fieldIdx := range sortFields {
		fieldSlice := strings.Split(sortFields[fieldIdx], ":")
		if len(fieldSlice) > 0 {
			if fullFieldName, ok := keywordMap[fieldSlice[0]]; ok {
				fieldSlice[0] = fullFieldName
			}
		}
		sortFields[fieldIdx] = strings.Join(fieldSlice, ":")
	}
	lr.OrderBy = strings.Join(sortFields, ",")
}

// GetFieldsNames retunrs slice of Fields names
func (uvr *UniqueValuesRequestV2) GetFieldsNames() []string {
	res := make([]string, 0, len(uvr.Fields))
	for fieldName := range uvr.Fields {
		res = append(res, fieldName)
	}
	return res
}

// ReplaceFieldsToKeywords replaces the original fields names to the .keyword if necessary
func (uvr *UniqueValuesRequestV2) ReplaceFieldsToKeywords(keywordMap map[string]string) {
	if uvr.FieldsReverseKeywordMap == nil {
		uvr.FieldsReverseKeywordMap = make(map[string]string, len(keywordMap))
	}
	for fieldName, fullFieldName := range keywordMap {
		for filterObjIdx := range uvr.InnerFilters {
			if fieldName != fullFieldName && uvr.InnerFilters[filterObjIdx][fieldName] != "" {
				uvr.InnerFilters[filterObjIdx][fullFieldName] = uvr.InnerFilters[filterObjIdx][fieldName]
				delete(uvr.InnerFilters[filterObjIdx], fieldName)
			}
		}
		if fieldName != fullFieldName {
			if _, ok := uvr.Fields[fieldName]; ok {
				uvr.Fields[fullFieldName] = uvr.Fields[fieldName]
				delete(uvr.Fields, fieldName)
			}
		}
		uvr.FieldsReverseKeywordMap[fullFieldName] = fieldName
	}
	for joinedFields := range uvr.Fields {
		fieldsNames := strings.Split(joinedFields, V2ListOperatorSeparator)
		if len(fieldsNames) == 1 {
			continue
		}
		for fieldIdx := range fieldsNames {
			if fullFieldName, ok := keywordMap[fieldsNames[fieldIdx]]; ok {
				fieldsNames[fieldIdx] = fullFieldName
			}
		}
		fullFelidsNames := strings.Join(fieldsNames, V2ListOperatorSeparator)
		if fullFelidsNames != joinedFields {
			uvr.FieldsReverseKeywordMap[fullFelidsNames] = joinedFields
			uvr.Fields[fullFelidsNames] = uvr.Fields[joinedFields]
			delete(uvr.Fields, joinedFields)
		}
	}
}

// ReplaceFieldsFromKeywords restores the original fields names from the .keyword if necessary
func (uvr *UniqueCardinalityResponseV2) ReplaceFieldsFromKeywords(keywordMap map[string]string) {
	for fullFieldName, fieldName := range keywordMap {
		if _, ok := uvr.Fields[fullFieldName]; fieldName != fullFieldName && ok {
			uvr.Fields[fieldName] = uvr.Fields[fullFieldName]
			delete(uvr.Fields, fullFieldName)
		}
	}
}

// fixOrAddAsInnerFilters adds the query scope params as inner filters of the request to adapt the right field names
func (qsp *QueryScopeParams) FixOrAddAsUniqueInnerFilters(reqObj *UniqueValuesRequestV2, instanceIDField, clusterNameField, namespaceField,
	wlidField, kindField, nameField string) {
	if len(reqObj.InnerFilters) == 0 {
		reqObj.InnerFilters = []map[string]string{{}}
	}
	for filterIdx := range reqObj.InnerFilters {
		if qsp.InstanceID != "" && instanceIDField != "" {
			if reqObj.InnerFilters[filterIdx][instanceIDField] != "" {
				reqObj.InnerFilters[filterIdx][instanceIDField] += "," + qsp.InstanceID
			} else {
				reqObj.InnerFilters[filterIdx][instanceIDField] = qsp.InstanceID
			}

		}
		if clusterNameField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, clusterNameField, qsp.Cluster)
		}
		if namespaceField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, namespaceField, qsp.Namespace)
		}
		if wlidField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, wlidField, qsp.WLIDs)
		}
		if kindField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, kindField, qsp.Kind)
		}
		if nameField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, nameField, qsp.Name)
		}
	}
}

func (qsp *QueryScopeParams) FixOrAddAsUniqueInnerFiltersMap(reqObj *UniqueValuesRequestV2, filters Filters) {
	if len(reqObj.InnerFilters) == 0 {
		reqObj.InnerFilters = []map[string]string{{}}
	}

	for filterIdx := range reqObj.InnerFilters {
		if qsp.InstanceID != "" && filters.InstanceIDField != "" {
			if reqObj.InnerFilters[filterIdx][filters.InstanceIDField] != "" {
				reqObj.InnerFilters[filterIdx][filters.InstanceIDField] += "," + qsp.InstanceID
			} else {
				reqObj.InnerFilters[filterIdx][filters.InstanceIDField] = qsp.InstanceID
			}

		}
		if filters.ClusterNameField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, filters.ClusterNameField, qsp.Cluster)
		}
		if filters.NamespaceField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, filters.NamespaceField, qsp.Namespace)
		}
		if filters.WlidField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, filters.WlidField, qsp.WLIDs)
		}
		if filters.KindField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, filters.KindField, qsp.Kind)
		}
		if filters.NameField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, filters.NameField, qsp.Name)
		}
		if filters.RepositoryField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, filters.RepositoryField, qsp.Name)
		}
		if filters.RegistryField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, filters.RegistryField, qsp.Name)
		}
		if filters.TagField != "" {
			insertQueryParam2InnerFilters(reqObj, filterIdx, filters.TagField, qsp.Name)
		}
	}
}

func insertQueryParam2InnerFilters(reqObj *UniqueValuesRequestV2, filterIdx int, field string, sliceQueryParam []string) {
	if reqObj.InnerFilters[filterIdx][field] != "" {
		reqObj.InnerFilters[filterIdx][field] = strings.Join(append(sliceQueryParam, reqObj.InnerFilters[filterIdx][field]), ",")
	} else {
		reqObj.InnerFilters[filterIdx][field] = strings.Join(sliceQueryParam, ",")
	}
}

// UnmarshalJSONObject --
func (ert *RespTotal) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "value":
		return dec.Int(&ert.Value)
	case "relation":
		return dec.String(&ert.Relation)
	}
	return nil
}

// NKeys --
func (ert *RespTotal) NKeys() int {
	return 2
}
