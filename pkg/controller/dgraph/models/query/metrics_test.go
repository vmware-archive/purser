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
	"os"
	"testing"

	"github.com/vmware/purser/pkg/controller/dgraph"

	"github.com/stretchr/testify/assert"
)

func mockSecondsSinceMonthStart() {
	secondsFromFirstOfCurrentMonth = func() string {
		return testSecondsSinceMonthStart
	}
}

func removeMocks() {
	secondsFromFirstOfCurrentMonth = getSecondsSinceMonthStart
	executeQuery = dgraph.ExecuteQuery
	executeQueryRaw = dgraph.ExecuteQueryRaw
}

// TestMain ...
func TestMain(m *testing.M) {
	mockSecondsSinceMonthStart()
	code := m.Run()
	removeMocks()
	os.Exit(code)
}

// TestGetQueryForGroupMetrics ...
func TestGetQueryForGroupMetrics(t *testing.T) {
	got := getQueryForGroupMetrics(testPodUIDList)
	expected := groupMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForAllGroupsData ...
func TestGetQueryForAllGroupsData(t *testing.T) {
	got := getQueryForAllGroupsData()
	expected := allGroupsDataTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForLogicalResources ...
func TestGetQueryForLogicalResources(t *testing.T) {
	got := getMetricsQueryForLogicalResources()
	expected := logicalResourcesMetricTestQuery
	assert.Equal(t, expected, got)
}

// TestGetQueryForPhysicalResources ...
func TestGetQueryForPhysicalResources(t *testing.T) {
	got := getMetricsQueryForPhysicalResources()
	expected := phycialResourcesMetricTestQuery
	assert.Equal(t, expected, got)
}
