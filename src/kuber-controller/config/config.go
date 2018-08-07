package config

import (
	"kuber-controller/buffering"
	"kuber-controller/client"
)

// Resource contains resource configuration
type Resource struct {
	Pod                   bool `json:"po"`
	Node                  bool `json:"node"`
	PersistentVolume      bool `json:"pv"`
	PersistentVolumeClaim bool `json:"pvc"`
	Service               bool `json:"service"`
	ReplicaSet            bool `json:"replicaset"`
	StatefulSet           bool `json:"statefulset"`
	Deployment            bool `json:"deployment"`
	Job                   bool `json:"job"`
	DaemonSet             bool `json:"daemonset"`
}

type Config struct {
	Resource         Resource `json:"resource"`
	RingBuffer       *buffering.RingBuffer
	Groupcrdclient   *client.GroupCrdClient
	Subscriberclient *client.SubscriberCrdClient
}
