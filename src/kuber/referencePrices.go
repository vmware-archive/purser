package main

import (
	"fmt"
)

var refcost = map[string]float64{}

func init() {
	// TODO: These are sample reference prices, prices changes based on region and OS type.
	// create an exhaustive list considering region and OS type.

	// memory optimized
	refcost["x1.16xlarge"] = 6.669
	refcost["x1.32xlarge"] = 13.338
	refcost["r4.large"] = 0.133
	refcost["r4.xlarge"] = 0.266
	refcost["r4.2xlarge"] = 0.532
	refcost["r4.4xlarge"] = 1.064
	refcost["r4.8xlarge"] = 2.128
	refcost["r4.16xlarge"] = 4.256

	// compute optimized
	refcost["c5.large"] = 0.085
	refcost["c5.xlarge"] = 0.17
	refcost["c5.2xlarge"] = 0.34
	refcost["c5.4xlarge"] = 0.68
	refcost["c5.9xlarge"] = 1.53
	refcost["c5.18xlarge"] = 3.06
	refcost["c5d.large"] = 0.096
	refcost["c5d.xlarge"] = 0.192
	refcost["c5d.2xlarge"] = 0.384
	refcost["c5d.4xlarge"] = 0.768
	refcost["c5d.9xlarge"] = 1.728
	refcost["c5d.18xlarge"] = 3.456
	refcost["c4.large"] = 0.1
	refcost["c4.xlarge"] = 0.199
	refcost["c4.2xlarge"] = 0.398
	refcost["c4.4xlarge"] = 0.796
	refcost["c4.8xlarge"] = 1.591
}

func getPriceForInstanceType(instanceType string) float64 {
	cost := refcost[instanceType]
	if cost == 0.0 {
		fmt.Printf("Price is not present for instance type - %s, returning default value...\n", instanceType)
		return 0.1
	}
	return cost
}
