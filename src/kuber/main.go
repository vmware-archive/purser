package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"kuber/client"
	"kuber/controller"
	"strings"
)

// ClientSetInstance helps in accessing kubernetes apis through client.
var ClientSetInstance *kubernetes.Clientset
var crdclient  *client.Crdclient

func main1()  {
	getClusterSummary()
}

func main() {
	inputs := os.Args[1:]
	inputs = inputs[1:]
	if len(inputs) == 4 && inputs[0] == "get" && inputs[1] == "cost" {
		if inputs[2] == "label" {
			getPodsCostForLabel(inputs[3])
		} else if inputs[2] == "pod" {
			getPodCost(inputs[3])
		} else if inputs[2] == "node" {
			getAllNodesCost()
		} else {
			printHelp()
		}
	} else if (len(inputs) == 4 && inputs[0] == "get" && inputs[1] == "resources") {
		if inputs[2] == "namespace" {
			group := controller.GetCrdByName(crdclient, inputs[3])
			if group != nil {
				controller.PrintGroup(group)
			} else {
				fmt.Printf("Group %s is not present\n", inputs[3])
			}
		} else if inputs[2] == "label" {
			if (!strings.Contains(inputs[3], "=")) {
				printHelp()
			}
			group := controller.GetCrdByName(crdclient, createGroupNameFromLabel(inputs[3]))
			if group != nil {
				controller.PrintGroup(group)
			} else {
				fmt.Printf("Group %s is not present\n", inputs[3])
			}
		}
	} else if (len(inputs) == 2 && inputs[0] == "get") {
		if inputs[1] == "summary" {
			getClusterSummary()
		} else if inputs[1] == "savings" {
			getSavings()
		} else {
			printHelp()
		}
	} else {
		printHelp()
	}
}

func createGroupNameFromLabel(input string) string {
	inp := strings.Split(input, "=")
	key := inp[0]
	val := inp[1]
	groupName := key + "." + val
	if strings.Contains(groupName, "/") {
		groupName = strings.Replace(groupName, "/", "-", -1)
	}
	groupName = strings.ToLower(groupName)
	return groupName
}

func main2()  {
	//controller.ListCrdInstances(crdclient)
	groupName := "apundlik1"
	group := controller.GetCrdByName(crdclient, groupName)
	//fmt.Println(group)
	if group != nil {
		controller.PrintGroup(group)
	} else {
		fmt.Printf("Group %s is not present\n", groupName)
	}
}

func init2() {
	crdclient = controller.GetApiExtensionClient()
}

func init() {
	var kubeconfig *string
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

	// Crd client
	crdclient = controller.GetApiExtensionClient()
}

func printHelp() {
	fmt.Printf("Try one of the following commands...\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get summary\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get resources namespace <Namespace>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get resources label <key=val>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost label <key=val>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost pod <pod name>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost node <node name>\n")
}
