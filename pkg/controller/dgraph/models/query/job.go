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

// RetrieveJobHierarchy returns hierarchy for a given daemonset
func RetrieveJobHierarchy(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for job, empty name is given")
		return JSONDataWrapper{}
	}
	query := `query {
		parent(func: has(isJob)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~job @filter(has(isPod)) {
				name
				type
			}
		}
	}`
	return getJSONDataFromQuery(query)
}

// RetrieveJobMetrics returns metrics for a given daemonset
func RetrieveJobMetrics(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for job, empty name is given")
		return JSONDataWrapper{}
	}
	query := `query {
		parent(func: has(isJob)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~job @filter(has(isPod)) {
				name
				type
				cpu: podCpu as cpuRequest
				memory: podMemory as memoryRequest
				storage: pvcStorage as storageRequest
				cpuCost: math(podCpu * ` + defaultCPUCostPerCPUPerHour + `)
				memoryCost: math(podMemory * ` + defaultMemCostPerGBPerHour + `)
				storageCost: math(pvcStorage * ` + defaultStorageCostPerGBPerHour + `)
			}
			cpu: cpu as sum(val(podCpu))
			memory: memory as sum(val(podMemory))
			storage: storage as sum(val(pvcStorage))
			cpuCost: math(cpu * ` + defaultCPUCostPerCPUPerHour + `)
			memoryCost: math(memory * ` + defaultMemCostPerGBPerHour + `)
			storageCost: math(storage * ` + defaultStorageCostPerGBPerHour + `)
		}
	}`
	return getJSONDataFromQuery(query)
}
