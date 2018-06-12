package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientSetInstance helps in accessing kubernetes apis through client.
var ClientSetInstance *kubernetes.Clientset

func main() {
	inputs := os.Args[1:]
	inputs = inputs[1:]
	if len(inputs) >= 4 && inputs[0] == "get" && inputs[1] == "cost" {
		if inputs[2] == "label" {
			getPodsCostForLabel(inputs[3])
		} else if inputs[2] == "pod" {
			getPodCost(inputs[3])
		} else if inputs[2] == "node" {
			getAllNodesCost()
			//fmt.Println("Work In Progress...")
		} else {
			printHelp()
		}
	} else {
		printHelp()
	}
}

func init() {
	var kubeconfig *string
	//fmt.Println(os.Environ())
	kubeconfig = flag.String("kubeconfig", os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"))
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	ClientSetInstance = clientset
}

func printHelp() {
	fmt.Printf("Try one of the following commands...\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost label <key=val>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost pod <pod name>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost node <node name>\n")
}
