package apis

import (
	"net/http"

	"github.com/docker/docker/api/types/registry"
)

// Connector - interface for any connector (BE/Portal and so on)
type Connector interface {

	//may used for a more generic httpsend interface based method
	GetBaseURL() string
	GetLoginObj() *LoginObject
	GetClient() *http.Client

	Login() error
	IsExpired() bool

	HTTPSend(httpverb string,
		endpoint string,
		payload []byte,
		f HTTPReqFunc,
		qryData interface{}) ([]byte, error)
}

type ImageScanCommand interface {
	GetWlid() string
	GetImageHash() string
	GetCreds() *registry.AuthConfig
	GetCredentialsList() []registry.AuthConfig
	SetCredentialsList([]registry.AuthConfig)
	GetArgs() map[string]interface{}
	SetArgs(map[string]interface{})
	GetSession() SessionChain
	SetSession(SessionChain)
	GetImageTag() string
	SetImageTag(string)
	GetJobID() string
	SetJobID(string)
	GetParentJobID() string
	SetParentJobID(string)
}
