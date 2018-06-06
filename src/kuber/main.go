package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	getPodsCommand      = "kubectl --kubeconfig=%s get pods -l %s -o json"
	getNodeCommand      = "kubectl --kubeconfig=%s get node %s -o json"
	nodeDescribeCommand = "kubectl --kubeconfig=%s describe node %s"
)

func getNodeDetails(nodeName string) Node {
	node := Node{}
	command := fmt.Sprintf(getNodeCommand, os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), nodeName)
	bytes := executeCommand(command)
	json := string(bytes)
	node.name = nodeName
	node.instanceType = gjson.Get(json, "metadata.labels.beta\\.kubernetes\\.io/instance-type").Str
	return node
}

func getPodsForLabel(label string) []Pod {
	pods := []Pod{}
	command := fmt.Sprintf(getPodsCommand, os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), label)
	bytes := executeCommand(command)
	json := string(bytes)
	items := gjson.Get(json, "items")

	items.ForEach(func(key, value gjson.Result) bool {
		name := value.Get("metadata.name")
		nodeName := value.Get("spec.nodeName")
		pod := Pod{name: name.Str, nodeName: nodeName.Str}
		pods = append(pods, pod)
		return true
	})
	return pods
}

func executeCommand(command string) []byte {
	slice := strings.Fields(command)
	cmd := exec.Command(slice[0], slice[1:]...)
	cmd.Env = os.Environ()
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return out.Bytes()
}

func printPodDetails(pods []Pod) {
	fmt.Println("===POD Details===")
	fmt.Println("POD Name \t\t\t\t\t Node Name")
	for _, value := range pods {
		fmt.Println(value.name + " \t" + value.nodeName)
	}
}

func printNodeDetails(nodes []Node) {
	fmt.Println("===Node Details===")
	fmt.Println("Node Name \t\t\t InstanceType")
	for _, value := range nodes {
		fmt.Println(value.name + " \t" + value.instanceType)
	}
}

func getNodeDetailsFromNodeDescribe(nodeName string) Node {
	command := fmt.Sprintf(nodeDescribeCommand, os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), nodeName)
	bytes := executeCommand(command)
	return parseNodeDescribe(bytes)
}

func collectNodes(nodes map[string]*Node) map[string]*Node {
	for key := range nodes {
		node := getNodeDetailsFromNodeDescribe(key)
		nodes[key] = &node
	}
	return nodes
}

func calculateCost(pods []Pod, nodes map[string]*Node) []Pod {
	i := 0
	for i <= len(pods)-1 {
		node := nodes[pods[i].nodeName]
		pods[i].nodeCostPercentage = (float64)(node.getPodResourcePercentage(pods[i].name))
		totalCost, cpuCost, memoryCost := 10.0, 3.0, 7.0
		podCost := Cost{}
		podCost.totalCost = pods[i].nodeCostPercentage * totalCost
		podCost.cpuCost = pods[i].nodeCostPercentage * cpuCost
		podCost.memoryCost = pods[i].nodeCostPercentage * memoryCost
		pods[i].cost = podCost
		i++
	}
	return pods
}

func printPodsVerbose(pods []Pod) {
	i := 0
	fmt.Printf("==Pods Cost Details==\n")
	for i <= len(pods)-1 {
		fmt.Printf("%-25s%s\n", "Pod Name:", pods[i].name)
		fmt.Printf("%-25s%s\n", "Node:", pods[i].nodeName)
		fmt.Printf("%-25s%.2f\n", "Pod Cost Percentage:", pods[i].nodeCostPercentage*100.0)
		fmt.Printf("%-25s\n", "Cost:")
		fmt.Printf("    %-21s%f$\n", "Total Cost:", pods[i].cost.totalCost)
		fmt.Printf("    %-21s%f$\n", "CPU Cost:", pods[i].cost.cpuCost)
		fmt.Printf("    %-21s%f$\n", "Memory Cost:", pods[i].cost.memoryCost)
		fmt.Printf("\n")
		i++
	}
}

func getPodsCostForLabel(label string) {
	pods := getPodsForLabel(label)
	nodes := map[string]*Node{}
	for _, val := range pods {
		nodes[val.nodeName] = nil
	}
	nodes = collectNodes(nodes)
	pods = calculateCost(pods, nodes)
	printPodsVerbose(pods)
}

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

func printHelp() {
	fmt.Printf("Try one of the following commands...\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost label <key=val>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost pod <pod name>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost node <node name>\n")
}
