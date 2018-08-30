/*
 * Copyright (c) 2018 VMware Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package purser_plugin

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	cost               *Cost
}

func getNodeDetailsFromClient(nodeName string) *Node {
	node, err := ClientSetInstance.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Node %s not found\n", nodeName)
		return nil
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting Node %s : %v\n", nodeName, statusError.ErrStatus.Message)
		return nil
	} else if err != nil {
		panic(err.Error())
	} else {
		var n Node
		n.name = node.GetObjectMeta().GetName()
		n.instanceType = node.GetObjectMeta().GetLabels()["beta.kubernetes.io/instance-type"]
		//fmt.Println(n)
		return &n
	}
}

func getAllNodeDetailsFromClient() []*Node {
	nodes, err := ClientSetInstance.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	//fmt.Printf("Total nodes = %d\n", len(nodes.Items))

	allNodes := []*Node{}
	i := 0
	for i < len(nodes.Items) {
		node := nodes.Items[i]
		n := Node{}
		n.name = node.GetObjectMeta().GetName()
		n.instanceType = node.GetObjectMeta().GetLabels()["beta.kubernetes.io/instance-type"]
		allNodes = append(allNodes, &n)
		i++
	}
	return allNodes
}

func GetClusterNodes() []v1.Node {
	nodes, err := ClientSetInstance.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return nodes.Items
}

func printNodeDetails(nodes []*Node) {
	fmt.Printf("%-40s%-20s%-30s\n", "NODE NAME", "INSTANCE TYPE", "TOTAL COST")
	for _, value := range nodes {
		fmt.Printf("%-40s%-20s%f$\n", value.name, value.instanceType, value.cost.totalCost)
	}
}

func getNodeDetailsFromNodeDescribe(nodeName string) *Node {
	command := fmt.Sprintf(nodeDescribeCommand, os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), nodeName)
	bytes := executeCommand(command)
	return parseNodeDescribe(bytes)
}

func collectNodes(nodes map[string]*Node) map[string]*Node {
	for key := range nodes {
		node := getNodeDetailsFromNodeDescribe(key)
		nodes[key] = node
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

func parseNodeDescribe(bytes []byte) *Node {
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
	return &node
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

func GetNodeType(node v1.Node) string {
	labels := node.Labels
	return labels["beta.kubernetes.io/instance-type"]
}
