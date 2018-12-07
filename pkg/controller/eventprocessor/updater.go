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

package eventprocessor

import (
	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
	"github.com/vmware/purser/pkg/controller/discovery/processor"

	log "github.com/Sirupsen/logrus"

	groups_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	groupsClient_v1 "github.com/vmware/purser/pkg/client/clientset/typed/groups/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// UpdateGroups retrieve all groups and updates them
func UpdateGroups(groupCRDClient *groupsClient_v1.GroupClient) {
	groups := processor.RetrieveGroupList(groupCRDClient, meta_v1.ListOptions{})
	if groups != nil {
		for _, group := range groups.Items {
			UpdateGroup(group)
		}
	}
}

func UpdateGroup(group *groups_v1.Group) {
	if group == nil {
		return
	}
	groupMetrics := getGroupMetrics(group)
	group.Spec.MTDMetrics = &groups_v1.GroupMetrics{
		CPURequest:    groupMetrics.MTDCpu,
		MemoryRequest: groupMetrics.MTDMemory,
		StorageClaim:  groupMetrics.MTDStorage,
	}
	group.Spec.PITMetrics = &groups_v1.GroupMetrics{
		CPURequest:    groupMetrics.PITCpu,
		MemoryRequest: groupMetrics.PITMemory,
		StorageClaim:  groupMetrics.PITStorage,
	}
	group.Spec.MTDCost = &groups_v1.Cost{
		CPUCost:     groupMetrics.CostCpu,
		MemoryCost:  groupMetrics.CostMemory,
		StorageCost: groupMetrics.CostStorage,
		TotalCost:   groupMetrics.CostCpu + groupMetrics.CostMemory + groupMetrics.CostStorage,
	}
}

func getGroupMetrics(group *groups_v1.Group) query.GroupMetrics {
	// for each label-expression retrieve UIDs of pods that satisfy the label-expression
	podUIDsFromExpressions := getPodUIDsFromAllExpressions(group.Spec.Expressions)

	// Across all the podUIDs computed from label-expressions, map pod's UID with number of occurrences of it
	podUIDsCounter := mapPodUIDsToNumberOfOccurences(podUIDsFromExpressions)

	// if number of occurrences of UID == number of expressions that means the pod satisfies all the expressions(i.e, AND)
	// get uid-query to retrieve such pods i.e, "uid1, uid2, uid2..."
	uidQueryForPods := getUIDQueryForPods(podUIDsCounter, len(group.Spec.Expressions))

	// given group metrics
	groupMetrics, err := query.RetrieveGroupMetricsFromPodUIDs(uidQueryForPods)
	if err != nil {
		log.Errorf("Unable to retrieve group metrics. UIDs: (%v)", uidQueryForPods)
		return query.GroupMetrics{}
	}
	return groupMetrics
}

// for each label-expression retrieve UIDs of pods that satisfy the label-expression
// appends results from each expression and returns the array of such results i.e,
// [[pod1-from-exp1, pod2-from-exp1], [pod1-from-exp2], [pod1-from-exp3, pod2-from-exp3, pod3-from-exp3]]
func getPodUIDsFromAllExpressions(expressions map[string]groups_v1.Selector) [][]string {
	var podsUIDsFromExpressions [][]string
	for _, selector := range expressions {
		podsFromSelector, err := query.RetrievePodsUIDsByLabelsFilter(selector.Labels)
		if err == nil {
			podsUIDsFromExpressions = append(podsUIDsFromExpressions, podsFromSelector)
		}
	}
	return podsUIDsFromExpressions
}

// across all the podUIDs computed from label-expressions, map pod's UID with number of occurrences of it and return map
func mapPodUIDsToNumberOfOccurences(podsFromExpressions [][]string) map[string]int {
	podsUIDsCounter := make(map[string]int)
	for _, podsFromExpression := range podsFromExpressions {
		for _, pod := range podsFromExpression {
			if _, isPresent := podsUIDsCounter[pod]; !isPresent {
				podsUIDsCounter[pod] = 0
			}
			podsUIDsCounter[pod]++
		}
	}
	return podsUIDsCounter
}

// returns UIDs for pods that satisfy (number of its occurrences == expressions count) i.e,
// if number of occurrences of UID == number of expressions that means the pod satisfies all the expressions(-> AND)
// returns uid-query(i.e, "uid1, uid2, uid2...") that can retrieve desired pods
func getUIDQueryForPods(podsUIDsCounter map[string]int, expressionsCount int) string {
	separator := ", "
	isFirst := true
	var uidQueryForPods string
	for podUID, count := range podsUIDsCounter {
		if count == expressionsCount {
			if !isFirst {
				uidQueryForPods += separator
			} else {
				isFirst = false
			}
			uidQueryForPods += podUID
		}
	}
	return uidQueryForPods
}
