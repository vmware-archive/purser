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

func mockDgraphForPodQueries(dgraphError bool) {
	executeQuery = func(query string, root interface{}) error {
		if dgraphError {
			return fmt.Errorf("unable to connect/retrieve data from dgraph")
		}

		if query == podsWithLabelFilterTestQuery {
			dummyPodList, ok := root.(*podRoot)
			if !ok {
				return fmt.Errorf("wrong root received")
			}
			dummyPodList.Pods = []models.Pod{
				{
					ID:   dgraph.ID{UID: "0x3e283"},
					Name: testPodName,
				},
			}
			return nil
		} else if query == allLivePodsTestQuery {
			dummyPodList, ok := root.(*podRoot)
			if !ok {
				return fmt.Errorf("wrong root received")
			}
			dummyPodList.Pods = []models.Pod{
				{
					ID:   dgraph.ID{UID: "0x3e283", Xid: "purser:" + testPodName},
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
	mockDgraphForPodQueries(testDgraphError)
	labels := make(map[string][]string)
	labels["k1"] = []string{"v1"}
	_, err := RetrievePodsUIDsByLabelsFilter(labels)
	assert.Error(t, err)
}

// TestRetrievePodsUIDsByLabelsFilter ...
func TestRetrievePodsUIDsByLabelsFilter(t *testing.T) {
	mockDgraphForPodQueries(testNoDgraphError)
	labels := make(map[string][]string)
	labels["k1"] = []string{"v1"}
	got, err := RetrievePodsUIDsByLabelsFilter(labels)
	expected := []string{"0x3e283"}
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

// TestRetrieveAllLivePodsWithDgraphError ...
func TestRetrieveAllLivePodsWithDgraphError(t *testing.T) {
	mockDgraphForPodQueries(testDgraphError)
	got := RetrieveAllLivePods()
	assert.Nil(t, got)
}

// TestRetrieveAllLivePods ...
func TestRetrieveAllLivePods(t *testing.T) {
	mockDgraphForPodQueries(testNoDgraphError)
	got := RetrieveAllLivePods()
	expected := []models.Pod{
		{
			ID:   dgraph.ID{UID: "0x3e283", Xid: "purser:" + testPodName},
			Name: testPodName,
		},
	}
	assert.Equal(t, expected, got)
}

func TestPodInteractionsErrorCase(t *testing.T) {
	mockDgraphForPodQueries(testNoDgraphError)
	gotAllOrphan := RetrievePodsInteractions("", true)
	gotAllNonOrphan := RetrievePodsInteractions("", false)
	gotWithName := RetrievePodsInteractions("pod-purser", false)
	_, err := RetrievePodsInteractionsForAllLivePodsWithCount()
	assert.Nil(t, gotAllOrphan)
	assert.Nil(t, gotAllNonOrphan)
	assert.Nil(t, gotWithName)
	assert.Error(t, err)
}
