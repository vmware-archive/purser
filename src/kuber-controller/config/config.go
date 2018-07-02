package config

// Resource contains resource configuration
type Resource struct {
	Pod  bool `json:"po"`
	Node bool `json:"node"`
}

// Config struct contains kubewatch configuration
type Config struct {
	//Handler Handler `json:"handler"`
	//Reason   []string `json:"reason"`
	Resource Resource `json:"resource"`
}
