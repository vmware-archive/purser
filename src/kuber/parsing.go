package main

import (
	"fmt"
	"strings"
)

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

func parseNodeDescribe(bytes []byte) {
	input := string(bytes)
	lines := strings.Split(input, "\n")
	fmt.Println(len(lines))
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
				metric.cpuRequest = words[2]
				metric.cpuLimit = words[4]
				metric.memoryRequest = words[6]
				metric.memoryLimit = words[8]
				podsResources[words[1]] = &metric
			}
		} else if flag == collectAllocatedResources {
			words := strings.Fields(val)
			metric := Metric{}
			metric.cpuRequest = words[0]
			metric.cpuLimit = words[2]
			metric.memoryRequest = words[4]
			metric.memoryLimit = words[6]
			node.allocatedResources = &metric
			flag = endOfCollection
		}
		i++
	}
	fmt.Println(node)
}
