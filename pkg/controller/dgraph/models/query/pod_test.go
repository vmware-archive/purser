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

	"github.com/stretchr/testify/assert"

	"github.com/vmware/purser/pkg/controller/dgraph"

	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

func mockDgraphForPodQueries(queryType string) {
	executeQuery = func(query string, root interface{}) error {
		if queryType == testLabelFilterPods {
			dummyPodList, ok := root.(*podRoot)
			if !ok {
				return fmt.Errorf("wrong root received")
			}
			dummyPodList.Pods = []models.Pod{
				{
					ID:   dgraph.ID{UID: testPodUID},
					Name: testPodName,
				},
			}
			return nil
		} else if queryType == testAlivePods {
			dummyPodList, ok := root.(*podRoot)
			if !ok {
				return fmt.Errorf("wrong root received")
			}
			dummyPodList.Pods = []models.Pod{
				{
					ID:   dgraph.ID{UID: testPodUID, Xid: testPodXID},
					Name: testPodName,
				},
			}
			return nil
		}

		return fmt.Errorf("no data found")
	}

	executeQueryRaw = func(query string) ([]byte, error) {
		return nil, fmt.Errorf("pod interactions err")
	}
}

// TestRetrievePodsUIDsByLabelsFilterWithError ...
func TestRetrievePodsUIDsByLabelsFilterWithError(t *testing.T) {
	mockDgraphForPodQueries(testWrongQuery)

	// input setup
	labels := make(map[string][]string)
	labels["k1"] = []string{"v1"}
	inputLabelFilter := CreateFilterFromListOfLabels(labels)

	_, err := RetrievePodsUIDsByLabelsFilter(inputLabelFilter)
	assert.Error(t, err)
}

// TestRetrievePodsUIDsByLabelsFilter ...
func TestRetrievePodsUIDsByLabelsFilter(t *testing.T) {
	mockDgraphForPodQueries(testLabelFilterPods)

	// input setup
	labels := make(map[string][]string)
	labels["k1"] = []string{"v1"}
	inputLabelFilter := CreateFilterFromListOfLabels(labels)

	got, err := RetrievePodsUIDsByLabelsFilter(inputLabelFilter)
	expected := []string{testPodUID}
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

// TestRetrieveAllLivePodsWithDgraphError ...
func TestRetrieveAllLivePodsWithDgraphError(t *testing.T) {
	mockDgraphForPodQueries(testWrongQuery)
	got := RetrieveAllLivePods()
	assert.Nil(t, got)
}

// TestRetrieveAllLivePods ...
func TestRetrieveAllLivePods(t *testing.T) {
	mockDgraphForPodQueries(testAlivePods)
	got := RetrieveAllLivePods()
	expected := []models.Pod{
		{
			ID:   dgraph.ID{UID: testPodUID, Xid: testPodXID},
			Name: testPodName,
		},
	}
	assert.Equal(t, expected, got)
}

func TestPodInteractionsErrorCase(t *testing.T) {
	mockDgraphForPodQueries(testPodInteractions)
	gotAllOrphan := RetrievePodsInteractions("", true)
	gotAllNonOrphan := RetrievePodsInteractions("", false)
	gotWithName := RetrievePodsInteractions(testPodName, false)
	_, err := RetrievePodsInteractionsForAllLivePodsWithCount()
	assert.Nil(t, gotAllOrphan)
	assert.Nil(t, gotAllNonOrphan)
	assert.Nil(t, gotWithName)
	assert.Error(t, err)
}

func TestGetPricePerResourceForPodWithError(t *testing.T) {
	mockDgraphForResourceQueries(testWrongQuery, testPodName, PodType)
	gotCPUPrice, gotMemoryPrice := getPricePerResourceForPod(testPodName)
	expectedCPUPrice, expectedMemoryPrice := models.DefaultCPUCostInFloat64, models.DefaultMemCostInFloat64
	assert.Equal(t, expectedCPUPrice, gotCPUPrice)
	assert.Equal(t, expectedMemoryPrice, gotMemoryPrice)
}

func TestGetPricePerResourceForPod(t *testing.T) {
	mockDgraphForResourceQueries(testPodPrices, testPodName, PodType)
	gotCPUPrice, gotMemoryPrice := getPricePerResourceForPod(testPodName)
	expectedCPUPrice, expectedMemoryPrice := testCPUPrice, testMemoryPrice
	assert.Equal(t, expectedCPUPrice, gotCPUPrice)
	assert.Equal(t, expectedMemoryPrice, gotMemoryPrice)
}
