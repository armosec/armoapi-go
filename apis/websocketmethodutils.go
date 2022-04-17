package apis

import (
	"fmt"
)

func SIDFallback(c *Command) {
	if c.GetID() == "" {
		sid, err := getSIDFromArgs(c.Args)
		if err != nil || sid == "" {
			return
		}
		c.Sid = sid
	}
}

func getSIDFromArgs(args map[string]interface{}) (string, error) {
	sidInterface, ok := args["sid"]
	if !ok {
		return "", nil
	}
	sid, ok := sidInterface.(string)
	if !ok || sid == "" {
		return "", fmt.Errorf("sid found in args but empty")
	}
	// if _, err := secrethandling.SplitSecretID(sid); err != nil {
	// 	return "", err
	// }
	return sid, nil
}
