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
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/utils"
)

// GroupMetrics structure
type GroupMetrics struct {
	PITCpu      float64
	PITMemory   float64
	PITStorage  float64
	MTDCpu      float64
	MTDMemory   float64
	MTDStorage  float64
	CostCPU     float64
	CostMemory  float64
	CostStorage float64
}

// RetrieveGroupMetricsFromPodUIDs ...
func RetrieveGroupMetricsFromPodUIDs(podsUIDs string) (GroupMetrics, error) {
	secondsSinceMonthStart := fmt.Sprintf("%f", utils.GetSecondsSince(utils.GetCurrentMonthStartTime()))
	query := `query {
		var(func: uid(` + podsUIDs +`)) {
			podCpu as cpuRequest
			podMemory as memoryRequest
			pvcStorage as storageRequest
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math((secondsSinceStart - secondsSinceEnd) / 3600)
			pitPodCPU as math(cond(isTerminated == 0, podCpu, 0.0))
			pitPodMemory as math(cond(isTerminated == 0, podMemory, 0.0))
			pitPvcStorage as math(cond(isTerminated == 0, pvcStorage, 0.0))
			mtdPodCPU as math(podCpu * durationInHours)
			mtdPodMemory as math(podMemory * durationInHours)
			mtdPvcStorage as math(pvcStorage * durationInHours)
			podCpuCost as math(mtdPodCPU * ` + defaultCPUCostPerCPUPerHour + `)
			podMemoryCost as math(mtdPodMemory * ` + defaultMemCostPerGBPerHour + `)
			podStorageCost as math(mtdPvcStorage * ` + defaultStorageCostPerGBPerHour + `)
		}
		
		group() {
			pitCPU: sum(val(pitPodCPU))
			pitMemory: sum(val(pitPodMemory))
			pitStorage: sum(val(pitPvcStorage))
			mtdCPU: sum(val(mtdPodCPU))
			mtdMemory: sum(val(mtdPodMemory))
			mtdStorage: sum(val(mtdPvcStorage))
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

func populateMetric(groupMetrics *GroupMetrics, key string, value float64) {
	logrus.Debugf("key: %s", key)
	switch key {
	case "pitCPU":
		groupMetrics.PITCpu = value
	case "pitMemory":
		groupMetrics.PITMemory = value
	case "pitStorage":
		groupMetrics.PITStorage = value
	case "mtdCPU":
		groupMetrics.MTDCpu = value
	case "mtdMemory":
		groupMetrics.MTDMemory = value
	case "mtdStorage":
		groupMetrics.MTDStorage = value
	case "cpuCost":
		groupMetrics.CostCPU = value
	case "memoryCost":
		groupMetrics.CostMemory = value
	case "storageCost":
		groupMetrics.CostStorage = value
	}
}
