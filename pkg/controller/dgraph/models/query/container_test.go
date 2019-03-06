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

	"github.com/vmware/purser/pkg/controller/dgraph"

	"github.com/stretchr/testify/assert"
)

func setupForContainerDgraphDataRetrieve(isDataToBeFoundInDgraph bool) {
	executeQuery = func(query string, root interface{}) error {
		if query == containerMetricTestQuery {
			if isDataToBeFoundInDgraph {
				dummyParentWrapper, ok := root.(*ParentWrapper)
				if !ok {
					return fmt.Errorf("wrong root received")
				}
				parent := Parent{
					Name:       testContainerName,
					Type:       "container",
					CPU:        0.5,
					Memory:     1,
					CPUCost:    0.000005,
					MemoryCost: 0.000004,
				}
				dummyParentWrapper.Parent = []Parent{parent}
				return nil
			} else {
				// No data found for the given container in dgraph
				return nil
			}
		} else if query == containerHierarchyTestQuery {
			dummyParentWrapper, ok := root.(*ParentWrapper)
			if !ok {
				return fmt.Errorf("wrong root received")
			}
			parent := Parent{
				Name: testContainerName,
				Type: "container",
			}
			dummyParentWrapper.Parent = []Parent{parent}
			return nil
		}
		return fmt.Errorf("wrong query received")
	}
}

func shutdownForContainerDgraphDataRetrieve() {
	executeQuery = dgraph.ExecuteQuery
}

// TestRetrieveContainerHierarchyWithNameEmpty ...
func TestRetrieveContainerHierarchyWithNameEmpty(t *testing.T) {
	got := RetrieveContainerHierarchy("")
	expected := JSONDataWrapper{}
	assert.Equal(t, expected, got)
}

// TestRetrieveContainerHierarchy ...
func TestRetrieveContainerHierarchy(t *testing.T) {
	setupForContainerDgraphDataRetrieve(testDataFoundInDgraph)
	defer shutdownForContainerDgraphDataRetrieve()

	got := RetrieveContainerHierarchy(testContainerName)
	expected := JSONDataWrapper{
		Data: ParentWrapper{
			Name: testContainerName,
			Type: "container",
		},
	}
	assert.Equal(t, expected, got)
}

// TestRetrieveContainerMetricsWithNameEmpty ...
func TestRetrieveContainerMetricsWithNameEmpty(t *testing.T) {
	got := RetrieveContainerMetrics("")
	expected := JSONDataWrapper{}
	assert.Equal(t, expected, got)
}

// TestRetrieveContainerMetrics ...
func TestRetrieveContainerMetrics(t *testing.T) {
	setupForContainerDgraphDataRetrieve(testDataFoundInDgraph)
	defer shutdownForContainerDgraphDataRetrieve()

	got := RetrieveContainerMetrics(testContainerName)
	expected := JSONDataWrapper{
		Data: ParentWrapper{
			Name:       testContainerName,
			Type:       "container",
			CPU:        0.5,
			Memory:     1,
			CPUCost:    0.000005,
			MemoryCost: 0.000004,
		},
	}
	assert.Equal(t, expected, got)
}

// TestRetrieveContainerMetricsWithErrorFromDgraph ...
func TestRetrieveContainerMetricsWithErrorFromDgraph(t *testing.T) {
	setupForContainerDgraphDataRetrieve(testDataFoundInDgraph)
	defer shutdownForContainerDgraphDataRetrieve()

	got := RetrieveContainerMetrics("container-wrong-name")
	expected := JSONDataWrapper{}
	assert.Equal(t, expected, got)
}

// TestRetrieveContainerMetricsWithNoDataFromDgraph ...
func TestRetrieveContainerMetricsWithNoDataFromDgraph(t *testing.T) {
	setupForContainerDgraphDataRetrieve(testNoDataFoundInDgraph)
	defer shutdownForContainerDgraphDataRetrieve()

	got := RetrieveContainerMetrics(testContainerName)
	expected := JSONDataWrapper{}
	assert.Equal(t, expected, got)
}
