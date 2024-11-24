package armotypes

import (
	"encoding/json"
	"errors"
)

var RegistryTypeMap = map[RegistryProvider]func() ContainerImageRegistry{
	AWS:    func() ContainerImageRegistry { return new(AWSImageRegistry) },
	Azure:  func() ContainerImageRegistry { return new(AzureImageRegistry) },
	Google: func() ContainerImageRegistry { return new(GoogleImageRegistry) },
	Harbor: func() ContainerImageRegistry { return new(HarborImageRegistry) },
	Quay:   func() ContainerImageRegistry { return new(QuayImageRegistry) },
	Nexus:  func() ContainerImageRegistry { return new(NexusImageRegistry) },
}

func UnmarshalRegistry(payload []byte) (ContainerImageRegistry, error) {
	var providerHolder struct {
		Provider string `json:"provider"`
	}
	if err := json.Unmarshal(payload, &providerHolder); err != nil {
		return nil, err
	}

	registry := RegistryTypeMap[RegistryProvider(providerHolder.Provider)]()
	if err := json.Unmarshal(payload, &registry); err != nil {
		return nil, err
	}
	return registry, nil
}

func (base *BaseContainerImageRegistry) ValidateBase() error {
	if base.ClusterName == "" {
		return errors.New("clusterName is empty")
	}
	return nil
}

func (b *BaseContainerImageRegistry) GetBase() *BaseContainerImageRegistry {
	return b
}
func (b *BaseContainerImageRegistry) SetBase(base *BaseContainerImageRegistry) {
	*b = *base
}

func (aws *AWSImageRegistry) MaskSecret() {
	aws.SecretAccessKey = ""
	aws.RoleARN = ""
}

func (aws *AWSImageRegistry) ExtractSecret() interface{} {
	return map[string]string{
		"registry":        aws.Registry,
		"registryRegion":  aws.RegistryRegion,
		"accessKeyID":     aws.AccessKeyID,
		"secretAccessKey": aws.SecretAccessKey,
		"roleARN":         aws.RoleARN,
	}
}

func (aws *AWSImageRegistry) FillSecret(value interface{}) error {
	secretMap, err := decodeSecretFromInterface[map[string]string](value)
	if err != nil {
		return err
	}
	aws.Registry = secretMap["registry"]
	aws.RegistryRegion = secretMap["registryRegion"]
	aws.AccessKeyID = secretMap["accessKeyID"]
	aws.SecretAccessKey = secretMap["secretAccessKey"]
	aws.RoleARN = secretMap["roleARN"]
	return nil
}

func (aws *AWSImageRegistry) Validate() error {
	if err := aws.GetBase().ValidateBase(); err != nil {
		return err
	}

	if aws.Registry == "" {
		return errors.New("registry is empty")
	}
	if aws.RegistryRegion == "" {
		return errors.New("registryRegion is empty")
	}
	if (aws.AccessKeyID == "" || aws.SecretAccessKey == "") && aws.RoleARN == "" {
		return errors.New("missing authentication data")
	}
	return nil
}

func (aws *AWSImageRegistry) GetDisplayName() string {
	return aws.Registry
}

func (azure *AzureImageRegistry) MaskSecret() {
	azure.AccessToken = ""
}

func (azure *AzureImageRegistry) ExtractSecret() interface{} {
	return map[string]string{
		"loginServer": azure.LoginServer,
		"username":    azure.Username,
		"accessToken": azure.AccessToken,
	}
}

func (azure *AzureImageRegistry) FillSecret(value interface{}) error {
	secretMap, err := decodeSecretFromInterface[map[string]string](value)
	if err != nil {
		return err
	}
	azure.LoginServer = secretMap["loginServer"]
	azure.Username = secretMap["username"]
	azure.AccessToken = secretMap["accessToken"]
	return nil
}

func (azure *AzureImageRegistry) Validate() error {
	if err := azure.GetBase().ValidateBase(); err != nil {
		return err
	}

	if azure.LoginServer == "" {
		return errors.New("loginServer is empty")
	}
	if azure.Username == "" {
		return errors.New("username is empty")
	}
	if azure.AccessToken == "" {
		return errors.New("accessToken is empty")
	}
	return nil
}

func (azure *AzureImageRegistry) GetDisplayName() string {
	return azure.LoginServer
}

func (google *GoogleImageRegistry) MaskSecret() {
	google.Key = nil
}

func (google *GoogleImageRegistry) ExtractSecret() interface{} {
	return map[string]interface{}{
		"registryURI": google.RegistryURI,
		"key":         google.Key,
	}
}

func (google *GoogleImageRegistry) FillSecret(value interface{}) error {
	secretMap, err := decodeSecretFromInterface[map[string]interface{}](value)
	if err != nil {
		return err
	}
	google.RegistryURI = secretMap["registryURI"].(string)
	google.Key = secretMap["key"].(map[string]interface{})
	return nil
}

func (google *GoogleImageRegistry) Validate() error {
	if err := google.GetBase().ValidateBase(); err != nil {
		return err
	}
	if google.RegistryURI == "" {
		return errors.New("registryURI is empty")
	}
	if len(google.Key) == 0 {
		return errors.New("json key is empty")
	}
	return nil
}

func (google *GoogleImageRegistry) GetDisplayName() string {
	return google.RegistryURI
}

func (harbor *HarborImageRegistry) MaskSecret() {
	harbor.Password = ""
}

func (harbor *HarborImageRegistry) ExtractSecret() interface{} {
	return map[string]string{
		"instanceURL": harbor.InstanceURL,
		"username":    harbor.Username,
		"password":    harbor.Password,
	}
}

func (harbor *HarborImageRegistry) FillSecret(value interface{}) error {
	secretMap, err := decodeSecretFromInterface[map[string]string](value)
	if err != nil {
		return err
	}
	harbor.InstanceURL = secretMap["instanceURL"]
	harbor.Username = secretMap["username"]
	harbor.Password = secretMap["password"]
	return nil
}

func (harbor *HarborImageRegistry) Validate() error {
	if err := harbor.GetBase().ValidateBase(); err != nil {
		return err
	}
	if harbor.InstanceURL == "" {
		return errors.New("instanceURL is empty")
	}
	if harbor.Username == "" {
		return errors.New("username is empty")
	}
	if harbor.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

func (harbor *HarborImageRegistry) GetDisplayName() string {
	return harbor.InstanceURL
}

const (
	containerRegistryName = "containerRegistryName"
	robotAccountName      = "robotAccountName"
	robotAccountToken     = "robotAccountToken"
)

func (quay *QuayImageRegistry) MaskSecret() {
	quay.RobotAccountToken = ""
}

func (quay *QuayImageRegistry) ExtractSecret() interface{} {
	return map[string]string{
		containerRegistryName: quay.ContainerRegistryName,
		robotAccountName:      quay.RobotAccountName,
		robotAccountToken:     quay.RobotAccountToken,
	}
}

func (quay *QuayImageRegistry) FillSecret(value interface{}) error {
	secretMap, err := decodeSecretFromInterface[map[string]string](value)
	if err != nil {
		return err
	}
	quay.ContainerRegistryName = secretMap[containerRegistryName]
	quay.RobotAccountName = secretMap[robotAccountName]
	quay.RobotAccountToken = secretMap[robotAccountToken]
	return nil
}

func (quay *QuayImageRegistry) Validate() error {
	if err := quay.GetBase().ValidateBase(); err != nil {
		return err
	}
	if quay.ContainerRegistryName == "" {
		return errors.New("container registry name is empty")
	}
	if quay.RobotAccountName == "" {
		return errors.New("robot account name is empty")
	}
	if quay.RobotAccountToken == "" {
		return errors.New("robot account token is empty")
	}
	return nil
}

func (quay *QuayImageRegistry) GetDisplayName() string {
	return quay.ContainerRegistryName
}

func (nexus *NexusImageRegistry) MaskSecret() {
	nexus.Password = ""
}

func (nexus *NexusImageRegistry) ExtractSecret() interface{} {
	return map[string]string{
		"registryURL": nexus.RegistryURL,
		"username":    nexus.Username,
		"password":    nexus.Password,
	}
}

func (nexus *NexusImageRegistry) FillSecret(value interface{}) error {
	secretMap, err := decodeSecretFromInterface[map[string]string](value)
	if err != nil {
		return err
	}
	nexus.RegistryURL = secretMap["registryURL"]
	nexus.Username = secretMap["username"]
	nexus.Password = secretMap["password"]
	return nil
}

func (nexus *NexusImageRegistry) Validate() error {
	if err := nexus.GetBase().ValidateBase(); err != nil {
		return err
	}
	if nexus.RegistryURL == "" {
		return errors.New("registry url is empty")
	}
	if nexus.Username == "" {
		return errors.New("username is empty")
	}
	if nexus.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

func (nexus *NexusImageRegistry) GetDisplayName() string {
	return nexus.RegistryURL
}

func decodeSecretFromInterface[T any](value interface{}) (T, error) {
	var res T
	if value == nil {
		return res, errors.New("got an empty value")
	}
	updatedJson, err := json.Marshal(value)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(updatedJson, &res)
	return res, err
}
