package armotypes

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const sep = "␟"

type File struct {
	Path       string         `json:"path,omitempty" bson:"path,omitempty"`
	Size       int64          `json:"size,omitempty" bson:"size,omitempty"`
	Hashes     FileHashes     `json:"hashes,omitempty" bson:"hashes,omitempty"`
	Timestamps FileTimestamps `json:"timestamps,omitempty" bson:"timestamps,omitempty"`
	Ownership  FileOwnership  `json:"ownership,omitempty" bson:"ownership,omitempty"`
	Attributes FileAttributes `json:"attributes,omitempty" bson:"attributes,omitempty"`
}

type CommPID struct {
	Comm string `json:"comm,omitempty" bson:"comm,omitempty"`
	PID  uint32 `json:"pid,omitempty" bson:"pid,omitempty"`
}

// MarshalText implements encoding.TextMarshaler
func (c CommPID) MarshalText() ([]byte, error) {
	s := fmt.Sprintf("%s%s%d", c.Comm, sep, c.PID)
	return []byte(s), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (c *CommPID) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), sep)
	if len(parts) != 2 {
		return fmt.Errorf("invalid CommPID representation: %q", text)
	}
	var err error
	c.Comm = parts[0]
	u64, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid PID value %q: %w", parts[1], err)
	}
	c.PID = uint32(u64)
	return nil
}

type Process struct {
	PID        uint32    `json:"pid,omitempty" bson:"pid,omitempty"`
	Cmdline    string    `json:"cmdline,omitempty" bson:"cmdline,omitempty"`
	Comm       string    `json:"comm,omitempty" bson:"comm,omitempty"`
	PPID       uint32    `json:"ppid,omitempty" bson:"ppid,omitempty"`
	Pcomm      string    `json:"pcomm,omitempty" bson:"pcomm,omitempty"`
	Hardlink   string    `json:"hardlink,omitempty" bson:"hardlink,omitempty"`
	Uid        *uint32   `json:"uid,omitempty" bson:"uid,omitempty"`
	Gid        *uint32   `json:"gid,omitempty" bson:"gid,omitempty"`
	UserName   string    `json:"userName,omitempty" bson:"userName,omitempty"`
	GroupName  string    `json:"groupName,omitempty" bson:"groupName,omitempty"`
	StartTime  time.Time `json:"startTime,omitempty" bson:"startTime,omitempty"`
	UpperLayer *bool     `json:"upperLayer,omitempty" bson:"upperLayer,omitempty"`
	Cwd        string    `json:"cwd,omitempty" bson:"cwd,omitempty"`
	Path       string    `json:"path,omitempty" bson:"path,omitempty"`
	// Deprecated: Use ChildrenMap instead
	Children    []Process            `json:"children,omitempty" bson:"children,omitempty"`
	ChildrenMap map[CommPID]*Process `json:"childrenMap,omitempty" bson:"childrenMap,omitempty"`
}

// MigrateToMap migrates the Children slice to ChildrenMap to accommodate for older versions of the Process struct
func (p *Process) MigrateToMap() {
	if p.ChildrenMap == nil {
		p.ChildrenMap = make(map[CommPID]*Process)
	}
	if len(p.Children) > 0 {
		for i := range p.Children {
			p.Children[i].MigrateToMap()
			commPID := CommPID{Comm: p.Children[i].Comm, PID: p.Children[i].PID}
			p.ChildrenMap[commPID] = &p.Children[i]
		}
		p.Children = nil
	}
}

type FileHashes struct {
	MD5    string `json:"md5,omitempty" bson:"md5,omitempty"`
	SHA1   string `json:"sha1,omitempty" bson:"sha1,omitempty"`
	SHA256 string `json:"sha256,omitempty" bson:"sha256,omitempty"`
}

type FileTimestamps struct {
	CreationTime     time.Time `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	ModificationTime time.Time `json:"modificationTime,omitempty" bson:"modificationTime,omitempty"`
	AccessTime       time.Time `json:"accessTime,omitempty" bson:"accessTime,omitempty"`
}

type FileOwnership struct {
	Uid       *uint32 `json:"uid,omitempty" bson:"uid,omitempty"`
	Gid       *uint32 `json:"gid,omitempty" bson:"gid,omitempty"`
	UserName  string  `json:"userName,omitempty" bson:"userName,omitempty"`
	GroupName string  `json:"groupName,omitempty" bson:"groupName,omitempty"`
}

type FileAttributes struct {
	Permissions string `json:"permissions,omitempty" bson:"permissions,omitempty"`
}
