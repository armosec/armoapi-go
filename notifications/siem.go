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

// SIEMConfig defines the interface for SIEM provider configurations
type SIEMConfig interface {
	Validate() error
}

type SumoLogicConfig struct {
	HttpSourceAddress string `json:"httpSourceAddress" bson:"httpSourceAddress"`
}

func (c *SumoLogicConfig) Validate() error {
	if c.HttpSourceAddress == "" {
		return fmt.Errorf("httpSourceAddress is required")
	}
	return nil
}

type SplunkConfig struct {
	URL   string `json:"url" bson:"url"`
	Port  string `json:"port" bson:"port"`
	Token string `json:"token" bson:"token"`
}

func (c *SplunkConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("url is required")
	}
	if c.Port == "" {
		return fmt.Errorf("port is required")
	}
	if c.Token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}

type MicrosoftSentinelConfig struct {
	WorkSpaceID string `json:"workSpaceID" bson:"workSpaceID"`
	PrimaryKey  string `json:"primaryKey" bson:"primaryKey"`
}

func (c *MicrosoftSentinelConfig) Validate() error {
	if c.WorkSpaceID == "" {
		return fmt.Errorf("workSpaceID is required")
	}
	if c.PrimaryKey == "" {
		return fmt.Errorf("primaryKey is required")
	}
	return nil
}

type WebhookConfig struct {
	WebhookURL string             `json:"webhookURL"`
	Headers    *map[string]string `json:"headers,omitempty"`
}

func (c *WebhookConfig) Validate() error {
	if c.WebhookURL == "" {
		return fmt.Errorf("webhookURL is required")
	}
	return nil
}

type SIEMIntegration struct {
	armotypes.PortalBase `json:",inline" bson:"inline"`
	Provider             SIEMProvider           `json:"provider" bson:"provider"`
	Configuration        map[string]interface{} `json:"configuration" bson:"configuration"`
	IsEnabled            bool                   `json:"isEnabled" bson:"isEnabled"`
	TestMessageStatus    TestMessageStatus      `json:"testMessageStatus" bson:"testMessageStatus"`
	UpdatedBy            string                 `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	CreationTime         string                 `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
}

type SumoLogicRequest struct {
	GUID          string           `json:"guid"`
	Name          string           `json:"name"`
	IsEnabled     bool             `json:"isEnabled"`
	Configuration *SumoLogicConfig `json:"configuration"`
}

func (r *SumoLogicRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	return r.Configuration.Validate()
}

type SplunkRequest struct {
	GUID          string       `json:"guid"`
	Name          string       `json:"name"`
	IsEnabled     bool         `json:"isEnabled"`
	Configuration SplunkConfig `json:"configuration"`
}

func (r *SplunkRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	return r.Configuration.Validate()
}

type MicrosoftSentinelRequest struct {
	GUID          string                  `json:"guid"`
	Name          string                  `json:"name"`
	IsEnabled     bool                    `json:"isEnabled"`
	Configuration MicrosoftSentinelConfig `json:"configuration"`
}

func (r *MicrosoftSentinelRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	return r.Configuration.Validate()
}

type WebhookRequest struct {
	GUID          string        `json:"guid"`
	Name          string        `json:"name"`
	IsEnabled     bool          `json:"isEnabled"`
	Configuration WebhookConfig `json:"configuration"`
}

func (r *WebhookRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	return r.Configuration.Validate()
}

type SiemIntegrationDeleteRequest struct {
	GUID string `json:"guid"`
}

func (s *SiemIntegrationDeleteRequest) Validate() error {
	if s.GUID == "" {
		return fmt.Errorf("guid is required")
	}
	return nil
}

func (s *SIEMIntegration) GetProvider() SIEMProvider {
	return s.Provider
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
		config.HttpSourceAddress = httpSourceAddress
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
		"httpSourceAddress": config.HttpSourceAddress,
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
