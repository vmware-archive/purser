package config

import (
	"sync"

	"github.com/vmware/purser/pkg/client"
	group_client "github.com/vmware/purser/pkg/client/clientset/typed/groups/v1"
	subscriber_client "github.com/vmware/purser/pkg/client/clientset/typed/subscriber/v1"
	"github.com/vmware/purser/pkg/controller"
	"github.com/vmware/purser/pkg/controller/buffering"
	"github.com/vmware/purser/pkg/controller/utils"
)

// Setup initialzes the controller configuration
func Setup(conf *controller.Config) {
	// initialize client for api extension server
	conf.Kubeclient = utils.GetKubeclient("")
	conf.Resource = controller.Resource{
		Pod:                   true,
		Node:                  true,
		PersistentVolume:      true,
		PersistentVolumeClaim: true,
		ReplicaSet:            true,
		Deployment:            true,
		StatefulSet:           true,
		DaemonSet:             true,
		Job:                   true,
		Service:               true,
	}
	conf.RingBuffer = &buffering.RingBuffer{Size: buffering.BufferSize, Mutex: &sync.Mutex{}}
	clientset, clusterConfig := client.GetAPIExtensionClient()
	conf.Groupcrdclient = group_client.NewGroupClient(clientset, clusterConfig)
	conf.Subscriberclient = subscriber_client.NewSubscriberClient(clientset, clusterConfig)
}
