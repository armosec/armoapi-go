package apis

import (
	"encoding/json"

	"github.com/armosec/armoapi-go/armotypes"
)

func getCVEExceptionByDEsignator(backendConn *BackendConnector, cusGUID string, designators *armotypes.PortalDesignator) ([]armotypes.VulnerabilityExceptionPolicy, error) {

	var vulnerabilityExceptionPolicy []armotypes.VulnerabilityExceptionPolicy
	bytes, err := backendConn.HTTPSend("GET", "v1/armoVulnerabilityExceptions", nil, MapQueryWithoutSortKeys, designators.Attributes)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &vulnerabilityExceptionPolicy)
	if err != nil {
		return nil, err
	}

	return vulnerabilityExceptionPolicy, nil
}

func BackendGetCVEExceptionByDEsignator(baseURL string, cusGUID string, designators *armotypes.PortalDesignator) ([]armotypes.VulnerabilityExceptionPolicy, error) {
	backendConn, err := MakePublicBackendConnector(baseURL)
	if err != nil {
		return nil, err
	}

	vulnerabilityExceptionPolicyList, err := getCVEExceptionByDEsignator(backendConn, cusGUID, designators)
	if err != nil {
		return nil, err
	}
	return vulnerabilityExceptionPolicyList, nil
}
