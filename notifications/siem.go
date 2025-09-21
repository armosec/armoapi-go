package notifications

import (
	"fmt"

	"github.com/armosec/armoapi-go/armotypes"
)

type SIEMProvider string

const (
	SIEMProviderSplunk            SIEMProvider = "splunk"
	SIEMProviderSumo              SIEMProvider = "sumoLogic"
	SIEMProviderMicrosoftSentinel SIEMProvider = "microsoftSentinel"
	SIEMProviderWebhook           SIEMProvider = "webhook"
)

type TestMessageStatus string

const (
	TestMessageStatusSuccess TestMessageStatus = "successful"
	TestMessageStatusFailure TestMessageStatus = "failed"
)

type SumoLogicConfig struct {
	HTTPSourceAddress string `json:"httpSourceAddress" bson:"httpSourceAddress"`
}

type SplunkConfig struct {
	URL   string `json:"url" bson:"url"`
	Port  string `json:"port" bson:"port"`
	Token string `json:"token" bson:"token"`
}

type MicrosoftSentinelConfig struct {
	WorkSpaceID string `json:"workSpaceID" bson:"workSpaceID"`
	PrimaryKey  string `json:"primaryKey" bson:"primaryKey"`
}

type SIEMIntegration struct {
	armotypes.PortalBase `json:",inline" bson:"inline"`
	Name                 string                 `json:"name" bson:"name"`
	CustomerGUID         string                 `json:"customerGUID" bson:"customerGUID"`
	Provider             SIEMProvider           `json:"provider" bson:"provider"`
	Configuration        map[string]interface{} `json:"configuration" bson:"configuration"`
	IsEnabled            bool                   `json:"isEnabled" bson:"isEnabled"`
	TestMessageStatus    TestMessageStatus      `json:"testMessageStatus" bson:"testMessageStatus"`
}

type WebhookConfig struct {
	WebhookURL string `json:"webhookURL"`
}

type SumoLogicRequest struct {
	GUID              string            `json:"guid"`
	Name              string            `json:"name"`
	HttpSourceAddress string            `json:"httpSourceAddress"`
	TestMessageStatus TestMessageStatus `json:"testMessageStatus"`
}

func (r *SumoLogicRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.HttpSourceAddress == "" {
		return fmt.Errorf("httpSourceAddress is required")
	}
	return nil
}

type SplunkRequest struct {
	GUID              string            `json:"guid"`
	Name              string            `json:"name"`
	URL               string            `json:"url"`
	Port              string            `json:"port"`
	Token             string            `json:"token"`
	TestMessageStatus TestMessageStatus `json:"testMessageStatus"`
}

func (r *SplunkRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.URL == "" {
		return fmt.Errorf("url is required")
	}
	if r.Port == "" {
		return fmt.Errorf("port is required")
	}
	if r.Token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}

type MicrosoftSentinelRequest struct {
	GUID              string            `json:"guid"`
	Name              string            `json:"name"`
	WorkSpaceID       string            `json:"workSpaceID"`
	PrimaryKey        string            `json:"primaryKey"`
	TestMessageStatus TestMessageStatus `json:"testMessageStatus"`
}

func (r *MicrosoftSentinelRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.WorkSpaceID == "" {
		return fmt.Errorf("workSpaceID is required")
	}
	if r.PrimaryKey == "" {
		return fmt.Errorf("primaryKey is required")
	}
	return nil
}

type WebhookRequest struct {
	GUID              string            `json:"guid"`
	Name              string            `json:"name"`
	WebhookURL        string            `json:"webhookURL"`
	TestMessageStatus TestMessageStatus `json:"testMessageStatus"`
}

func (r *WebhookRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.WebhookURL == "" {
		return fmt.Errorf("webhookURL is required")
	}
	return nil
}

type DeleteRequest struct {
	GUID string `json:"guid"`
}

func (s *SIEMIntegration) GetProvider() SIEMProvider {
	return s.Provider
}

func (s *SIEMIntegration) GetCustomerGUID() string {
	return s.CustomerGUID
}

func (s *SIEMIntegration) GetTestMessageStatus() TestMessageStatus {
	return s.TestMessageStatus
}

func (s *SIEMIntegration) SetTestMessageStatus(status TestMessageStatus) {
	s.TestMessageStatus = status
}

func (s *SIEMIntegration) GetWebhookConfig() (*WebhookConfig, error) {
	if s.Provider != SIEMProviderWebhook {
		return nil, fmt.Errorf("invalid provider for Webhook config: %s", s.Provider)
	}
	config := &WebhookConfig{}
	return config, nil
}

func (s *SIEMIntegration) GetSumoLogicConfig() (*SumoLogicConfig, error) {
	if s.Provider != SIEMProviderSumo {
		return nil, fmt.Errorf("invalid provider for SumoLogic config: %s", s.Provider)
	}
	config := &SumoLogicConfig{}
	if httpSourceAddress, ok := s.Configuration["httpSourceAddress"].(string); ok {
		config.HTTPSourceAddress = httpSourceAddress
	}
	return config, nil
}

func (s *SIEMIntegration) GetSplunkConfig() (*SplunkConfig, error) {
	if s.Provider != SIEMProviderSplunk {
		return nil, fmt.Errorf("invalid provider for Splunk config: %s", s.Provider)
	}
	config := &SplunkConfig{}
	if url, ok := s.Configuration["url"].(string); ok {
		config.URL = url
	}
	if port, ok := s.Configuration["port"].(string); ok {
		config.Port = port
	}
	if token, ok := s.Configuration["token"].(string); ok {
		config.Token = token
	}
	return config, nil
}

func (s *SIEMIntegration) GetMicrosoftSentinelConfig() (*MicrosoftSentinelConfig, error) {
	if s.Provider != SIEMProviderMicrosoftSentinel {
		return nil, fmt.Errorf("invalid provider for Microsoft Sentinel config: %s", s.Provider)
	}
	config := &MicrosoftSentinelConfig{}
	if workSpaceID, ok := s.Configuration["workSpaceID"].(string); ok {
		config.WorkSpaceID = workSpaceID
	}
	if primaryKey, ok := s.Configuration["primaryKey"].(string); ok {
		config.PrimaryKey = primaryKey
	}
	return config, nil
}

func (s *SIEMIntegration) SetSumoLogicConfig(config *SumoLogicConfig) {
	if config == nil {
		return
	}
	s.Provider = SIEMProviderSumo
	s.Configuration = map[string]interface{}{
		"httpSourceAddress": config.HTTPSourceAddress,
	}
}

func (s *SIEMIntegration) SetSplunkConfig(config *SplunkConfig) {
	if config == nil {
		return
	}
	s.Provider = SIEMProviderSplunk
	s.Configuration = map[string]interface{}{
		"url":   config.URL,
		"port":  config.Port,
		"token": config.Token,
	}
}

func (s *SIEMIntegration) SetMicrosoftSentinelConfig(config *MicrosoftSentinelConfig) {
	if config == nil {
		return
	}
	s.Provider = SIEMProviderMicrosoftSentinel
	s.Configuration = map[string]interface{}{
		"workSpaceID": config.WorkSpaceID,
		"primaryKey":  config.PrimaryKey,
	}
}

func (s *SIEMIntegration) SetWebhookConfig(config *WebhookConfig) {
	if config == nil {
		return
	}
	s.Provider = SIEMProviderWebhook
	s.Configuration = map[string]interface{}{
		"webhookURL": config.WebhookURL,
	}
}

func (s *SIEMIntegration) GetName() string {
	return s.Name
}

func (s *SIEMIntegration) SetName(name string) {
	s.Name = name
}

func (s *SIEMIntegration) GetGUID() string {
	return s.GUID
}

func (s *SIEMIntegration) SetGUID(guid string) {
	s.GUID = guid
}
