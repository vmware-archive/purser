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
	"github.com/vmware/purser/pkg/controller/dgraph/models"

	"github.com/Sirupsen/logrus"
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

type groupsRoot struct {
	Groups []models.Group `json:"groups,omitempty"`
}

type groupJSONMetrics struct {
	JSONMetrics []map[string]float64 `json:"group"`
}

// RetrieveGroupsData returns list of models.Group objects in json format
// error is not nil if any failure is encountered
func RetrieveGroupsData() ([]models.Group, error) {
	query := getQueryForAllGroupsData()

	newRoot := groupsRoot{}
	err := executeQuery(query, &newRoot)
	if err != nil {
		return []models.Group{}, err
	}
	return newRoot.Groups, nil
}

// RetrieveGroupMetricsFromPodUIDs ...
func RetrieveGroupMetricsFromPodUIDs(podsUIDs string) (GroupMetrics, error) {
	query := getQueryForGroupMetrics(podsUIDs)

	newRoot := groupJSONMetrics{}
	err := executeQuery(query, &newRoot)
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
	case "livePods":
		groupMetrics.PodsCount = int(value)
	}
}
