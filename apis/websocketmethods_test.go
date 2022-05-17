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
