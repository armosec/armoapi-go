package armotypes

import (
	"fmt"
	"strings"
	"time"
)

type DataType string

// PortalCache is an auxiliary structure to store cache data
type PortalCache struct {
	GUID         string      `json:"guid" bson:"guid"`
	Name         string      `json:"name,omitempty" bson:"name,omitempty"`
	DataType     DataType    `json:"dataType,omitempty" bson:"dataType,omitempty"`
	Data         interface{} `json:"data,omitempty" bson:"data,omitempty"`
	CreationTime string      `json:"creationTime" bson:"creationTime"`
	UpdatedTime  string      `json:"lastUpdated,omitempty" bson:"lastUpdated,omitempty"`
	ExpiryTime   time.Time   `json:"expiryTime,omitempty" bson:"expiryTime,omitempty"`
}

func (c *PortalCache) SetExpiryTime(expiryTime time.Time) {
	c.ExpiryTime = expiryTime
}

func (c *PortalCache) SetTTL(ttl time.Duration) {
	c.ExpiryTime = time.Now().Add(ttl)
}

func GetCacheData[T any](cache *PortalCache) (T, error) {
	if cache.Data == nil {
		var val T
		return val, fmt.Errorf("cache data is nil")
	}
	val, ok := cache.Data.(T)
	if !ok {
		return val, fmt.Errorf("cache data is not of type %T", val)
	}
	return val, nil
}

// Generates cache data-type form properties
func MakeCacheDataTypeV1(service, customerGUID, domain, propose, version string) DataType {
	return DataType(fmt.Sprintf("datatypeV1:%s-%s-%s-%s-%s", service, customerGUID, domain, propose, version))
}

// Parse cache data-type
func ParseCacheDataTypeV1(dataType DataType) (service, customerGUID, domain, propose, version string, err error) {
	parts := strings.Split(string(dataType), "-")
	if len(parts) != 5 {
		err = fmt.Errorf("failed to parse dataType: %s", dataType)
		return
	}
	prefix := "datatypeV1:"
	if !strings.HasPrefix(string(dataType), prefix) {
		err = fmt.Errorf("invalid dataType format: %s", dataType)
		return
	}
	service = parts[0][len(prefix):]
	customerGUID = parts[1]
	domain = parts[2]
	propose = parts[3]
	version = parts[4]
	return
}
