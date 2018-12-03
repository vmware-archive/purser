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
	"github.com/vmware/purser/pkg/controller/utils"
)

// RetrieveNamespaceHierarchy returns hierarchy for a given namespace
func RetrieveNamespaceHierarchy(name string) JSONDataWrapper {
	if name == All {
		return RetrieveClusterHierarchy(Logical)
	}

	query := `query {
		parent(func: has(isNamespace)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~namespace @filter(has(isDeployment) OR has(isStatefulset) OR has(isJob) OR has(isDaemonset) OR (has(isReplicaset) AND (NOT has(deployment)))) {
				name
				type
			}
        }
    }`
	return getJSONDataFromQuery(query)
}

// RetrieveNamespaceMetrics returns metrics for a given namespace
func RetrieveNamespaceMetrics(name string) JSONDataWrapper {
	if name == All {
		return RetrieveClusterHierarchy(Logical)
	}

	secondsSinceMonthStart := fmt.Sprintf("%f", utils.GetSecondsSince(utils.GetCurrentMonthStartTime()))
	query := `query {
		ns as var(func: has(isNamespace)) @filter(eq(name, "` + name + `")) {
			childs as ~namespace @filter(has(isDeployment) OR has(isStatefulset) OR has(isJob) OR has(isDaemonset) OR (has(isReplicaset) AND (NOT has(deployment)))) {
				name
				type
				~deployment @filter(has(isReplicaset)) {
					name
					type
					~replicaset @filter(has(isPod)) {
						name
						type
				        replicasetPodCpu as cpuRequest
				        replicasetPodMemory as memoryRequest
						replicasetPvcStorage as storageRequest
						replicasetPodST as startTime
						replicasetPodSTSeconds as math(since(replicasetPodST))
						replicasetPodSecondsSinceStart as math(cond(replicasetPodSTSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, replicasetPodSTSeconds))
						replicasetPodET as endTime
						replicasetPodIsTerminated as count(endTime)
						replicasetPodSecondsSinceEnd as math(cond(replicasetPodIsTerminated == 0, 0.0, since(replicasetPodET)))
						replicasetPodDurationInHours as math((replicasetPodSecondsSinceStart - replicasetPodSecondsSinceEnd) / 60)
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
				~statefulset @filter(has(isPod)) {
                    name
                    type
                    statefulsetPodCpu as cpuRequest
                    statefulsetPodMemory as memoryRequest
					statefulsetPvcStorage as storageRequest
					statefulsetPodST as startTime
					statefulsetPodSTSeconds as math(since(statefulsetPodST))
					statefulsetPodSecondsSinceStart as math(cond(statefulsetPodSTSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, statefulsetPodSTSeconds))
					statefulsetPodET as endTime
					statefulsetPodIsTerminated as count(endTime)
					statefulsetPodSecondsSinceEnd as math(cond(statefulsetPodIsTerminated == 0, 0.0, since(statefulsetPodET)))
					statefulsetPodDurationInHours as math((statefulsetPodSecondsSinceStart - statefulsetPodSecondsSinceEnd) / 60)
					statefulsetPodCpuCost as math(statefulsetPodCpu * statefulsetPodDurationInHours * ` + defaultCPUCostPerCPUPerHour + `)
					statefulsetPodMemoryCost as math(statefulsetPodMemory * statefulsetPodDurationInHours * ` + defaultMemCostPerGBPerHour + `)
					statefulsetPvcStorageCost as math(statefulsetPvcStorage * statefulsetPodDurationInHours * ` + defaultStorageCostPerGBPerHour + `)
                }
				~job @filter(has(isPod)) {
                    name
                    type
                    jobPodCpu as cpuRequest
                    jobPodMemory as memoryRequest
					jobPvcStorage as jobRequest
					jobPodST as startTime
					jobPodSTSeconds as math(since(jobPodST))
					jobPodSecondsSinceStart as math(cond(jobPodSTSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, jobPodSTSeconds))
					jobPodET as endTime
					jobPodIsTerminated as count(endTime)
					jobPodSecondsSinceEnd as math(cond(jobPodIsTerminated == 0, 0.0, since(jobPodET)))
					jobPodDurationInHours as math((jobPodSecondsSinceStart - jobPodSecondsSinceEnd) / 60)
					jobPodCpuCost as math(jobPodCpu * jobPodDurationInHours * ` + defaultCPUCostPerCPUPerHour + `)
					jobPodMemoryCost as math(jobPodMemory * jobPodDurationInHours * ` + defaultMemCostPerGBPerHour + `)
					jobPvcStorageCost as math(jobPvcStorage * jobPodDurationInHours * ` + defaultStorageCostPerGBPerHour + `)
                }
				~daemonset @filter(has(isPod)) {
                    name
                    type
                    daemonsetPodCpu as cpuRequest
                    daemonsetPodMemory as memoryRequest
					daemonsetPvcStorage as daemonsetRequest
					daemonsetPodST as startTime
					daemonsetPodSTSeconds as math(since(daemonsetPodST))
					daemonsetPodSecondsSinceStart as math(cond(daemonsetPodSTSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, daemonsetPodSTSeconds))
					daemonsetPodET as endTime
					daemonsetPodIsTerminated as count(endTime)
					daemonsetPodSecondsSinceEnd as math(cond(daemonsetPodIsTerminated == 0, 0.0, since(daemonsetPodET)))
					daemonsetPodDurationInHours as math((daemonsetPodSecondsSinceStart - daemonsetPodSecondsSinceEnd) / 60)
					daemonsetPodCpuCost as math(daemonsetPodCpu * daemonsetPodDurationInHours * ` + defaultCPUCostPerCPUPerHour + `)
					daemonsetPodMemoryCost as math(daemonsetPodMemory * daemonsetPodDurationInHours * ` + defaultMemCostPerGBPerHour + `)
					daemonsetPvcStorageCost as math(daemonsetPvcStorage * daemonsetPodDurationInHours * ` + defaultStorageCostPerGBPerHour + `)
                }
				~replicaset @filter(has(isPod)) {
                    name
                    type
                    replicasetSimplePodCpu as cpuRequest
                    replicasetSimplePodMemory as memoryRequest
					replicasetSimplePvcStorage as replicasetRequest
					replicasetSimplePodST as startTime
					replicasetSimplePodSTSeconds as math(since(replicasetSimplePodST))
					replicasetSimplePodSecondsSinceStart as math(cond(replicasetSimplePodSTSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, replicasetSimplePodSTSeconds))
					replicasetSimplePodET as endTime
					replicasetSimplePodIsTerminated as count(endTime)
					replicasetSimplePodSecondsSinceEnd as math(cond(replicasetSimplePodIsTerminated == 0, 0.0, since(replicasetSimplePodET)))
					replicasetSimplePodDurationInHours as math((replicasetSimplePodSecondsSinceStart - replicasetSimplePodSecondsSinceEnd) / 60)
					replicasetSimplePodCpuCost as math(replicasetSimplePodCpu * replicasetSimplePodDurationInHours * ` + defaultCPUCostPerCPUPerHour + `)
					replicasetSimplePodMemoryCost as math(replicasetSimplePodMemory * replicasetSimplePodDurationInHours * ` + defaultMemCostPerGBPerHour + `)
					replicasetSimplePvcStorageCost as math(replicasetSimplePvcStorage * replicasetSimplePodDurationInHours * ` + defaultStorageCostPerGBPerHour + `)
                }
				sumReplicasetSimplePodCpu as sum(val(replicasetSimplePodCpu))
				sumDaemonsetPodCpu as sum(val(daemonsetPodCpu))
				sumJobPodCpu as sum(val(jobPodCpu))
				sumStatefulsetPodCpu as sum(val(statefulsetPodCpu))
				sumDeploymentPodCpu as sum(val(deploymentReplicasetCpu))
				namespaceChildCpu as math(sumReplicasetSimplePodCpu + sumDaemonsetPodCpu + sumJobPodCpu + sumStatefulsetPodCpu + sumDeploymentPodCpu)

				sumReplicasetSimplePodMemory as sum(val(replicasetSimplePodMemory))
				sumDaemonsetPodMemory as sum(val(daemonsetPodMemory))
				sumJobPodMemory as sum(val(jobPodMemory))
				sumStatefulsetPodMemory as sum(val(statefulsetPodMemory))
				sumDeploymentPodMemory as sum(val(deploymentReplicasetMemory))
				namespaceChildMemory as math(sumReplicasetSimplePodMemory + sumDaemonsetPodMemory + sumJobPodMemory + sumStatefulsetPodMemory + sumDeploymentPodMemory)

				sumReplicasetSimplePvcStorage as sum(val(replicasetSimplePvcStorage))
				sumDaemonsetPvcStorage as sum(val(daemonsetPvcStorage))
				sumJobPvcStorage as sum(val(jobPvcStorage))
				sumStatefulsetPvcStorage as sum(val(statefulsetPvcStorage))
				sumDeploymentPvcStorage as sum(val(deploymentReplicasetStorage))
				namespaceChildStorage as math(sumReplicasetSimplePvcStorage + sumDaemonsetPvcStorage + sumJobPvcStorage + sumStatefulsetPvcStorage + sumDeploymentPvcStorage)
			
				sumReplicasetSimplePodCpuCost as sum(val(replicasetSimplePodCpuCost))
				sumDaemonsetPodCpuCost as sum(val(daemonsetPodCpuCost))
				sumJobPodCpuCost as sum(val(jobPodCpuCost))
				sumStatefulsetPodCpuCost as sum(val(statefulsetPodCpuCost))
				sumDeploymentPodCpuCost as sum(val(deploymentReplicasetCpuCost))
				namespaceChildCpuCost as math(sumReplicasetSimplePodCpuCost + sumDaemonsetPodCpuCost + sumJobPodCpuCost + sumStatefulsetPodCpuCost + sumDeploymentPodCpuCost)

				sumReplicasetSimplePodMemoryCost as sum(val(replicasetSimplePodMemoryCost))
				sumDaemonsetPodMemoryCost as sum(val(daemonsetPodMemoryCost))
				sumJobPodMemoryCost as sum(val(jobPodMemoryCost))
				sumStatefulsetPodMemoryCost as sum(val(statefulsetPodMemoryCost))
				sumDeploymentPodMemoryCost as sum(val(deploymentReplicasetMemoryCost))
				namespaceChildMemoryCost as math(sumReplicasetSimplePodMemoryCost + sumDaemonsetPodMemoryCost + sumJobPodMemoryCost + sumStatefulsetPodMemoryCost + sumDeploymentPodMemoryCost)

				sumReplicasetSimplePvcStorageCost as sum(val(replicasetSimplePvcStorageCost))
				sumDaemonsetPvcStorageCost as sum(val(daemonsetPvcStorageCost))
				sumJobPvcStorageCost as sum(val(jobPvcStorageCost))
				sumStatefulsetPvcStorageCost as sum(val(statefulsetPvcStorageCost))
				sumDeploymentPvcStorageCost as sum(val(deploymentReplicasetStorageCost))
				namespaceChildStorageCost as math(sumReplicasetSimplePvcStorageCost + sumDaemonsetPvcStorageCost + sumJobPvcStorageCost + sumStatefulsetPvcStorageCost + sumDeploymentPvcStorageCost)
			}
			namespaceCpu as sum(val(namespaceChildCpu))
			namespaceMemory as sum(val(namespaceChildMemory))
			namespaceStorage as sum(val(namespaceChildStorage))
			namespaceCpuCost as sum(val(namespaceChildCpuCost))
			namespaceMemoryCost as sum(val(namespaceChildMemoryCost))
			namespaceStorageCost as sum(val(namespaceChildStorageCost))
		}

		parent(func: uid(ns)) {
			name
            type
			children: ~namespace @filter(uid(childs)) {
				name
				type
				cpu: val(namespaceChildCpu)
				memory: val(namespaceChildMemory)
				storage: val(namespaceChildStorage)
				cpuCost: val(namespaceChildCpuCost)
				memoryCost: val(namespaceChildMemoryCost)
				storageCost: val(namespaceChildStorageCost)
			}
			cpu: val(namespaceCpu)
			memory: val(namespaceMemory)
			storage: val(namespaceStorage)
			cpuCost: val(namespaceCpuCost)
			memoryCost: val(namespaceMemoryCost)
			storageCost: val(namespaceStorageCost)
        }
    }`
	return getJSONDataFromQuery(query)
}

// getJSONDataFromQuery executes query and wraps the data in a desired structure(JSONDataWrapper)
func getJSONDataFromQuery(query string) JSONDataWrapper {
	parentRoot := ParentWrapper{}
	err := dgraph.ExecuteQuery(query, &parentRoot)
	if err != nil || len(parentRoot.Parent) == 0 {
		logrus.Errorf("Unable to execute query, err: (%v), length of output: (%d)", err, len(parentRoot.Parent))
		return JSONDataWrapper{}
	}
	root := JSONDataWrapper{
		Data: ParentWrapper{
			Name:        parentRoot.Parent[0].Name,
			Type:        parentRoot.Parent[0].Type,
			Children:    parentRoot.Parent[0].Children,
			CPU:         parentRoot.Parent[0].CPU,
			Memory:      parentRoot.Parent[0].Memory,
			Storage:     parentRoot.Parent[0].Storage,
			CPUCost:     parentRoot.Parent[0].CPUCost,
			MemoryCost:  parentRoot.Parent[0].MemoryCost,
			StorageCost: parentRoot.Parent[0].StorageCost,
		},
	}
	return root
}
