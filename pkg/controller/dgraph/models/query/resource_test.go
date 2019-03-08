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

func mockDgraphForPodPrices(queryType string) {
	executeQuery = func(query string, root interface{}) error {
		if queryType == testWrongQuery {
			return fmt.Errorf("unable to retrieve data from dgraph")
		}

		newRoot, ok := root.(*podRoot)
		if !ok {
			return fmt.Errorf("wrong pod root received")
		}
		pod := models.Pod{
			CPUPrice:    testCPUPrice,
			MemoryPrice: testMemoryPrice,
		}
		newRoot.Pods = []models.Pod{pod}
		return nil
	}
}

func mockDgraphForResourceQueries(queryType string) {
	executeQuery = func(query string, root interface{}) error {
		dummyParentWrapper, ok := root.(*ParentWrapper)
		if !ok {
			return fmt.Errorf("wrong root received")
		}

		var parent Parent
		if queryType == testMetrics {
			firstPodWithMetrics := Children{
				Name:        "pod-purser-1",
				Type:        PodType,
				CPU:         0.25,
				Memory:      0.1,
				Storage:     1.2,
				CPUCost:     0.024,
				MemoryCost:  0.09,
				StorageCost: 0.1,
			}
			secondPodWithMetrics := Children{
				Name:        "pod-purser-2",
				Type:        PodType,
				CPU:         0.15,
				Memory:      0.2,
				Storage:     0.2,
				CPUCost:     0.014,
				MemoryCost:  0.19,
				StorageCost: 0.01,
			}
			parent = Parent{
				Name:        testDaemonsetName,
				Type:        DaemonsetType,
				Children:    []Children{firstPodWithMetrics, secondPodWithMetrics},
				CPU:         0.40,
				Memory:      0.28,
				Storage:     1.4,
				CPUCost:     0.038,
				MemoryCost:  0.28,
				StorageCost: 0.11,
			}
			dummyParentWrapper.Parent = []Parent{parent}
			return nil
		} else if queryType == testHierarchy {
			firstPod := Children{
				Name: "pod-purser-1",
				Type: PodType,
			}
			secondPod := Children{
				Name: "pod-purser-2",
				Type: PodType,
			}
			parent = Parent{
				Name:     testDaemonsetName,
				Type:     DaemonsetType,
				Children: []Children{firstPod, secondPod},
			}
			dummyParentWrapper.Parent = []Parent{parent}
			return nil
		}
		return fmt.Errorf("unable to retrieve data from dgraph")
	}
}

// TestRetrieveResourceHierarchyWithNameEmpty ...
func TestRetrieveResourceHierarchyWithNameEmpty(t *testing.T) {
	input := &Resource{
		Check:       DaemonsetCheck,
		Type:        DaemonsetType,
		Name:        "",
		ChildFilter: IsPodFilter,
	}
	got := input.RetrieveResourceHierarchy()
	expected := JSONDataWrapper{}
	assert.Equal(t, expected, got)
}

// TestRetrieveResourceHierarchy ...
func TestRetrieveResourceHierarchy(t *testing.T) {
	mockDgraphForResourceQueries(testHierarchy)

	input := &Resource{
		Check:       DaemonsetCheck,
		Type:        DaemonsetType,
		Name:        testDaemonsetName,
		ChildFilter: IsPodFilter,
	}
	got := input.RetrieveResourceHierarchy()

	firstPod := Children{
		Name: "pod-purser-1",
		Type: PodType,
	}
	secondPod := Children{
		Name: "pod-purser-2",
		Type: PodType,
	}
	expected := JSONDataWrapper{
		Data: ParentWrapper{
			Name:     testDaemonsetName,
			Type:     DaemonsetType,
			Children: []Children{firstPod, secondPod},
		},
	}
	assert.Equal(t, expected, got)
}

// TestRetrieveResourceHierarchyWithDgraphError ...
func TestRetrieveResourceHierarchyWithDgraphError(t *testing.T) {
	mockDgraphForResourceQueries(testWrongQuery)

	input := &Resource{
		Check:       DaemonsetCheck,
		Type:        DaemonsetType,
		Name:        testDaemonsetName,
		ChildFilter: IsPodFilter,
	}
	got := input.RetrieveResourceHierarchy()
	expected := JSONDataWrapper{}
	assert.Equal(t, expected, got)
}

// TestRetrieveResourceMetricsWithNameEmpty ...
func TestRetrieveResourceMetricsWithNameEmpty(t *testing.T) {
	input := &Resource{
		Check: DaemonsetCheck,
		Type:  DaemonsetType,
		Name:  "",
	}
	got := input.RetrieveResourceMetrics()
	expected := JSONDataWrapper{}
	assert.Equal(t, expected, got)
}

// TestRetrieveDaemonsetMetrics ...
func TestRetrieveDaemonsetMetrics(t *testing.T) {
	mockDgraphForResourceQueries(testMetrics)

	input := &Resource{
		Check: DaemonsetCheck,
		Type:  DaemonsetType,
		Name:  testDaemonsetName,
	}
	got := input.RetrieveResourceMetrics()

	firstPodWithMetrics := Children{
		Name:        "pod-purser-1",
		Type:        PodType,
		CPU:         0.25,
		Memory:      0.1,
		Storage:     1.2,
		CPUCost:     0.024,
		MemoryCost:  0.09,
		StorageCost: 0.1,
	}
	secondPodWithMetrics := Children{
		Name:        "pod-purser-2",
		Type:        PodType,
		CPU:         0.15,
		Memory:      0.2,
		Storage:     0.2,
		CPUCost:     0.014,
		MemoryCost:  0.19,
		StorageCost: 0.01,
	}
	expected := JSONDataWrapper{
		Data: ParentWrapper{
			Name:        testDaemonsetName,
			Type:        DaemonsetType,
			Children:    []Children{firstPodWithMetrics, secondPodWithMetrics},
			CPU:         0.40,
			Memory:      0.28,
			Storage:     1.4,
			CPUCost:     0.038,
			MemoryCost:  0.28,
			StorageCost: 0.11,
		},
	}
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsPod ...
func TestGetQueryForResourceMetricsPod(t *testing.T) {
	mockDgraphForPodPrices(testPodName)

	input := &Resource{
		Check: PodCheck,
		Type:  PodType,
		Name:  testPodName,
	}
	got := input.getQueryForResourceMetrics()
	expected := podMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsPod ...
func TestGetQueryForResourceMetricsPodWithError(t *testing.T) {
	mockDgraphForPodPrices(testWrongQuery)

	input := &Resource{
		Check: PodCheck,
		Type:  PodType,
		Name:  "pod-wrong",
	}
	got := input.getQueryForResourceMetrics()
	expected := podMetricTestQuery
	assert.NotEqual(t, got, expected)
}

// TestGetQueryForResourceMetricsDeployment ...
func TestGetQueryForResourceMetricsDeployment(t *testing.T) {
	input := &Resource{
		Check: DeploymentCheck,
		Type:  DeploymentType,
		Name:  testDeploymentName,
	}
	got := input.getQueryForResourceMetrics()
	expected := deploymentMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsNamespace ...
func TestGetQueryForResourceMetricsNamespace(t *testing.T) {
	input := &Resource{
		Check: NamespaceCheck,
		Type:  NamespaceType,
		Name:  testNamespaceName,
	}
	got := input.getQueryForResourceMetrics()
	expected := namespaceMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsNode ...
func TestGetQueryForResourceMetricsNode(t *testing.T) {
	input := &Resource{
		Check: NodeCheck,
		Type:  NodeType,
		Name:  testNodeName,
	}
	got := input.getQueryForResourceMetrics()
	expected := nodeMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsPV ...
func TestGetQueryForResourceMetricsPV(t *testing.T) {
	input := &Resource{
		Check: PVCheck,
		Type:  PVType,
		Name:  testPVName,
	}
	got := input.getQueryForResourceMetrics()
	expected := pvMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsPVC ...
func TestGetQueryForResourceMetricsPVC(t *testing.T) {
	input := &Resource{
		Check: PVCCheck,
		Type:  PVCType,
		Name:  testPVCName,
	}
	got := input.getQueryForResourceMetrics()
	expected := pvcMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsContainer ...
func TestGetQueryForResourceMetricsContainer(t *testing.T) {
	input := &Resource{
		Check: ContainerCheck,
		Type:  ContainerType,
		Name:  testContainerName,
	}
	got := input.getQueryForResourceMetrics()
	expected := containerMetricTestQuery
	assert.Equal(t, expected, got)
}
