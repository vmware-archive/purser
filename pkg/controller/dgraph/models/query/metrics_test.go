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

	"github.com/stretchr/testify/assert"
)

const mockSecondsSinceMonthStart = "1.45"

func setupForMetricQueryTesting() {
	secondsFromFirstOfCurrentMonth = func() string {
		return mockSecondsSinceMonthStart
	}
}

func shutdownForMetricQueryTesting() {
	secondsFromFirstOfCurrentMonth = getSecondsSinceMonthStart
}

func TestMain(m *testing.M) {
	setupForMetricQueryTesting()
	code := m.Run()
	shutdownForMetricQueryTesting()
	os.Exit(code)
}

func TestGetQueryForDeploymentMetrics(t *testing.T) {
	got := getQueryForDeploymentMetrics("deployment-purser")
	expected := deploymentMetricTestQuery
	assert.Equal(t, expected, got)
}

func TestGetQueryForGroupMetrics(t *testing.T) {
	got := getQueryForGroupMetrics("0x3e283, 0x3e288")
	expected := groupMetricTestQuery
	assert.Equal(t, expected, got)
}

func TestGetQueryForAllGroupsData(t *testing.T) {
	got := getQueryForAllGroupsData()
	expected := allGroupsDataTestQuery
	assert.Equal(t, expected, got)
}

func TestGetQueryForNamespaceMetrics(t *testing.T) {
	got := getQueryForNamespaceMetrics("namespace-default")
	expected := namespaceMetricTestQuery
	assert.Equal(t, expected, got)
}

func TestGetQueryForNodeMetrics(t *testing.T) {
	got := getQueryForNodeMetrics("node-default")
	expected := nodeMetricTestQuery
	assert.Equal(t, expected, got)
}

func TestGetQueryForPodMetrics(t *testing.T) {
	got := getQueryForPodMetrics("pod-purser-dgraph-0", "0.24", "0.1")
	expected := podMetricTestQuery
	assert.Equal(t, expected, got)
}

func TestGetQueryForPVMetrics(t *testing.T) {
	got := getQueryForPVMetrics("pv-datadir-purser-dgraph")
	expected := pvMetricTestQuery
	assert.Equal(t, expected, got)
}

func TestGetQueryForPVCMetrics(t *testing.T) {
	got := getQueryForPVCMetrics("pvc-datadir-purser-dgraph")
	expected := pvcMetricTestQuery
	assert.Equal(t, expected, got)
}

func TestGetQueryForContainerMetrics(t *testing.T) {
	got := getQueryForContainerMetrics("container-purser-controller")
	expected := containerMetricTestQuery
	assert.Equal(t, expected, got)
}

func TestGetQueryForLogicalResources(t *testing.T) {
	got := getQueryForLogicalResources()
	expected := logicalResourcesMetricTestQuery
	assert.Equal(t, expected, got)
}

func TestGetQueryForPhysicalResources(t *testing.T) {
	got := getQueryForPhysicalResources()
	expected := phycialResourcesMetricTestQuery
	assert.Equal(t, expected, got)
}
