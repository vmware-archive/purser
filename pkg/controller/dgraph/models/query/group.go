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
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/utils"
)

// GroupMetrics structure
type GroupMetrics struct {
	PITCpu         float64
	PITMemory      float64
	PITStorage     float64
	PITCpuLimit    float64
	PITMemoryLimit float64
	MTDCpu         float64
	MTDMemory      float64
	MTDStorage     float64
	MTDCpuLimit    float64
	MTDMemoryLimit float64
	CostCPU        float64
	CostMemory     float64
	CostStorage    float64
	PodsCount      int
}

// RetrieveGroupsData returns list of models.Group objects in json format
// error is not nil if any failure is encountered
func RetrieveGroupsData() ([]models.Group, error) {
	query := `query {
		groups(func: has(isGroup)) {
			name
			podsCount
			mtdCPU
			mtdMemory
			mtdStorage
			cpu
			memory
			storage
			mtdCPUCost
			mtdMemoryCost
			mtdStorageCost
			mtdCost
		}
	}`

	type root struct {
		Groups []models.Group `json:groups`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(query, &newRoot)
	return newRoot.Groups, err
}

// RetrieveGroupMetricsFromPodUIDs ...
func RetrieveGroupMetricsFromPodUIDs(podsUIDs string) (GroupMetrics, error) {
	secondsSinceMonthStart := fmt.Sprintf("%f", utils.GetSecondsSince(utils.GetCurrentMonthStartTime()))
	query := `query {
		var(func: uid(` + podsUIDs + `)) {
			podCpu as cpuRequest
			podMemory as memoryRequest
			pvcStorage as storageRequest
			podCpuLimit as cpuLimit
			podMemoryLimit as memoryLimit
			cpuRequestCount as count(cpuRequest)
			memoryRequestCount as count(memoryRequest)
			storageRequestCount as count(storageRequest)
			cpuLimitCount as count(cpuLimit)
			memoryLimitCount as count(memoryLimit)
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math((secondsSinceStart - secondsSinceEnd) / 3600)
			pitPodCPU as math(cond(isTerminated == 0, cond(cpuRequestCount > 0, podCpu, 0.0), 0.0))
			pitPodMemory as math(cond(isTerminated == 0, cond(memoryRequestCount > 0, podMemory, 0.0), 0.0))
			pitPvcStorage as math(cond(isTerminated == 0, cond(storageRequestCount > 0, pvcStorage, 0.0), 0.0))
			pitPodCPULimit as math(cond(isTerminated == 0, cond(cpuLimitCount > 0, podCpuLimit, 0.0), 0.0))
			pitPodMemoryLimit as math(cond(isTerminated == 0, cond(memoryLimitCount > 0, podMemoryLimit, 0.0), 0.0))
			mtdPodCPU as math(podCpu * durationInHours)
			mtdPodMemory as math(podMemory * durationInHours)
			mtdPvcStorage as math(pvcStorage * durationInHours)
			mtdPodCPULimit as math(podCpuLimit * durationInHours)
			mtdPodMemoryLimit as math(podMemoryLimit * durationInHours)
			podCpuCost as math(mtdPodCPU * ` + defaultCPUCostPerCPUPerHour + `)
			podMemoryCost as math(mtdPodMemory * ` + defaultMemCostPerGBPerHour + `)
			podStorageCost as math(mtdPvcStorage * ` + defaultStorageCostPerGBPerHour + `)
		}
		
		group() {
			pitCPU: sum(val(pitPodCPU))
			pitMemory: sum(val(pitPodMemory))
			pitStorage: sum(val(pitPvcStorage))
			pitCPULimit: sum(val(pitPodCPULimit))
			pitMemoryLimit: sum(val(pitPodMemoryLimit))
			mtdCPU: sum(val(mtdPodCPU))
			mtdMemory: sum(val(mtdPodMemory))
			mtdStorage: sum(val(mtdPvcStorage))
			mtdCPULimit: sum(val(mtdPodCPULimit))
			mtdMemoryLimit: sum(val(mtdPodMemoryLimit))
			cpuCost: sum(val(podCpuCost))
			memoryCost: sum(val(podMemoryCost))
			storageCost: sum(val(podStorageCost))
		}
	}`

	type root struct {
		JSONMetrics []map[string]float64 `json:"group"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(query, &newRoot)
	if err != nil {
		return GroupMetrics{}, err
	}
	return convertToGroupMetrics(newRoot.JSONMetrics), nil
}

func convertToGroupMetrics(jsonMetrics []map[string]float64) GroupMetrics {
	var groupMetrics GroupMetrics
	for _, data := range jsonMetrics {
		for key, value := range data {
			populateMetric(&groupMetrics, key, value)
			break
		}
	}
	logrus.Debugf("JSON metrics: (%v), Group metrics: (%v)", jsonMetrics, groupMetrics)
	return groupMetrics
}

// nolint: gocyclo
func populateMetric(groupMetrics *GroupMetrics, key string, value float64) {
	logrus.Debugf("key: %s", key)
	switch key {
	case "pitCPU":
		groupMetrics.PITCpu = value
	case "pitMemory":
		groupMetrics.PITMemory = value
	case "pitStorage":
		groupMetrics.PITStorage = value
	case "pitCPULimit":
		groupMetrics.PITCpuLimit = value
	case "pitMemoryLimit":
		groupMetrics.PITMemoryLimit = value
	case "mtdCPU":
		groupMetrics.MTDCpu = value
	case "mtdMemory":
		groupMetrics.MTDMemory = value
	case "mtdStorage":
		groupMetrics.MTDStorage = value
	case "mtdCPULimit":
		groupMetrics.MTDCpuLimit = value
	case "mtdMemoryLimit":
		groupMetrics.MTDMemoryLimit = value
	case "cpuCost":
		groupMetrics.CostCPU = value
	case "memoryCost":
		groupMetrics.CostMemory = value
	case "storageCost":
		groupMetrics.CostStorage = value
	}
}
