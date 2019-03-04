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

	"github.com/vmware/purser/pkg/controller/dgraph/models"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/utils"
)

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

	secondsSinceMonthStart := fmt.Sprintf("%f", utils.GetSecondsSince(utils.GetCurrentMonthStartTime()))
	query := `query {
		parent(func: has(isNode)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~node @filter(has(isPod)) {
				name
				type
				cpu: podCpu as cpuRequest
				memory: podMemory as memoryRequest
				storage: pvcStorage as storageRequest
				stChild as startTime
				stSecondsChild as math(since(stChild))
				secondsSinceStartChild as math(cond(stSecondsChild > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, stSecondsChild))
				etChild as endTime
				isTerminatedChild as count(endTime)
				secondsSinceEndChild as math(cond(isTerminatedChild == 0, 0.0, since(etChild)))
				durationInHoursChild as math(cond(secondsSinceStartChild > secondsSinceEndChild, (secondsSinceStartChild - secondsSinceEndChild) / 3600, 0.0))
				podPricePerCPU as cpuPrice
				podPricePerMemory as memoryPrice
				cpuCost: math(podCpu * durationInHoursChild * podPricePerCPU)
				memoryCost: math(podMemory * durationInHoursChild * podPricePerMemory)
				storageCost: math(pvcStorage * durationInHoursChild * ` + models.DefaultStorageCostPerGBPerHour + `)
			}
			cpu: cpu as cpuCapacity
			memory: memory as memoryCapacity
			storage: storage as sum(val(pvcStorage))
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math(cond(secondsSinceStart > secondsSinceEnd, (secondsSinceStart - secondsSinceEnd) / 3600, 0.0))
			pricePerCPU as cpuPrice
			pricePerMemory as memoryPrice
			cpuCost: math(cpu * durationInHours * pricePerCPU)
			memoryCost: math(memory * durationInHours * pricePerMemory)
			storageCost: math(storage * durationInHours * ` + models.DefaultStorageCostPerGBPerHour + `)
		}
	}`
	return getJSONDataFromQuery(query)
}
