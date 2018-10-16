package client

import (
	"flag"
	"log"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

const environment = "dev"

// GetClientConfig returns rest config, if path not specified assume in cluster config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	log.Println("Using In cluster config.")
	return rest.InClusterConfig()
}

// GetAPIExtensionClient returns an instance of CRD client.
func GetAPIExtensionClient() (*apiextcs.Clientset, *rest.Config) {
	var config *rest.Config
	var err error

	if environment == "dev" {
		kubeconf := flag.String("kubeconf", "/Users/hkatyal/Downloads/project/upgrade-config-1.9", "path to Kubernetes config file")
		flag.Parse()
		config, err = GetClientConfig(*kubeconf)
	} else {
		config, err = GetClientConfig("")
	}

	if err != nil {
		log.Println(err)
		panic(err.Error())
	}

	// create clientset and create our CRD, this only need to run once
	clientset, err := apiextcs.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, config
}
