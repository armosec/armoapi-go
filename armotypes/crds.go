package armotypes

type GenericCRD[T any] struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Spec       T      `json:"spec"`
}
