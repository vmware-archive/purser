package main

type Metric struct {
	cpuRequest    string
	cpuLimit      string
	memoryRequest string
	memoryLimit   string
}

// Node Information
type Node struct {
	name               string
	instanceType       string
	allocatedResources *Metric
	podsResources      map[string]*Metric
}
