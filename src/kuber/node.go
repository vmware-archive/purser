package main

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
