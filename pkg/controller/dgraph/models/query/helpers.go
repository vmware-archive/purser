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
	"github.com/vmware/purser/pkg/controller/utils"
)

func getSecondsSinceMonthStart() string {
	return fmt.Sprintf("%f", utils.GetSecondsSince(utils.GetCurrentMonthStartTime()))
}

func getQueryForMetricsComputationWithAliasAndVariables(suffix string) string {
	return `name
			type
			cpu: cpu` + suffix + ` as cpuRequest
			memory: memory` + suffix + ` as memoryRequest
			storage: storage` + suffix + ` as storageRequest
			` + getQueryForTimeComputation(suffix) + `
			` + getQueryForCostWithPriceWithAliasAndVariables(suffix)
}

func getQueryForMetricsComputationWithAlias(suffix string) string {
	return `name
			type
			cpu: cpu` + suffix + ` as cpuRequest
			memory: memory` + suffix + ` as memoryRequest
			storage: storage` + suffix + ` as storageRequest
			` + getQueryForTimeComputation(suffix) + `
			` + getQueryForCostWithPriceWithAlias(suffix)
}

func getQueryForMetricsComputation(suffix string) string {
	return `cpu` + suffix + ` as cpuRequest
			memory` + suffix + ` as memoryRequest
			storage` + suffix + ` as storageRequest
			` + getQueryForTimeComputation(suffix) + `
			` + getQueryForCostWithPrice(suffix)
}

func getQueryForTimeComputation(suffix string) string {
	secondsSinceMonthStart := getSecondsSinceMonthStart()
	return `st` + suffix + ` as startTime
			stSeconds` + suffix + ` as math(since(st` + suffix + `))
			secondsSinceStart` + suffix + ` as math(cond(stSeconds` + suffix + ` > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, stSeconds` + suffix + `))
			et` + suffix + ` as endTime
			isTerminated` + suffix + ` as count(endTime)
			secondsSinceEnd` + suffix + ` as math(cond(isTerminated == 0, 0.0, since(et` + suffix + `)))
			durationInHours` + suffix + ` as math(cond(secondsSinceStart` + suffix + ` > secondsSinceEnd` + suffix + `, (secondsSinceStart` + suffix + ` - secondsSinceEnd` + suffix + `) / 3600, 0.0))`
}

func getQueryForCostWithPriceWithAliasAndVariables(suffix string) string {
	return `pricePerCPU` + suffix + ` as cpuPrice
			pricePerMemory` + suffix + ` as memoryPrice
			cpuCost: cpuCost` + suffix + ` as math(cpu` + suffix + ` * durationInHours` + suffix + ` * pricePerCPU` + suffix + `)
			memoryCost: memoryCost` + suffix + ` as math(memory` + suffix + ` * durationInHours` + suffix + ` * pricePerMemory` + suffix + `)
			storageCost: storageCost` + suffix + ` as math(storage` + suffix + ` * durationInHours` + suffix + ` * ` + models.DefaultStorageCostPerGBPerHour + `)`
}

func getQueryForCostWithPriceWithAlias(suffix string) string {
	return `pricePerCPU` + suffix + ` as cpuPrice
			pricePerMemory` + suffix + ` as memoryPrice
			cpuCost: math(cpu` + suffix + ` * durationInHours` + suffix + ` * pricePerCPU` + suffix + `)
			memoryCost: math(memory` + suffix + ` * durationInHours` + suffix + ` * pricePerMemory` + suffix + `)
			storageCost: math(storage` + suffix + ` * durationInHours` + suffix + ` * ` + models.DefaultStorageCostPerGBPerHour + `)`
}

func getQueryForCostWithPrice(suffix string) string {
	return `pricePerCPU` + suffix + ` as cpuPrice
			pricePerMemory` + suffix + ` as memoryPrice
			cpuCost` + suffix + ` as math(cpu` + suffix + ` * durationInHours` + suffix + ` * pricePerCPU` + suffix + `)
			memoryCost` + suffix + ` as math(memory` + suffix + ` * durationInHours` + suffix + ` * pricePerMemory` + suffix + `)
			storageCost` + suffix + ` as math(storage` + suffix + ` * durationInHours` + suffix + ` * ` + models.DefaultStorageCostPerGBPerHour + `)`
}

func getQueryForAggregatingChildMetricsWithAlias(childSuffix string) string {
	return `name
			type
			cpu: sum(val(cpu` + childSuffix + `))
			memory: sum(val(memory` + childSuffix + `))
			storage: sum(val(storage` + childSuffix + `))
			cpuCost: sum(val(cpuCost` + childSuffix + `))
			memoryCost: sum(val(memoryCost` + childSuffix + `))
			storageCost: sum(val(storageCost` + childSuffix + `))`
}

func getQueryForAggregatingChildMetrics(parentSuffix, childSuffix string) string {
	return `cpu` + parentSuffix + ` as sum(val(cpu` + childSuffix + `))
			memory` + parentSuffix + ` as sum(val(memory` + childSuffix + `))
			storage` + parentSuffix + ` as sum(val(storage` + childSuffix + `))
			cpuCost` + parentSuffix + ` as sum(val(cpuCost` + childSuffix + `))
			memoryCost` + parentSuffix + ` as sum(val(memoryCost` + childSuffix + `))
			storageCost` + parentSuffix + ` as sum(val(storageCost` + childSuffix + `))`
}

func getQueryFromSubQueryWithAlias(suffix string) string {
	return `name
			type
			cpu: val(cpu` + suffix + `)
			memory: val(memory` + suffix + `)
			storage: val(storage` + suffix + `)
			cpuCost: val(cpuCost` + suffix + `)
			memoryCost: val(memoryCost` + suffix + `)
			storageCost: val(storageCost` + suffix + `)`
}

func getQueryForHierarchy(parentCheck, parentType, parentName, childFilter string) string {
	return `query {
		parent(func: has(` + parentCheck + `)) @filter(eq(name, "` + parentName + `")) {
			name
			type
			children: ~` + parentType + ` ` + childFilter + ` {
				name
				type
			}
		}
	}`
}
