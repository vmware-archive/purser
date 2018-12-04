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

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/utils"
)

// RetrieveDaemonsetHierarchy returns hierarchy for a given daemonset
func RetrieveDaemonsetHierarchy(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for daemonset, empty name is given")
		return JSONDataWrapper{}
	}
	query := `query {
		parent(func: has(isDaemonset)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~daemonset @filter(has(isPod)) {
				name
				type
			}
		}
	}`
	return getJSONDataFromQuery(query)
}

// RetrieveDaemonsetMetrics returns metrics for a given daemonset
func RetrieveDaemonsetMetrics(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for daemonset, empty name is given")
		return JSONDataWrapper{}
	}
	secondsSinceMonthStart := fmt.Sprintf("%f", utils.GetSecondsSince(utils.GetCurrentMonthStartTime()))
	query := `query {
		parent(func: has(isDaemonset)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~daemonset @filter(has(isPod)) {
				name
				type
				cpu: podCpu as cpuRequest
				memory: podMemory as memoryRequest
				storage: pvcStorage as storageRequest
				st as startTime
				stSeconds as math(since(st))
				secondsSinceStart as math(cond(stSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, stSeconds))
				et as endTime
				isTerminated as count(endTime)
				secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
				durationInHours as math((secondsSinceStart - secondsSinceEnd) / 3600)
				cpuCost: podCpuCost as math(podCpu * durationInHours * ` + defaultCPUCostPerCPUPerHour + `)
				memoryCost: podMemCost as math(podMemory * durationInHours * ` + defaultMemCostPerGBPerHour + `)
				storageCost: pvcStorageCost as math(pvcStorage * durationInHours * ` + defaultStorageCostPerGBPerHour + `)
			}
			cpu: sum(val(podCpu))
			memory: sum(val(podMemory))
			storage: sum(val(pvcStorage))
			cpuCost: sum(val(podCpuCost))
			memoryCost: sum(val(podMemCost))
			storageCost: sum(val(pvcStorageCost))
		}
	}`
	return getJSONDataFromQuery(query)
}
