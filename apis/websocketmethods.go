package apis

import (
	"encoding/json"

	"github.com/armosec/armoapi-go/armotypes"
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
		return armotypes.DesignatorsToken
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

func (safeMode *SafeMode) Json() string {
	b, err := json.Marshal(*safeMode)
	if err != nil {
		return ""
	}
	return string(b)
}
