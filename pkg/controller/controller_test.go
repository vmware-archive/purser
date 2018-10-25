package controller

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/client"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	subscriber_v1 "github.com/vmware/purser/pkg/client/clientset/typed/subscriber/v1"
)

// TestCrdFlow executes the CRD flow.
func TestCrdFlow(t *testing.T) {
	clientset, clusterConfig := client.GetAPIExtensionClient("")
	subcrdclient := subscriber_v1.NewSubscriberClient(clientset, clusterConfig)
	ListSubscriberCrdInstances(subcrdclient)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

// ListSubscriberCrdInstances fetches list of subscriber CRD instances.
func ListSubscriberCrdInstances(crdclient *subscriber_v1.SubscriberClient) {
	items, err := crdclient.ListSubscriber(meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("List:\n%v\n", items)
}
