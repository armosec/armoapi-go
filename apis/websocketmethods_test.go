package apis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCronJobParams(t *testing.T) {
	c := Command{
		CommandName: TypeRunKubescapeJob,
	}
	cjp := CronJobParams{
		CronTabSchedule: "* * * * *",
	}
	c.SetCronJobParams(cjp)
	cjpr := c.GetCronJobParams()
	assert.Equal(t, "* * * * *", cjpr.CronTabSchedule)
}

func TestLabels(t *testing.T) {
	c := Command{
		CommandName: TypeAttachWorkload,
	}
	cjp := map[string]string{
		"app": "game",
	}
	c.SetLabels(cjp)
	cjpr := c.GetLabels()
	assert.Equal(t, "game", cjpr["app"])
}

func TestFieldSelector(t *testing.T) {
	c := Command{
		CommandName: TypeAttachWorkload,
	}
	cjp := map[string]string{
		"app": "game",
	}
	c.SetFieldSelector(cjp)
	cjpr := c.GetFieldSelector()
	assert.Equal(t, "game", cjpr["app"])
}
