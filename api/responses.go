package api

// ServiceDiscoveryResponse is returned when well known file is requested by a Paymail client
type ServiceDiscoveryResponse struct {
	Version      string                 `json:"bsvalias"`
	Capabilities map[string]interface{} `json:"capabilities"`
}

// PKIResponse is returned when a Paymail id is requested
type PKIResponse struct {
	Version string `json:"bsvalias"`
	Handle  string `json:"handle"`
	PubKey  string `json:"pubkey"`
}
