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
	"github.com/vmware/purser/pkg/controller/utils"
)

// RetrieveDeploymentHierarchy returns hierarchy for a given deployment
func RetrieveDeploymentHierarchy(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for deployment, empty name is given")
		return JSONDataWrapper{}
	}
	query := `query {
		parent(func: has(isDeployment)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~deployment @filter(has(isReplicaset)) {
				name
				type
			}
		}
	}`
	return getJSONDataFromQuery(query)
}

// RetrieveDeploymentMetrics returns metrics for a given deployment
func RetrieveDeploymentMetrics(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for deployment, empty name is given")
		return JSONDataWrapper{}
	}
	secondsSinceMonthStart := fmt.Sprintf("%f", utils.GetSecondsSince(utils.GetCurrentMonthStartTime()))
	query := `query {
		dep as var(func: has(isDeployment)) @filter(eq(name, "` + name + `")) {
			~deployment @filter(has(isReplicaset)) {
				~replicaset @filter(has(isPod)) {
					replicasetPodCpu as cpuRequest
					replicasetPodMemory as memoryRequest
					replicasetPvcStorage as storageRequest
					replicasetPodST as startTime
					replicasetPodSTSeconds as math(since(replicasetPodST))
					replicasetPodSecondsSinceStart as math(cond(replicasetPodSTSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, replicasetPodSTSeconds))
					replicasetPodET as endTime
					replicasetPodIsTerminated as count(endTime)
					replicasetPodSecondsSinceEnd as math(cond(replicasetPodIsTerminated == 0, 0.0, since(replicasetPodET)))
					replicasetPodDurationInHours as math((replicasetPodSecondsSinceStart - replicasetPodSecondsSinceEnd) / 3600)
					replicasetPodCpuCost as math(replicasetPodCpu * replicasetPodDurationInHours * ` + defaultCPUCostPerCPUPerHour + `)
					replicasetPodMemoryCost as math(replicasetPodMemory * replicasetPodDurationInHours * ` + defaultMemCostPerGBPerHour + `)
					replicasetPvcStorageCost as math(replicasetPvcStorage * replicasetPodDurationInHours * ` + defaultStorageCostPerGBPerHour + `)
				}
				deploymentReplicasetCpu as sum(val(replicasetPodCpu))
				deploymentReplicasetMemory as sum(val(replicasetPodMemory))
				deploymentReplicasetStorage as sum(val(replicasetPvcStorage))
				deploymentReplicasetCpuCost as sum(val(replicasetPodCpuCost))
				deploymentReplicasetMemoryCost as sum(val(replicasetPodMemoryCost))
				deploymentReplicasetStorageCost as sum(val(replicasetPvcStorageCost))
			}
			deploymentCpu as sum(val(deploymentReplicasetCpu))
			deploymentMemory as sum(val(deploymentReplicasetMemory))
			deploymentStorage as sum(val(deploymentReplicasetStorage))
			deploymentCpuCost as sum(val(deploymentReplicasetCpuCost))
			deploymentMemoryCost as sum(val(deploymentReplicasetMemoryCost))
			deploymentStorageCost as sum(val(deploymentReplicasetStorageCost))
		}

		parent(func: uid(dep)) {
			name
			type
			children: ~deployment @filter(has(isReplicaset)) {
				name
				type
				cpu: val(deploymentReplicasetCpu)
				memory: val(deploymentReplicasetMemory)
				storage: val(deploymentReplicasetStorage)
				cpuCost: val(deploymentReplicasetCpuCost)
				memoryCost: val(deploymentReplicasetMemoryCost)
				storageCost: val(deploymentReplicasetStorageCost)
			}
			cpu: val(deploymentCpu)
			memory: val(deploymentMemory)
			storage: val(deploymentStorage)
			cpuCost: val(deploymentCpuCost)
			memoryCost: val(deploymentMemoryCost)
			storageCost: val(deploymentStorageCost)
		}
	}`
	return getJSONDataFromQuery(query)
}
