package config

import (
	"kuber-controller/buffering"
	"kuber-controller/client"
)

// Resource contains resource configuration
type Resource struct {
	Pod        bool `json:"po"`
	Node       bool `json:"node"`
	Services   bool `json:"services"`
	ReplicaSet bool `json:"replicaset"`
	Deployment bool `json:"deployment"`
	Job        bool `json:"job"`
}

type Config struct {
	Resource Resource `json:"resource"`
	RingBuffer *buffering.RingBuffer
	Groupcrdclient *client.GroupCrdClient
	Subscriberclient *client.SubscriberCrdClient
}
