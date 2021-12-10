package run

// config represents the confile file
type config struct {
	JsonRpcUrl    string            `json:"json_rpc_url,omitempty"`
	WebsocketUrl  string            `json:"websocket_url,omitempty"`
	KeypairPath   string            `json:"keypair_path,omitempty"`
	AddressLabels map[string]string `json:"address_labels,omitempty"`
	Commitment    string            `json:"commitment,omitempty"`
}
