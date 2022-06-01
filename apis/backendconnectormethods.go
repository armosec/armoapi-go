package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var globalPublicBackendClient *http.Client
var globalPublicBackendClientLock = &sync.Mutex{}

func getPublicBackendClient() *http.Client {
	if globalPublicBackendClient == nil {
		globalPublicBackendClientLock.Lock()
		defer globalPublicBackendClientLock.Unlock()
		if globalPublicBackendClient == nil {
			globalPublicBackendClient = &http.Client{}
		}
	}
	return globalPublicBackendClient
}

func MakePublicBackendConnector(baseURL string) (*BackendConnector, error) {
	publicBackendClient := getPublicBackendClient()
	if err := ValidatePublicBEConnectorMakerInput(publicBackendClient, baseURL); err != nil {
		return nil, err
	}
	conn := &BackendConnector{BaseURL: baseURL, HTTPClient: publicBackendClient}
	return conn, nil
}

func MakeBackendConnector(client *http.Client, baseURL string, loginDetails *CustomerLoginDetails) (*BackendConnector, error) {
	if err := ValidateBEConnectorMakerInput(client, baseURL, loginDetails); err != nil {
		return nil, err
	}
	conn := &BackendConnector{BaseURL: baseURL, Credentials: loginDetails, HTTPClient: client}
	err := conn.Login()

	return conn, err
}

func ValidatePublicBEConnectorMakerInput(client *http.Client, baseURL string) error {
	if client == nil {
		return fmt.Errorf("You must provide an initialized httpclient")
	}
	if len(baseURL) == 0 {
		return fmt.Errorf("you must provide a valid backend url")
	}
	return nil
}

func ValidateBEConnectorMakerInput(client *http.Client, baseURL string, loginDetails *CustomerLoginDetails) error {
	var err error
	if err = ValidatePublicBEConnectorMakerInput(client, baseURL); err != nil {
		return err
	}
	if loginDetails == nil || (len(loginDetails.Email) == 0 && len(loginDetails.Password) == 0) {
		return fmt.Errorf("you must provide valid login details")
	}
	return nil
}

func (r *BackendConnector) Login() error {
	if !r.IsExpired() {
		return nil
	}

	loginInfoBytes, err := json.Marshal(r.Credentials)
	if err != nil {
		return fmt.Errorf("unable to marshal credentials properly")
	}

	beURL := fmt.Sprintf("%v/%v", r.BaseURL, "login")

	req, err := http.NewRequest("POST", beURL, bytes.NewReader(loginInfoBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Referer", strings.Replace(beURL, "dashbe", "cpanel", 1))
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read login response")
	}

	loginS := &BELoginResponse{}
	json.Unmarshal(body, &loginS)

	loginS.Cookies = resp.Cookies()
	r.BELoginResponse = loginS

	return nil
}

func (r *BackendConnector) IsExpired() bool {
	return r.BELoginResponse == nil || r.BELoginResponse.ToLoginObject().IsExpired()
}

func (r *BackendConnector) GetBaseURL() string {
	return r.BaseURL
}
func (r *BackendConnector) GetLoginObj() *LoginObject {
	return r.BELoginResponse.ToLoginObject()
}
func (r *BackendConnector) GetClient() *http.Client {
	return r.HTTPClient
}

func (r *BackendConnector) HTTPSend(httpverb string,
	endpoint string,
	payload []byte,
	f HTTPReqFunc,
	login bool,
	qryData interface{}) ([]byte, error) {

	beURL := fmt.Sprintf("%v/%v", r.GetBaseURL(), endpoint)
	req, err := http.NewRequest(httpverb, beURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if login {
		if r.IsExpired() {
			if err := r.Login(); err != nil {
				return nil, err
			}
		}
		loginobj := r.GetLoginObj()
		req.Header.Set("Authorization", loginobj.Authorization)
		q.Set("customerGUID", loginobj.GUID)
		for _, cookie := range loginobj.Cookies {
			req.AddCookie(cookie)
		}
	}

	req.URL.RawQuery = q.Encode()
	f(req, qryData)
	resp, err := r.GetClient().Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status code: %d, status: %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
