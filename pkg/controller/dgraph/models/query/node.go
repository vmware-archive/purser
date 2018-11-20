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

package query

import "github.com/Sirupsen/logrus"

// RetrieveNodeHierarchy returns hierarchy for a given node
func RetrieveNodeHierarchy(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for node, empty name is given")
		return JSONDataWrapper{}
	}

	query := `query {
		parent(func: has(isNode)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~node @filter(has(isPod)) {
				name
				type
			}
        }
    }`
	return getJSONDataFromQuery(query)
}

// RetrieveNodeMetrics returns metrics for a given node
func RetrieveNodeMetrics(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for node, empty name is given")
		return JSONDataWrapper{}
	}

	query := `query {
		parent(func: has(isNode)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~node @filter(has(isPod)) {
				name
				type
				cpu: podCpu as cpuRequest
				memory: podMemory as memoryRequest
				storage: podStorage as storageRequest
				cpuCost: math(podCpu * ` + defaultCPUCostPerCPUPerHour + `)
				memoryCost: math(podMemory * ` + defaultMemCostPerGBPerHour + `)
				storageCost: math(podStorage * ` + defaultStorageCostPerGBPerHour + `)
			}
			cpu: cpu as cpuCapacity
			memory: memory as memoryCapacity
			storage: storage as sum(val(podStorage))
			cpuCost: math(cpu * ` + defaultCPUCostPerCPUPerHour + `)
			memoryCost: math(memory * ` + defaultMemCostPerGBPerHour + `)
			storageCost: math(storage * ` + defaultStorageCostPerGBPerHour + `)
		}
	}`
	return getJSONDataFromQuery(query)
}
