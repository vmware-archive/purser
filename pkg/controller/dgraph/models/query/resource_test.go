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

var firstPod, secondPod, firstPodWithMetrics, secondPodWithMetrics Children

func setupForResources() {
	firstPod = Children{
		Name: "pod-purser-1",
		Type: PodType,
	}
	secondPod = Children{
		Name: "pod-purser-2",
		Type: PodType,
	}
	firstPodWithMetrics = Children{
		Name:        "pod-purser-1",
		Type:        PodType,
		CPU:         0.25,
		Memory:      0.1,
		Storage:     1.2,
		CPUCost:     0.024,
		MemoryCost:  0.09,
		StorageCost: 0.1,
	}
	secondPodWithMetrics = Children{
		Name:        "pod-purser-2",
		Type:        PodType,
		CPU:         0.15,
		Memory:      0.2,
		Storage:     0.2,
		CPUCost:     0.014,
		MemoryCost:  0.19,
		StorageCost: 0.01,
	}
}

func mockDgraphForResourceQueries(isHierarchy, dgraphError bool) {
	setupForResources()
	executeQuery = func(query string, root interface{}) error {
		if dgraphError {
			return fmt.Errorf("error while executing query")
		}

		if query == podPriceTestQuery {
			newRoot, ok := root.(*podRoot)
			if !ok {
				return fmt.Errorf("wrong pod root received")
			}
			pod := models.Pod{
				CPUPrice:    0.24,
				MemoryPrice: 0.1,
			}
			newRoot.Pods = []models.Pod{pod}
			return nil
		}

		dummyParentWrapper, ok := root.(*ParentWrapper)
		if !ok {
			return fmt.Errorf("wrong root received")
		}

		var parent Parent
		if !isHierarchy {
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
		} else {
			parent = Parent{
				Name:     testDaemonsetName,
				Type:     DaemonsetType,
				Children: []Children{firstPod, secondPod},
			}
		}
		dummyParentWrapper.Parent = []Parent{parent}
		return nil
	}
}

// TestRetrieveResourceHierarchyWithNameEmpty ...
func TestRetrieveResourceHierarchyWithNameEmpty(t *testing.T) {
	got := RetrieveResourceHierarchy(DaemonsetCheck, DaemonsetType, "", IsPodFilter)
	expected := JSONDataWrapper{}
	assert.Equal(t, expected, got)
}

// TestRetrieveResourceHierarchy ...
func TestRetrieveResourceHierarchy(t *testing.T) {
	mockDgraphForResourceQueries(testHierarchy, testNoDgraphError)

	got := RetrieveResourceHierarchy(DaemonsetCheck, DaemonsetType, testDaemonsetName, IsPodFilter)
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
	mockDgraphForResourceQueries(testHierarchy, testDgraphError)

	got := RetrieveResourceHierarchy(DaemonsetCheck, DaemonsetType, testDaemonsetName, IsPodFilter)
	expected := JSONDataWrapper{}
	assert.Equal(t, expected, got)
}

// TestRetrieveResourceMetricsWithNameEmpty ...
func TestRetrieveResourceMetricsWithNameEmpty(t *testing.T) {
	got := RetrieveResourceMetrics(DaemonsetCheck, DaemonsetType, "")
	expected := JSONDataWrapper{}
	assert.Equal(t, expected, got)
}

// TestRetrieveDaemonsetMetrics ...
func TestRetrieveDaemonsetMetrics(t *testing.T) {
	mockDgraphForResourceQueries(testMetrics, testNoDgraphError)

	got := RetrieveResourceMetrics(DaemonsetCheck, DaemonsetType, testDaemonsetName)
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

// TestGetQueryForResourceMetricsDeployment ...
func TestGetQueryForResourceMetricsDeployment(t *testing.T) {
	got := getQueryForResourceMetrics(DeploymentCheck, DeploymentType, testDeploymentName)
	expected := deploymentMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsNamespace ...
func TestGetQueryForResourceMetricsNamespace(t *testing.T) {
	got := getQueryForResourceMetrics(NamespaceCheck, NamespaceType, testNamespaceName)
	expected := namespaceMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsNode ...
func TestGetQueryForResourceMetricsNode(t *testing.T) {
	got := getQueryForResourceMetrics(NodeCheck, NodeType, testNodeName)
	expected := nodeMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsPV ...
func TestGetQueryForResourceMetricsPV(t *testing.T) {
	got := getQueryForResourceMetrics(PVCheck, PVType, testPVName)
	expected := pvMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsPVC ...
func TestGetQueryForResourceMetricsPVC(t *testing.T) {
	got := getQueryForResourceMetrics(PVCCheck, PVCType, testPVCName)
	expected := pvcMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsContainer ...
func TestGetQueryForResourceMetricsContainer(t *testing.T) {
	got := getQueryForResourceMetrics(ContainerCheck, ContainerType, testContainerName)
	expected := containerMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsPod ...
func TestGetQueryForResourceMetricsPod(t *testing.T) {
	mockDgraphForResourceQueries(testMetrics, testNoDgraphError)
	got := getQueryForResourceMetrics(PodCheck, PodType, testPodName)
	expected := podMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForResourceMetricsPod ...
func TestGetQueryForResourceMetricsPodWithError(t *testing.T) {
	mockDgraphForResourceQueries(testMetrics, testNoDgraphError)
	got := getQueryForResourceMetrics(PodCheck, PodType, "pod-wrong")
	expected := podMetricTestQuery
	assert.NotEqual(t, got, expected)
}
