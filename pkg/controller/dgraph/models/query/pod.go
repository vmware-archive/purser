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
	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

type podRoot struct {
	Pods []models.Pod `json:"pod"`
}

// RetrieveAllLivePods will return all pods without endTime in dgraph. Error is returned if any
// failure is encountered in the process.
func RetrieveAllLivePods() []models.Pod {
	query := getAllLivePodsQuery()
	newRoot := podRoot{}
	err := executeQuery(query, &newRoot)
	if err != nil {
		logrus.Errorf("unable to retrieve all live pods: %v", err)
		return nil
	}
	return newRoot.Pods
}

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

	result, err := executeQueryRaw(query)
	if err != nil {
		logrus.Errorf("Error while retrieving query for pods interactions. Name: (%v), isOrphan: (%v), error: (%v)", name, isOrphan, err)
		return nil
	}
	return result
}

func getPricePerResourceForPod(name string) (float64, float64) {
	query := `query {
		pod(func: has(isPod)) @filter(eq(name, "` + name + `")) {
			cpuPrice
			memoryPrice
		}
	}`
	newRoot := podRoot{}
	err := executeQuery(query, &newRoot)
	if err != nil || len(newRoot.Pods) < 1 {
		logrus.Errorf("err: %v", err)
		return models.DefaultCPUCostInFloat64, models.DefaultMemCostInFloat64
	}
	pod := newRoot.Pods[0]
	return pod.CPUPrice, pod.MemoryPrice
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
	err := executeQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Pods, nil
}

// RetrievePodsUIDsByLabelsFilter returns pods satisfying the filter conditions for labels (OR logic only)
func RetrievePodsUIDsByLabelsFilter(labels map[string][]string) ([]string, error) {
	labelFilter := createFilterFromListOfLabels(labels)
	q := getQueryForPodsWithLabelFilter(labelFilter)
	newRoot := podRoot{}
	err := executeQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return removeDuplicates(newRoot.Pods), nil
}

func removeDuplicates(pods []models.Pod) []string {
	duplicateChecker := make(map[string]bool)
	var podsUIDs []string
	for _, pod := range pods {
		if _, isPresent := duplicateChecker[pod.UID]; !isPresent {
			podsUIDs = append(podsUIDs, pod.UID)
			duplicateChecker[pod.UID] = true
		}
	}
	return podsUIDs
}
