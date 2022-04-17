package apis

import (
	"encoding/json"

	"github.com/armosec/armoapi-go/armotypes"
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
	if c.Args != nil {
		if ilabels, ok := c.Args["labels"]; ok {
			labels := map[string]string{}
			if b, e := json.Marshal(ilabels); e == nil {
				if e = json.Unmarshal(b, &labels); e == nil {
					return labels
				}
			}
		}
	}
	return map[string]string{}
}

func (c *Command) SetLabels(labels map[string]string) {
	if c.Args == nil {
		c.Args = make(map[string]interface{})
	}
	c.Args["labels"] = labels
}

func (c *Command) GetFieldSelector() map[string]string {
	if c.Args != nil {
		if ilabels, ok := c.Args["fieldSelector"]; ok {
			labels := map[string]string{}
			if b, e := json.Marshal(ilabels); e == nil {
				if e = json.Unmarshal(b, &labels); e == nil {
					return labels
				}
			}
		}
	}
	return map[string]string{}
}

func (c *Command) SetFieldSelector(labels map[string]string) {
	if c.Args == nil {
		c.Args = make(map[string]interface{})
	}
	c.Args["fieldSelector"] = labels
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
