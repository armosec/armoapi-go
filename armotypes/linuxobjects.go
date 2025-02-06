package armotypes

import "time"

type File struct {
	Path       string         `json:"path,omitempty" bson:"path,omitempty"`
	Size       int64          `json:"size,omitempty" bson:"size,omitempty"`
	Hashes     FileHashes     `json:"hashes,omitempty" bson:"hashes,omitempty"`
	Timestamps FileTimestamps `json:"timestamps,omitempty" bson:"timestamps,omitempty"`
	Ownership  FileOwnership  `json:"ownership,omitempty" bson:"ownership,omitempty"`
	Attributes FileAttributes `json:"attributes,omitempty" bson:"attributes,omitempty"`
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
	Children   []Process `json:"children,omitempty" bson:"children,omitempty"`
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
