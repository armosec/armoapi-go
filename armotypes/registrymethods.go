package armotypes

import (
	"encoding/json"
	"errors"
)

var RegistryTypeMap = map[RegistryProvider]ContainerImageRegistry{
	AWS:    new(AWSImageRegistry),
	Azure:  new(AzureImageRegistry),
	Google: new(GoogleImageRegistry),
	Harbor: new(HarborImageRegistry),
	Quay:   new(QuayImageRegistry),
}

func UnmarshalRegistry(payload []byte) (ContainerImageRegistry, error) {
	var providerHolder struct {
		Provider string `json:"provider"`
	}
	if err := json.Unmarshal(payload, &providerHolder); err != nil {
		return nil, err
	}

	registry := RegistryTypeMap[RegistryProvider(providerHolder.Provider)]
	if err := json.Unmarshal(payload, &registry); err != nil {
		return nil, err
	}
	return registry, nil
}

func (base *BaseContainerImageRegistry) ValidateBase() error {
	if base.ClusterName == "" {
		return errors.New("clusterName is empty")
	}
	if len(base.Repositories) == 0 {
		return errors.New("repositories is empty")
	}
	return nil
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
	secretMap, err := decodeSecretMapFromInterface(value)
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

func (aws *AWSImageRegistry) GetBase() *BaseContainerImageRegistry {
	return &aws.BaseContainerImageRegistry
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
	secretMap, err := decodeSecretMapFromInterface(value)
	if err != nil {
		return err
	}
	azure.LoginServer = secretMap["loginServer"]
	azure.Username = secretMap["username"]
	azure.AccessToken = secretMap["accessToken"]
	return nil
}

func (azure *AzureImageRegistry) GetBase() *BaseContainerImageRegistry {
	return &azure.BaseContainerImageRegistry
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

func (google *GoogleImageRegistry) MaskSecret() {

}

func (google *GoogleImageRegistry) ExtractSecret() interface{} {
	return map[string]string{
		"registryURI": google.RegistryURI,
	}
}

func (google *GoogleImageRegistry) FillSecret(value interface{}) error {
	secretMap, err := decodeSecretMapFromInterface(value)
	if err != nil {
		return err
	}
	google.RegistryURI = secretMap["registryURI"]
	return nil
}

func (google *GoogleImageRegistry) GetBase() *BaseContainerImageRegistry {
	return &google.BaseContainerImageRegistry
}

func (google *GoogleImageRegistry) Validate() error {
	if err := google.GetBase().ValidateBase(); err != nil {
		return err
	}
	if google.RegistryURI == "" {
		return errors.New("registryURI is empty")
	}
	return nil
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
	secretMap, err := decodeSecretMapFromInterface(value)
	if err != nil {
		return err
	}
	harbor.InstanceURL = secretMap["instanceURL"]
	harbor.Username = secretMap["username"]
	harbor.Password = secretMap["password"]
	return nil
}

func (harbor *HarborImageRegistry) GetBase() *BaseContainerImageRegistry {
	return &harbor.BaseContainerImageRegistry
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

func (quay *QuayImageRegistry) MaskSecret() {
	quay.RobotAccountToken = ""
}

func (quay *QuayImageRegistry) ExtractSecret() interface{} {
	return map[string]string{
		"containerRegistryName": quay.ContainerRegistryName,
		"robotAccountName":      quay.RobotAccountName,
		"robotAccountToken":     quay.RobotAccountToken,
	}
}

func (quay *QuayImageRegistry) FillSecret(value interface{}) error {
	secretMap, err := decodeSecretMapFromInterface(value)
	if err != nil {
		return err
	}
	quay.ContainerRegistryName = secretMap["containerRegistryName"]
	quay.RobotAccountName = secretMap["robotAccountName"]
	quay.RobotAccountToken = secretMap["robotAccountToken"]
	return nil
}

func (quay *QuayImageRegistry) GetBase() *BaseContainerImageRegistry {
	return &quay.BaseContainerImageRegistry
}

func (quay *QuayImageRegistry) Validate() error {
	if err := quay.GetBase().ValidateBase(); err != nil {
		return err
	}
	if quay.RobotAccountName == "" {
		return errors.New("robotAccountName is empty")
	}
	if quay.RobotAccountToken == "" {
		return errors.New("robotAccountToken is empty")
	}
	return nil
}

func decodeSecretMapFromInterface(value interface{}) (map[string]string, error) {
	var res map[string]string
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
