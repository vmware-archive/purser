package config

// Resource contains resource configuration
type Resource struct {
	Deployment            bool `json:"deployment"`
	ReplicationController bool `json:"rc"`
	ReplicaSet            bool `json:"rs"`
	DaemonSet             bool `json:"ds"`
	Services              bool `json:"svc"`
	Pod                   bool `json:"po"`
	Job                   bool `json:"job"`
	PersistentVolume      bool `json:"pv"`
	Namespace             bool `json:"ns"`
	Secret                bool `json:"secret"`
	Ingress               bool `json:"ing"`
}

// Config struct contains kubewatch configuration
type Config struct {
	//Handler Handler `json:"handler"`
	//Reason   []string `json:"reason"`
	Resource Resource `json:"resource"`
}
