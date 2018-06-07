package main

import (
	"fmt"
)

var computeRefCost = map[string]float64{}

var storageRefCost = map[string]float64{}

func init() {
	// TODO: These are sample reference prices, prices changes based on region and OS type.
	// create an exhaustive list considering region and OS type.

	// memory optimized
	computeRefCost["x1.16xlarge"] = 6.669
	computeRefCost["x1.32xlarge"] = 13.338
	computeRefCost["r4.large"] = 0.133
	computeRefCost["r4.xlarge"] = 0.266
	computeRefCost["r4.2xlarge"] = 0.532
	computeRefCost["r4.4xlarge"] = 1.064
	computeRefCost["r4.8xlarge"] = 2.128
	computeRefCost["r4.16xlarge"] = 4.256

	// compute optimized
	computeRefCost["c5.large"] = 0.085
	computeRefCost["c5.xlarge"] = 0.17
	computeRefCost["c5.2xlarge"] = 0.34
	computeRefCost["c5.4xlarge"] = 0.68
	computeRefCost["c5.9xlarge"] = 1.53
	computeRefCost["c5.18xlarge"] = 3.06
	computeRefCost["c5d.large"] = 0.096
	computeRefCost["c5d.xlarge"] = 0.192
	computeRefCost["c5d.2xlarge"] = 0.384
	computeRefCost["c5d.4xlarge"] = 0.768
	computeRefCost["c5d.9xlarge"] = 1.728
	computeRefCost["c5d.18xlarge"] = 3.456
	computeRefCost["c4.large"] = 0.1
	computeRefCost["c4.xlarge"] = 0.199
	computeRefCost["c4.2xlarge"] = 0.398
	computeRefCost["c4.4xlarge"] = 0.796
	computeRefCost["c4.8xlarge"] = 1.591

	// storage ref costs
	storageRefCost["gp2"] = 0.1
	storageRefCost["io1"] = 0.125
	storageRefCost["st1"] = 0.045
	storageRefCost["sc1"] = 0.025

}

func getPriceForInstanceType(instanceType string) float64 {
	cost := computeRefCost[instanceType]
	if cost == 0.0 {
		fmt.Printf("Price is not present for instance type - %s, returning default value...\n", instanceType)
		return 0.1
	}
	return cost
}

func getPriceForVolumeType(volumeType string) float64 {
	cost := storageRefCost[volumeType]
	if cost == 0.0 {
		fmt.Printf("Price is not present for volume type - %s, returning default value...\n", volumeType)
		return 0.1
	}
	return cost
}
