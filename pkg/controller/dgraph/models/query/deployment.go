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

import "github.com/Sirupsen/logrus"

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

// RetrieveDeploymentMetrics returns hierarchy for a given deployment
func RetrieveDeploymentMetrics(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for deployment, empty name is given")
		return JSONDataWrapper{}
	}
	query := `query {
		dep as var(func: has(isDeployment)) @filter(eq(name, "` + name + `")) {
			~deployment @filter(has(isReplicaset)) {
				~replicaset @filter(has(isPod)) {
					replicasetPodCpu as cpuRequest
					replicasetPodMemory as memoryRequest
					replicasetPvcStorage as storageRequest
				}
				deploymentReplicasetCpu as sum(val(replicasetPodCpu))
				deploymentReplicasetMemory as sum(val(replicasetPodMemory))
				deploymentReplicasetStorage as sum(val(replicasetPvcStorage))
			}
			deploymentCpu as sum(val(deploymentReplicasetCpu))
			deploymentMemory as sum(val(deploymentReplicasetMemory))
			deploymentStorage as sum(val(deploymentReplicasetStorage))
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
				cpuCost: math(deploymentReplicasetCpu * ` + defaultCPUCostPerCPUPerHour + `)
				memoryCost: math(deploymentReplicasetMemory * ` + defaultMemCostPerGBPerHour + `)
				storageCost: math(deploymentReplicasetStorage * ` + defaultStorageCostPerGBPerHour + `)
			}
			cpu: val(deploymentCpu)
			memory: val(deploymentMemory)
			storage: val(deploymentStorage)
			cpuCost: math(deploymentCpu * ` + defaultCPUCostPerCPUPerHour + `)
			memoryCost: math(deploymentMemory * ` + defaultMemCostPerGBPerHour + `)
			storageCost: math(deploymentStorage * ` + defaultStorageCostPerGBPerHour + `)
		}
	}`
	return getJSONDataFromQuery(query)
}
