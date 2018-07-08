package controller

import (
	"flag"
	"fmt"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	//"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//"kuber-controller/client"
	//"kuber-controller/crd"
	//"kuber-controller/metrics"
	"time"
	"kuber/client"
	"kuber/crd"
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

func GetApiExtensionClient2() *client.Crdclient {
	//TODO: replace config with --kubeconfig parameter
	kubeconf := flag.String("kubeconf", "/Users/gurusreekanthc/.kube/config", "path to Kubernetes config file")
	flag.Parse()

	config, err := GetClientConfig(*kubeconf)
	if err != nil {
		panic(err.Error())
	}

	// create clientset and create our CRD, this only need to run once
	_, err = apiextcs.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a new clientset which include our CRD schema
	crdcs, scheme, err := crd.NewClient(config)
	if err != nil {
		panic(err)
	}

	// Create a CRD client interface
	crdclient := client.CrdClient(crdcs, scheme, "default")

	return crdclient
}

func ListCrdInstances(crdclient *client.Crdclient) {
	// List all Example objects
	items, err := crdclient.List(meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("List:\n%s\n", items)
}

func GetCrdByName(crdclient *client.Crdclient, groupName string, groupType string) *crd.Group {
	group, err := crdclient.Get(groupName)

	if err == nil {
		return group
	} else if apierrors.IsNotFound(err) {
		// create group if not exist
		//return CreateCRDInstance(crdclient, groupName, groupType)
		return nil
	} else {
		panic(err)
	}
}