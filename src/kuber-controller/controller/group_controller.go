package controller

import (
	"flag"
	"fmt"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"kuber-controller/client"
	"kuber-controller/crd"
	"kuber-controller/metrics"
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

func CreateCRDInstance(crdclient *client.Crdclient, groupName string, groupType string) *crd.Group {
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
	return result
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
		return CreateCRDInstance(crdclient, groupName, groupType)
	} else {
		panic(err)
	}
}

func UpdateNamespaceGroupCrd(crdclient *client.Crdclient, groupName string, groupType string, pod string,
	metric *metrics.Metrics) {

	group := GetCrdByName(crdclient, groupName, groupType)
	existingPods := group.Spec.PodsMetrics

	if existingPods == nil {
		existingPods = map[string]*metrics.Metrics{}
	}

	existingPods[pod] = metric
	group.Spec.PodsMetrics = existingPods
	group.Spec.AllocatedResources = calculatedAggregatedPodMetric(existingPods)
	group.Name = groupName

	//fmt.Println(group)
	_, err := crdclient.Update(group)

	if err != nil {
		fmt.Printf("There is a panic while updating the crd for group = %s\n", groupName)
		panic(err)
	} else {
		fmt.Printf("Updating the crd for group = %s is successful\n", groupName)
	}
}

func calculatedAggregatedPodMetric(met map[string]*metrics.Metrics) *metrics.Metrics {
	cpuLimit := &resource.Quantity{}
	memoryLimit := &resource.Quantity{}
	cpuRequest := &resource.Quantity{}
	memoryRequest := &resource.Quantity{}
	for _, c := range met {
		cpuLimit.Add(*c.CpuLimit)
		memoryLimit.Add(*c.MemoryLimit)
		cpuRequest.Add(*c.CpuRequest)
		memoryRequest.Add(*c.MemoryRequest)
	}
	return &metrics.Metrics{
		CpuLimit:      cpuLimit,
		MemoryLimit:   memoryLimit,
		CpuRequest:    cpuRequest,
		MemoryRequest: memoryRequest,
	}
}
