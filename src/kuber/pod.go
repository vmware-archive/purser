package main

// Cost details
type Cost struct {
	totalCost  float64
	cpuCost    float64
	memoryCost float64
}

// Pod Information
type Pod struct {
	name               string
	nodeName           string
	nodeCostPercentage float64
	cost               Cost
}
