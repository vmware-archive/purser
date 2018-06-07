package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Metric details
type Metric struct {
	cpuRequest    float64
	cpuLimit      float64
	memoryRequest float64
	memoryLimit   float64
}

func (node *Node) getPodResourcePercentage(pod string) float64 {
	podMetrics := node.podsResources[pod]
	if podMetrics == nil {
		return 0.0
	}
	return podMetrics.cpuRequest / (float64)(node.allocatedResources.cpuRequest)
}

// Node Information
type Node struct {
	name               string
	instanceType       string
	allocatedResources *Metric
	podsResources      map[string]*Metric
}

func getNodeDetails(nodeName string) Node {
	node := Node{}
	command := fmt.Sprintf(getNodeByName, os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), nodeName)
	bytes := executeCommand(command)
	json := string(bytes)
	node.name = nodeName
	node.instanceType = gjson.Get(json, "metadata.labels.beta\\.kubernetes\\.io/instance-type").Str
	return node
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

type state int

const (
	begin state = 1 + iota
	labelStart
	collectLabels
	podsStatus
	collectPodsMetrics
	collectAllocatedResources
	endOfCollection
)

func parseNodeDescribe(bytes []byte) Node {
	input := string(bytes)
	lines := strings.Split(input, "\n")
	node := Node{}

	flag := begin
	i := 0
	length := len(lines)
	podsResources := map[string]*Metric{}
	node.podsResources = podsResources
	for i < length {
		val := lines[i]
		if flag == begin {
			if strings.HasPrefix(val, "Name:") {
				words := strings.Fields(val)
				node.name = words[1]
				flag = labelStart
			}
		} else if flag == labelStart {
			if strings.HasPrefix(val, "Labels:") {
				flag = collectLabels
			}
		} else if flag == collectLabels {
			if strings.HasPrefix(val, "Annotations:") {
				flag = podsStatus
			} else if strings.Contains(val, "beta.kubernetes.io/instance-type") {
				words := strings.Split(val, "=")
				node.instanceType = words[1]
				flag = podsStatus
			}
		} else if flag == podsStatus {
			if strings.HasPrefix(val, "Non-terminated Pods:") {
				i = i + 2
				flag = collectPodsMetrics
			}
		} else if flag == collectPodsMetrics {
			if strings.HasPrefix(val, "Allocated resources:") {
				flag = collectAllocatedResources
				i = i + 3
			} else {
				words := strings.Fields(val)
				metric := Metric{}
				metric.cpuRequest = convertToMillis(words[2])
				metric.cpuLimit = convertToMillis(words[4])
				metric.memoryRequest = convertToMi(words[6])
				metric.memoryLimit = convertToMi(words[8])
				podsResources[words[1]] = &metric
			}
		} else if flag == collectAllocatedResources {
			words := strings.Fields(val)
			metric := Metric{}
			metric.cpuRequest = convertToMillis(words[0])
			metric.cpuLimit = convertToMillis(words[2])
			metric.memoryRequest = convertToMi(words[4])
			metric.memoryLimit = convertToMi(words[6])
			node.allocatedResources = &metric
			flag = endOfCollection
		}
		i++
	}
	return node
}

func convertToMillis(input string) float64 {
	number := input
	if strings.HasSuffix(input, "m") {
		number = input[:len(input)-2]
	}
	s, err := strconv.ParseFloat(number, 64)
	if err != nil {
		fmt.Println(err)
	}
	return s
}

func convertToMi(input string) float64 {
	divisor := 1024.0 * 1024.0
	number := input
	if strings.HasSuffix(input, "Mi") {
		number = input[:len(input)-2]
		divisor = 1.0
	} else if strings.HasSuffix(input, "Gi") {
		number = input[:len(input)-2]
		divisor = 1 / (1024.0)
	} else if strings.HasSuffix(input, "Ki") {
		number = input[:len(input)-2]
		divisor = 1024.0
	}
	s, err := strconv.ParseFloat(number, 64)
	if err != nil {
		fmt.Println(err)
	}
	return s / divisor
}
