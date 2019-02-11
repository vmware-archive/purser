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
	"time"

	"github.com/vmware/purser/pkg/controller/dgraph/models"

	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
	"github.com/vmware/purser/pkg/controller/utils"

	log "github.com/Sirupsen/logrus"

	groups_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	groupsClient_v1 "github.com/vmware/purser/pkg/client/clientset/typed/groups/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// UpdateGroups retrieve all groups and updates them
func UpdateGroups(groupCRDClient *groupsClient_v1.GroupClient) {
	log.Infof("Started updating groups")
	groups := utils.RetrieveGroupList(groupCRDClient, meta_v1.ListOptions{})
	if groups == nil {
		log.Debugf("GroupList is nil")
		return
	}
	log.Debugf("Retrieved groups of length: %d", len(groups.Items))
	for _, group := range groups.Items {
		UpdateGroup(group, groupCRDClient)
	}
}

// UpdateGroup given a group it updates its spec with metrics
func UpdateGroup(group *groups_v1.Group, groupCRDClient *groupsClient_v1.GroupClient) {
	if group == nil {
		log.Warn("Received empty group to update")
		return
	}
	groupMetrics := getGroupMetrics(group)
	log.Debugf("GroupMetrics computed from dgraph data: (%v)", groupMetrics)
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
		CPUCost:     groupMetrics.CostCPU,
		MemoryCost:  groupMetrics.CostMemory,
		StorageCost: groupMetrics.CostStorage,
		TotalCost:   groupMetrics.CostCPU + groupMetrics.CostMemory + groupMetrics.CostStorage,
	}
	group.Spec.LastUpdated = time.Now()
	_, err := groupCRDClient.Update(group)
	if err != nil {
		log.Errorf("unable to update group: (%s), error: (%v)", group.Name, err)
	} else {
		log.Debugf("Updated group spec: (%v)", group.Spec)
		log.Infof("Group spec is updated with metrics for group: (%s)", group.Name)
		_, err = models.CreateOrUpdateGroup(group, groupMetrics.PodsCount)
		if err != nil {
			log.Errorf("unable to create or update group in dgraph: (%s), error: (%v)", group.Name, err)
		}
	}
}

func getGroupMetrics(group *groups_v1.Group) query.GroupMetrics {
	log.Debugf("Group: (%v), expressions: (%v)", group.Name, group.Spec.Expressions)

	// for each label-expression retrieve UIDs of pods that satisfy the label-expression
	podUIDsFromExpressions := getPodUIDsFromAllExpressions(group.Spec.Expressions)
	log.Debugf("Group: (%v), pod uids: (%v)", group.Name, podUIDsFromExpressions)

	// Across all the podUIDs computed from label-expressions, map pod's UID with number of occurrences of it
	podUIDsCounter := mapPodUIDsToNumberOfOccurences(podUIDsFromExpressions)
	log.Debugf("Group: (%v), pod uids counter: (%v)", group.Name, podUIDsCounter)

	// if number of occurrences of UID == number of expressions that means the pod satisfies all the expressions(i.e, AND)
	// get uid-query to retrieve such pods i.e, "uid1, uid2, uid2..."
	uidQueryForPods := getUIDQueryForPods(podUIDsCounter, len(group.Spec.Expressions))
	log.Debugf("Group: (%v), uidQuery: (%v)", group.Name, uidQueryForPods)

	// get group metrics
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
func getPodUIDsFromAllExpressions(expressions map[string]map[string][]string) [][]string {
	var podsUIDsFromExpressions [][]string
	for _, selector := range expressions {
		podsUIDsFromSelector, err := query.RetrievePodsUIDsByLabelsFilter(selector)
		if err == nil {
			podsUIDsFromExpressions = append(podsUIDsFromExpressions, podsUIDsFromSelector)
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
