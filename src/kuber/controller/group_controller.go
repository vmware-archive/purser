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
	"kuber/client"
	"kuber/crd"
	"os"
)

// return rest config, if path not specified assume in cluster config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func GetApiExtensionClient() *client.Crdclient {
	kubeconf := flag.String("kubeconf", os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), "path to Kubernetes config file")
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

func GetCrdByName(crdclient *client.Crdclient, groupName string) *crd.Group {
	group, err := crdclient.Get(groupName)

	if err == nil {
		return group
	} else if apierrors.IsNotFound(err) {
		return nil
	} else {
		panic(err)
	}
}

func PrintGroup(group *crd.Group) {
	fmt.Printf("%-25s%s\n", "Group Name:", group.Name)
	fmt.Printf("%-25s\n", "Resources:")
	fmt.Printf("             %-25s%s\n", "Cpu Limit:", group.Spec.AllocatedResources.CpuLimit)
	fmt.Printf("             %-25s%s\n", "Memory Limit:", group.Spec.AllocatedResources.MemoryLimit)
	fmt.Printf("             %-25s%s\n", "Cpu Request:", group.Spec.AllocatedResources.CpuRequest)
	fmt.Printf("             %-25s%s\n", "Memory Request:", group.Spec.AllocatedResources.MemoryRequest)
}