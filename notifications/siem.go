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
	Provider             string                 `json:"provider" bson:"provider"`
	Configuration        map[string]interface{} `json:"configuration" bson:"configuration"`
	IsEnabled            bool                   `json:"isEnabled" bson:"isEnabled"`
}

func (s *SIEMIntegration) GetProvider() string {
	return s.Provider
}

func (s *SIEMIntegration) GetCustomerGUID() string {
	return s.CustomerGUID
}


func (s *SIEMIntegration) GetSumoLogicConfig() (*SumoLogicConfig, error) {
	if s.Provider != string(SIEMProviderSumo) {
		return nil, fmt.Errorf("invalid provider for SumoLogic config: %s", s.Provider)
	}
	config := &SumoLogicConfig{}
	if httpSourceAddress, ok := s.Configuration["httpSourceAddress"].(string); ok {
		config.HTTPSourceAddress = httpSourceAddress
	}
	return config, nil
}

func (s *SIEMIntegration) GetSplunkConfig() (*SplunkConfig, error) {
	if s.Provider != string(SIEMProviderSplunk) {
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
	if s.Provider != string(SIEMProviderMicrosoftSentinel) {
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
	s.Provider = string(SIEMProviderSumo)
	s.Configuration = map[string]interface{}{
		"httpSourceAddress": config.HTTPSourceAddress,
	}
}

func (s *SIEMIntegration) SetSplunkConfig(config *SplunkConfig) {
	if config == nil {
		return
	}
	s.Provider = string(SIEMProviderSplunk)
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
	s.Provider = string(SIEMProviderMicrosoftSentinel)
	s.Configuration = map[string]interface{}{
		"workSpaceID": config.WorkSpaceID,
		"primaryKey":  config.PrimaryKey,
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
