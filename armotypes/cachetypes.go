package armotypes

import (
	"time"
)

type DataType string

// PortalCache is an auxiliary structure to store cache data
type PortalCache[T any] struct {
	GUID         string    `json:"guid" bson:"guid"`
	Name         string    `json:"name,omitempty" bson:"name,omitempty"`
	DataType     DataType  `json:"dataType,omitempty" bson:"dataType,omitempty"`
	Data         T         `json:"data,omitempty" bson:"data,omitempty"`
	CreationTime string    `json:"creationTime" bson:"creationTime"`
	UpdatedTime  string    `json:"lastUpdated,omitempty" bson:"lastUpdated,omitempty"`
	ExpiryTime   time.Time `json:"expiryTime,omitempty" bson:"expiryTime,omitempty"`
}

func (c *PortalCache[T]) SetExpiryTime(expiryTime time.Time) {
	c.ExpiryTime = expiryTime
}

func (c *PortalCache[T]) SetTTL(ttl time.Duration) {
	c.ExpiryTime = time.Now().UTC().Add(ttl)
}

func (c *PortalCache[T]) GetTimestampFieldName() string {
	return "creationTime"
}
