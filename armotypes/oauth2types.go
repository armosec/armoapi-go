package armotypes

type ConnectedStatus string

const (
	Connected    ConnectedStatus = "connected"
	Disconnected ConnectedStatus = "disconnected"
)

type ProviderConnectionStatus struct {
	Status ConnectedStatus `json:"status"`
}
