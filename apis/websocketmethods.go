package apis

import (
	"encoding/json"

	"github.com/armosec/armoapi-go/identifiers"
	"github.com/docker/docker/api/types/registry"
)

const (
	CommandDeprecatedArgsJobParams string = "kubescapeJobParams"

	commandArgsJobParams     string = "jobParams"
	commandArgsLabels        string = "labels"
	commandArgsFieldSelector string = "fieldSelector"
)

func (c *Command) DeepCopy() *Command {
	newCommand := &Command{}
	newCommand.CommandName = c.CommandName
	newCommand.ResponseID = c.ResponseID
	newCommand.Wlid = c.Wlid
	newCommand.WildWlid = c.WildWlid
	newCommand.Designators = c.Designators
	if c.Args != nil {
		newCommand.Args = make(map[string]interface{})
		for i, j := range c.Args {
			newCommand.Args[i] = j
		}
	}
	return newCommand
}

func (c *Command) GetLabels() map[string]string {
	labels := map[string]string{}
	if f := c.GetArg(commandArgsLabels); f != nil {
		b, err := json.Marshal(f)
		if err != nil {
			return labels
		}
		if err := json.Unmarshal(b, &labels); err != nil {
			return labels
		}
	}
	return labels

}

func (c *Command) SetLabels(labels map[string]string) {
	c.SetArg(commandArgsLabels, labels)
}

func (c *Command) GetFieldSelector() map[string]string {
	fieldSelector := map[string]string{}
	if f := c.GetArg(commandArgsFieldSelector); f != nil {
		b, err := json.Marshal(f)
		if err != nil {
			return fieldSelector
		}
		if err := json.Unmarshal(b, &fieldSelector); err != nil {
			return fieldSelector
		}
	}
	return fieldSelector
}

func (c *Command) SetFieldSelector(labels map[string]string) {
	c.SetArg(commandArgsFieldSelector, labels)
}
func (c *Command) SetCronJobParams(cjParams CronJobParams) {
	c.SetArg(commandArgsJobParams, cjParams)
}

func (c *Command) GetCronJobParams() *CronJobParams {
	cjParams := &CronJobParams{}
	if icjParams := c.GetArg(commandArgsJobParams); icjParams != nil {
		b, err := json.Marshal(icjParams)
		if err != nil {
			return cjParams
		}
		if err := json.Unmarshal(b, cjParams); err != nil {
			return cjParams
		}
	}
	return cjParams
}

func (c *Command) SetArg(key string, value interface{}) {
	if c.Args == nil {
		c.Args = make(map[string]interface{})
	}
	c.Args[key] = value
}

func (c *Command) GetArg(key string) interface{} {
	if c.Args == nil {
		return nil
	}
	v, ok := c.Args[key]
	if !ok {
		return nil
	}
	return v
}

func (c *Command) GetID() string {
	if len(c.Designators) > 0 {
		return identifiers.DesignatorsToken
	}
	if c.WildWlid != "" {
		return c.WildWlid
	}
	if c.WildSid != "" {
		return c.WildSid
	}
	if c.Wlid != "" {
		return c.Wlid
	}
	if c.Sid != "" {
		return c.Sid
	}
	return ""
}

func (c *Command) Json() string {
	b, _ := json.Marshal(*c)
	return string(b)
}

// RegistryScanCommand implementation for ImageScanCommand interface

var _ ImageScanCommand = &RegistryScanCommand{}

func (r *RegistryScanCommand) GetWlid() string {
	return ""
}

func (r *RegistryScanCommand) GetCredentialsList() []registry.AuthConfig {
	return r.ImageScanParams.Credentialslist
}

func (r *RegistryScanCommand) SetCredentialsList(credentialslist []registry.AuthConfig) {
	r.ImageScanParams.Credentialslist = credentialslist
}

func (r *RegistryScanCommand) GetArgs() map[string]interface{} {
	return r.ImageScanParams.Args
}

func (r *RegistryScanCommand) SetArgs(args map[string]interface{}) {
	r.ImageScanParams.Args = args
}

func (r *RegistryScanCommand) GetSession() SessionChain {
	return r.ImageScanParams.Session
}

func (r *RegistryScanCommand) SetSession(session SessionChain) {
	r.ImageScanParams.Session = session
}

func (r *RegistryScanCommand) GetImageTag() string {
	return r.ImageScanParams.ImageTag
}

func (r *RegistryScanCommand) SetImageTag(imageTag string) {
	r.ImageScanParams.ImageTag = imageTag
}

func (r *RegistryScanCommand) GetJobID() string {
	return r.ImageScanParams.JobID
}

func (r *RegistryScanCommand) SetJobID(jobID string) {
	r.ImageScanParams.JobID = jobID
}

func (r *RegistryScanCommand) GetParentJobID() string {
	return r.ImageScanParams.ParentJobID
}

func (r *RegistryScanCommand) SetParentJobID(parentJobID string) {
	r.ImageScanParams.ParentJobID = parentJobID
}

func (r *RegistryScanCommand) GetCreds() *registry.AuthConfig {
	return nil
}

func (r *RegistryScanCommand) GetImageHash() string {
	return ""
}

// WebsocketScanCommand implementation for ImageScanCommand interface

var _ ImageScanCommand = &WebsocketScanCommand{}

func (c *WebsocketScanCommand) GetCredentialsList() []registry.AuthConfig {
	return c.Credentialslist
}

func (c *WebsocketScanCommand) SetCredentialsList(credentialslist []registry.AuthConfig) {
	c.Credentialslist = credentialslist
}

func (c *WebsocketScanCommand) GetArgs() map[string]interface{} {
	return c.Args
}

func (c *WebsocketScanCommand) SetArgs(args map[string]interface{}) {
	c.Args = args
}

func (c *WebsocketScanCommand) GetSession() SessionChain {
	return c.Session
}

func (c *WebsocketScanCommand) SetSession(session SessionChain) {
	c.Session = session
}

func (c *WebsocketScanCommand) GetImageTag() string {
	return c.ImageTag
}

func (c *WebsocketScanCommand) SetImageTag(imageTag string) {
	c.ImageTag = imageTag
}

func (c *WebsocketScanCommand) GetJobID() string {
	return c.JobID
}

func (c *WebsocketScanCommand) SetJobID(jobID string) {
	c.JobID = jobID
}

func (c *WebsocketScanCommand) GetParentJobID() string {
	return c.ParentJobID
}

func (c *WebsocketScanCommand) SetParentJobID(parentJobID string) {
	c.ParentJobID = parentJobID
}

func (c *WebsocketScanCommand) GetImageHash() string {
	return c.ImageHash
}

func (c *WebsocketScanCommand) GetCreds() *registry.AuthConfig {
	return c.Credentials
}

func (c *WebsocketScanCommand) GetWlid() string {
	return c.Wlid
}
