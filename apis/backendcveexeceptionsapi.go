package apis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/armosec/armoapi-go/armotypes"
	httputils "github.com/armosec/utils-go/httputils"
)

func getCVEExceptionsURL(backendURL string, cusGUID string, designators *armotypes.PortalDesignator) (*url.URL, error) {
	expURL, err := url.Parse(backendURL)
	if err != nil {
		return nil, err
	}
	expURL.Scheme = "https"
	expURL.Path = path.Join(expURL.Path, "v1/armoVulnerabilityExceptions")
	qValues := expURL.Query()
	for k, v := range designators.Attributes {
		qValues.Add(k, v)
	}
	expURL.RawQuery = qValues.Encode()
	return expURL, nil
}

func getCVEExceptionByDEsignator(backendURL string, cusGUID string, designators *armotypes.PortalDesignator) ([]armotypes.VulnerabilityExceptionPolicy, error) {

	var vulnerabilityExceptionPolicy []armotypes.VulnerabilityExceptionPolicy

	url, err := getCVEExceptionsURL(backendURL, cusGUID, designators)
	if err != nil {
		return nil, err
	}

	resp, err := httputils.HttpGet(http.DefaultClient, url.String(), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("getCVEExceptionByDEsignator: resp.StatusCode %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &vulnerabilityExceptionPolicy)
	if err != nil {
		return nil, err
	}

	return vulnerabilityExceptionPolicy, nil
}

func BackendGetCVEExceptionByDEsignator(baseURL string, cusGUID string, designators *armotypes.PortalDesignator) ([]armotypes.VulnerabilityExceptionPolicy, error) {
	vulnerabilityExceptionPolicyList, err := getCVEExceptionByDEsignator(baseURL, cusGUID, designators)
	if err != nil {
		return nil, err
	}
	return vulnerabilityExceptionPolicyList, nil
}
