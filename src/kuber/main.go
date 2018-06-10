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
			fmt.Println("Work In Progress...")
		} else if inputs[2] == "node" {
			fmt.Println("Work In Progress...")
		} else {
			printHelp()
		}
	} else {
		printHelp()
	}
}

func main2() {
	//collectPersistentVolume("pvc-22197ba2-6a10-11e8-9bc2-0270c9080a70")
	//collectPersistentVolumeClaim("vrbc-adapter-volume-1-1-569-vrbc-adapter-statefulset-1-1-569-2")
	//getPodsForLabelThroughClient("app=vrbc-transformer")
	//pods := getPodsForLabelThroughClient("app=vrbc-adapter")
	//printPodsVerbose(pods)
	getNodeDetailsFromClient("ip-172-20-34-236.ec2.internal")
}

func init() {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", "/Users/gurusreekanthc/staging-config-1.9", "/Users/gurusreekanthc/staging-config-1.9")
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
