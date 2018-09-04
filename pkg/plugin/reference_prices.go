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

package plugin

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

	// storage optimized
	computeRefCost["i3.large"] = 0.156
	computeRefCost["i3.xlarge"] = 0.312
	computeRefCost["i3.2xlarge"] = 0.624
	computeRefCost["i3.4xlarge"] = 1.248
	computeRefCost["i3.8xlarge"] = 2.496
	computeRefCost["i3.16xlarge"] = 4.992
	computeRefCost["i3.metal"] = 4.992
	computeRefCost["h1.2xlarge"] = 0.468
	computeRefCost["h1.4xlarge"] = 0.936
	computeRefCost["h1.8xlarge"] = 1.872
	computeRefCost["h1.16xlarge"] = 3.744
	computeRefCost["d2.xlarge"] = 0.69
	computeRefCost["d2.2xlarge"] = 1.38
	computeRefCost["d2.4xlarge"] = 2.76
	computeRefCost["d2.8xlarge"] = 5.52

	// general purpose
	computeRefCost["t2.nano"] = 0.0058
	computeRefCost["t2.micro"] = 0.0116
	computeRefCost["t2.small"] = 0.023
	computeRefCost["t2.medium"] = 0.0464
	computeRefCost["t2.large"] = 0.0928
	computeRefCost["t2.xlarge"] = 0.1856
	computeRefCost["t2.2xlarge"] = 0.3712
	computeRefCost["m5.large"] = 0.096
	computeRefCost["m5.xlarge"] = 0.192
	computeRefCost["m5.2xlarge"] = 0.384
	computeRefCost["m5.4xlarge"] = 0.768
	computeRefCost["m5.12xlarge"] = 2.304
	computeRefCost["m5.24xlarge"] = 4.608
	computeRefCost["m5d.large"] = 0.113
	computeRefCost["m5d.xlarge"] = 0.226
	computeRefCost["m5d.2xlarge"] = 0.452
	computeRefCost["m5d.4xlarge"] = 0.904
	computeRefCost["m5d.12xlarge"] = 2.712
	computeRefCost["m5d.24xlarge"] = 5.424
	computeRefCost["m4.large"] = 0.1
	computeRefCost["m4.xlarge"] = 0.2
	computeRefCost["m4.2xlarge"] = 0.4
	computeRefCost["m4.4xlarge"] = 0.8
	computeRefCost["m4.10xlarge"] = 2
	computeRefCost["m4.16xlarge"] = 3.2

	// new vals
	computeRefCost["m3.large"] = 0.1
	computeRefCost["m3.xlarge"] = 0.2
	computeRefCost["m3.2xlarge"] = 0.4
	computeRefCost["m3.4xlarge"] = 0.8
	computeRefCost["m3.10xlarge"] = 2
	computeRefCost["m3.16xlarge"] = 3.2

	// gpu
	computeRefCost["p2.xlarge"] = 0.9
	computeRefCost["p2.8xlarge"] = 7.2
	computeRefCost["p2.16xlarge"] = 14.4
	computeRefCost["p3.2xlarge"] = 3.06
	computeRefCost["p3.8xlarge"] = 12.24
	computeRefCost["p3.16xlarge"] = 24.48
	computeRefCost["g3.4xlarge"] = 1.14
	computeRefCost["g3.8xlarge"] = 2.28
	computeRefCost["g3.16xlarge"] = 4.56

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
		return 0.1
	}
	return cost
}
