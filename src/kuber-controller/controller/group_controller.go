package controller

import (
	"flag"
	"fmt"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"kuber-controller/client"
	"kuber-controller/crd"
	"time"
)

// return rest config, if path not specified assume in cluster config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func GetApiExtensionClient() *client.Crdclient {
	//TODO: replace config with --kubeconfig parameter
	kubeconf := flag.String("kubeconf", "/Users/gurusreekanthc/.kube/config", "path to Kubernetes config file")
	flag.Parse()

	config, err := GetClientConfig(*kubeconf)
	if err != nil {
		panic(err.Error())
	}

	// create clientset and create our CRD, this only need to run once
	clientset, err := apiextcs.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// note: if the CRD exist our CreateCRD function is set to exit without an error
	err = crd.CreateCRD(clientset)
	if err != nil {
		panic(err)
	}

	// Wait for the CRD to be created before we use it (only needed if its a new one)
	time.Sleep(3 * time.Second)

	// Create a new clientset which include our CRD schema
	crdcs, scheme, err := crd.NewClient(config)
	if err != nil {
		panic(err)
	}

	// Create a CRD client interface
	crdclient := client.CrdClient(crdcs, scheme, "default")

	return crdclient
}

func CreateCRDInstance(crdclient *client.Crdclient, groupName string, groupType string) {
	// Create a new Example object and write to k8s
	example := &crd.Group{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: groupName,
			//Labels: map[string]string{"mylabel": "test"},
		},
		Spec: crd.GroupSpec{
			Name: groupName,
			Type: groupType,
		},
		Status: crd.GroupStatus{
			State:   "created",
			Message: "Done",
		},
	}

	result, err := crdclient.Create(example)
	if err == nil {
		fmt.Printf("CREATED: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else {
		panic(err)
	}
}

func ListCrdInstances(crdclient *client.Crdclient) {
	// List all Example objects
	items, err := crdclient.List(meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("List:\n%s\n", items)
}
