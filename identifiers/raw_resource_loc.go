package identifiers

// S3 object path; support in bytes range
type S3ObjectPath struct {
	Bucket string         `json:"bucket"`
	Key    string         `json:"key"`
	Range  *S3ObjectRange `json:"range,omitempty"`
}

type S3ObjectRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}
