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
	"testing"

	"github.com/vmware/purser/pkg/controller/dgraph/models"

	"github.com/stretchr/testify/assert"
)

var dummyGroup models.Group

func mockMetricsMap(key string, value float64) map[string]float64 {
	metrics := make(map[string]float64)
	metrics[key] = value
	return metrics
}

func mockDgraphForGroupQueries(dgraphError bool) {
	executeQuery = func(query string, root interface{}) error {
		if dgraphError {
			return fmt.Errorf("unable to connect/retrieve data from dgraph")
		}

		if query == allGroupsDataTestQuery {
			dummyGroupList, ok := root.(*groupsRoot)
			if !ok {
				return fmt.Errorf("wrong root received")
			}
			dummyGroup = models.Group{
				Name:           "group-purser",
				PodsCount:      3,
				MtdCPU:         50.1,
				MtdMemory:      31.7,
				MtdStorage:     300,
				CPU:            4.1,
				Memory:         3.4,
				Storage:        10,
				MtdCPUCost:     5.01,
				MtdMemoryCost:  3.17,
				MtdStorageCost: 3,
				MtdCost:        11.18,
			}
			dummyGroupList.Groups = []models.Group{dummyGroup}
			return nil
		} else if query == groupMetricTestQuery {
			groupMetrics, ok := root.(*groupJSONMetrics)
			if !ok {
				return fmt.Errorf("wrong root received")
			}
			var jsonMetrics []map[string]float64
			jsonMetrics = append(jsonMetrics, mockMetricsMap("pitCPU", 1.3))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("pitMemory", 2.4))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("pitStorage", 2))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("pitCPULimit", 1.4))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("pitMemoryLimit", 2.5))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("mtdCPU", 13.1))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("mtdMemory", 24.2))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("mtdStorage", 20))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("mtdCPULimit", 14.1))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("mtdMemoryLimit", 25.2))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("cpuCost", 1.31))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("memoryCost", 2.42))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("storageCost", 0.21))
			jsonMetrics = append(jsonMetrics, mockMetricsMap("livePods", 2))
			groupMetrics.JSONMetrics = jsonMetrics
			return nil
		}

		return fmt.Errorf("no data found")
	}
}

// TestRetrieveGroupsDataWithDgraphError ...
func TestRetrieveGroupsDataWithDgraphError(t *testing.T) {
	mockDgraphForGroupQueries(testDgraphError)
	_, err := RetrieveGroupsData()
	assert.Error(t, err)
}

// TestRetrieveGroupsData ...
func TestRetrieveGroupsData(t *testing.T) {
	mockDgraphForGroupQueries(testNoDgraphError)
	got, err := RetrieveGroupsData()
	expected := []models.Group{{
		Name:           "group-purser",
		PodsCount:      3,
		MtdCPU:         50.1,
		MtdMemory:      31.7,
		MtdStorage:     300,
		CPU:            4.1,
		Memory:         3.4,
		Storage:        10,
		MtdCPUCost:     5.01,
		MtdMemoryCost:  3.17,
		MtdStorageCost: 3,
		MtdCost:        11.18,
	}}
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

// TestGroupMetricsFromPodUIDsWithDgraphError ...
func TestGroupMetricsFromPodUIDsWithDgraphError(t *testing.T) {
	mockDgraphForGroupQueries(testDgraphError)
	_, err := RetrieveGroupMetricsFromPodUIDs("")
	assert.Error(t, err)
}

// TestGroupMetricsFromPodUIDs ...
func TestGroupMetricsFromPodUIDs(t *testing.T) {
	mockDgraphForGroupQueries(testNoDgraphError)
	got, err := RetrieveGroupMetricsFromPodUIDs(testPodUIDs)
	expected := GroupMetrics{
		PITCpu:         1.3,
		PITMemory:      2.4,
		PITStorage:     2,
		PITCpuLimit:    1.4,
		PITMemoryLimit: 2.5,
		MTDCpu:         13.1,
		MTDMemory:      24.2,
		MTDStorage:     20,
		MTDCpuLimit:    14.1,
		MTDMemoryLimit: 25.2,
		CostCPU:        1.31,
		CostMemory:     2.42,
		CostStorage:    0.21,
		PodsCount:      2,
	}
	assert.Equal(t, expected, got)
	assert.NoError(t, err)
}
