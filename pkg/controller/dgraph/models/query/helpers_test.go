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
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetQueryForMetricsComputationWithAliasAndVariables ...
func TestGetQueryForMetricsComputationWithAliasAndVariables(t *testing.T) {
	got := getQueryForMetricsComputationWithAliasAndVariables("Pod")
	expected := testQueryForMetricsComputationWithAliasAndVariables
	assert.Equal(t, expected, got)
}

// TestGetQueryForAggregatingChildMetricsWithAlias ...
func TestGetQueryForAggregatingChildMetricsWithAlias(t *testing.T) {
	got := getQueryForAggregatingChildMetricsWithAlias("Pod")
	expected := testQueryForAggregatingChildMetricsWithAlias
	assert.Equal(t, expected, got)
}

// TestGetQueryForPodParentMetrics ...
func TestGetQueryForPodParentMetrics(t *testing.T) {
	got := getQueryForPodParentMetrics("isJob", "job", "job-purser")
	expected := testQueryForPodParentMetrics
	assert.Equal(t, expected, got)
}

// TestGetQueryForHierarchy ...
func TestGetQueryForHierarchy(t *testing.T) {
	got := getQueryForHierarchy("isNode", "node", "node-minikube", "@filter(has(isPod))")
	expected := testQueryForHierarchy
	assert.Equal(t, expected, got)
}

// TestGetSecondsSinceMonthStart ...
func TestGetSecondsSinceMonthStart(t *testing.T) {
	maxSecondsInAMonth := 2678400.0
	got := getSecondsSinceMonthStart()
	gotFloat, err := strconv.ParseFloat(got, 64)
	assert.NoError(t, err, "unable to convert secondsSinceMonthStart to float64")
	assert.False(t, gotFloat > maxSecondsInAMonth, "secondsSinceMonthStart can't be greater than 2678400")
	assert.False(t, gotFloat < 0, "secondsSinceMonthStart can't be less than 0")
}
