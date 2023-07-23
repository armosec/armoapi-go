package armotypes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetUpdatedTime(t *testing.T) {
	now := time.Now()
	nowString := now.UTC().Format(time.RFC3339)
	validDateString := "2022-12-26T15:05:23Z"
	validDate, _ := time.Parse(time.RFC3339, validDateString)

	type testCase struct {
		name     string
		time     *time.Time
		expected PortalBase
	}
	testTable := []testCase{
		{
			name:     "valid time",
			time:     &validDate,
			expected: PortalBase{UpdatedTime: validDateString},
		},
		{
			name:     "default time",
			time:     nil,
			expected: PortalBase{UpdatedTime: nowString},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			p := PortalBase{}
			p.SetUpdatedTime(test.time)
			assert.Equal(t, test.expected, p)
		})
	}
}

func TestValidateContainerScanID(t *testing.T) {
	type args struct {
		containerScanID string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid containerScanID",
			args: args{
				containerScanID: "9711c327-1a08-487e-b24a-72128712ef2d",
			},
			want: true,
		},
		{
			name: "containerScanID with a slash is invalid",
			args: args{
				containerScanID: "foo/bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ValidateContainerScanID(tt.args.containerScanID), "ValidateContainerScanID(%v)", tt.args.containerScanID)
		})
	}
}
