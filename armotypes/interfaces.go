package armotypes

type VolumeInfo interface {
	GetProvider() string
	GetAccountId() string
	GetInstanceId() string
	GetVolumeId() string
	GetVolumeScanId() string
}
