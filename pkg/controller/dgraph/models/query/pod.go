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

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
	"github.com/vmware/purser/pkg/controller/utils"
)

// RetrievePodsInteractions returns inbound and outbound interactions of a pod
func RetrievePodsInteractions(name string, isOrphan bool) []byte {
	var query string
	if name == All {
		if isOrphan {
			query = `query {
				pods(func: has(isPod)) {
					name
					outbound: pod {
						name
					}
					inbound: ~pod @filter(has(isPod)) {
						name
					}
				}
			}`
		} else {
			query = `query {
				pods(func: has(isPod)) @filter(has(pod)) {
					name
					outbound: pod {
						name
					}
					inbound: ~pod @filter(has(isPod)) {
						name
					}
				}
			}`
		}
	} else {
		query = `query {
			pods(func: has(isPod)) @filter(eq(name, "` + name + `")) {
				name
				outbound: pod {
					name
				}
				inbound: ~pod @filter(has(isPod)) {
					name
				}
			}
		}`
	}

	result, err := dgraph.ExecuteQueryRaw(query)
	if err != nil {
		logrus.Errorf("Error while retrieving query for pods interactions. Name: (%v), isOrphan: (%v), error: (%v)", name, isOrphan, err)
		return nil
	}
	return result
}

// RetrievePodHierarchy returns hierarchy for a given pod
func RetrievePodHierarchy(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for pod, empty name is given")
		return JSONDataWrapper{}
	}
	query := `query {
		parent(func: has(isPod)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~pod @filter(has(isContainer)) {
				name
				type
			}
		}
	}`
	return getJSONDataFromQuery(query)
}

// RetrievePodMetrics returns metrics for a given pod
func RetrievePodMetrics(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for pod, empty name is given")
		return JSONDataWrapper{}
	}
	secondsSinceMonthStart := fmt.Sprintf("%f", utils.GetSecondsSince(utils.GetCurrentMonthStartTime()))
	query := `query {
		parent(func: has(isPod)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~pod @filter(has(isContainer)) {
				name
				type
				stChild as startTime
				stSecondsChild as math(since(stChild))
				secondsSinceStartChild as math(cond(stSecondsChild > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, stSecondsChild))
				etChild as endTime
				isTerminatedChild as count(endTime)
				secondsSinceEndChild as math(cond(isTerminatedChild == 0, 0.0, since(etChild)))
				durationInHoursChild as math((secondsSinceStartChild - secondsSinceEndChild) / 3600)
				cpu: cpu as cpuRequest
				memory: memory as memoryRequest
				cpuCost: math(cpu * durationInHoursChild * ` + defaultCPUCostPerCPUPerHour + `)
				memoryCost: math(memory * durationInHoursChild * ` + defaultMemCostPerGBPerHour + `)
			}
			cpu: podCpu as cpuRequest
			memory: podMemory as memoryRequest
			storage: pvcStorage as storageRequest
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math((secondsSinceStart - secondsSinceEnd) / 3600)
			cpuCost: math(podCpu * durationInHours * ` + defaultCPUCostPerCPUPerHour + `)
			memoryCost: math(podMemory * durationInHours * ` + defaultMemCostPerGBPerHour + `)
			storageCost: math(pvcStorage * durationInHours * ` + defaultStorageCostPerGBPerHour + `)
		}
	}`
	return getJSONDataFromQuery(query)
}

// RetrievePodsInteractionsForAllLivePodsWithCount returns all pods in the dgraph
func RetrievePodsInteractionsForAllLivePodsWithCount() ([]models.Pod, error) {
	q := `query {
		pods(func: has(isPod)) @filter((NOT has(endTime))) {
			name
			pod {
				name
				count
			}
			cid: ~pod @filter(has(isService)) {
				name
			}
		}
	}`

	type root struct {
		Pods []models.Pod `json:"pods"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Pods, nil
}

// RetrievePodsByLabelsFilter returns pods satisfying the filter conditions for labels (OR logic only)
func RetrievePodsByLabelsFilter(labels map[string]string) ([]models.Pod, error) {
	labelFilter := createLabelFilter(labels)
	secondsSinceMonthStart := fmt.Sprintf("%f", utils.GetSecondsSince(utils.GetCurrentMonthStartTime()))
	q := `query {
		var(func: has(isLabel)) @filter(` + labelFilter + `) {
            podUIDs as ~label @filter(has(isPod)) {
				name
			}
		}
		pods(func: uid(podUIDs)) {
			uid
			name
			type
			cpu: podCpu as cpuRequest
			memory: podMemory as memoryRequest
			storage: pvcStorage as storageRequest
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math((secondsSinceStart - secondsSinceEnd) / 3600)
			cpuCost: math(podCpu * durationInHours * ` + defaultCPUCostPerCPUPerHour + `)
			memoryCost: math(podMemory * durationInHours * ` + defaultMemCostPerGBPerHour + `)
			storageCost: math(pvcStorage * durationInHours * ` + defaultStorageCostPerGBPerHour + `)
		}
	}`
	type root struct {
		Pods []models.Pod `json:"pods"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Pods, nil
}
