package containerscan

import "time"

type TimeValueCordindate struct {
	Value       interface{} `json:"value"`
	Timestamp   time.Time   `json:"timestamp"`
	UniqueValue interface{} `json:"uniqueValue"`
}

type ContainerSummmaryTimeValueCordindate struct {
	TimeValueCordindate `json:",inline"`
	ImageTag            string `json:"imageTag"`
	ImageHash           string `json:"imageHash"`
}

type SeverityTimeValue struct {
	Cords    []ContainerSummmaryTimeValueCordindate `json:"cords"`
	Severity string                                 `json:"severity"`
}
