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
	"github.com/vmware/purser/pkg/controller/dgraph"
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
			        }
					deploymentReplicasetCpu as sum(val(replicasetPodCpu))
			        deploymentReplicasetMemory as sum(val(replicasetPodMemory))
					deploymentReplicasetStorage as sum(val(replicasetPvcStorage))
                }
				~statefulset @filter(has(isPod)) {
                    name
                    type
                    statefulsetPodCpu as cpuRequest
                    statefulsetPodMemory as memoryRequest
					statefulsetPvcStorage as storageRequest
                }
				~job @filter(has(isPod)) {
                    name
                    type
                    jobPodCpu as cpuRequest
                    jobPodMemory as memoryRequest
					jobPvcStorage as jobRequest
                }
				~daemonset @filter(has(isPod)) {
                    name
                    type
                    daemonsetPodCpu as cpuRequest
                    daemonsetPodMemory as memoryRequest
					daemonsetPvcStorage as daemonsetRequest
                }
				~replicaset @filter(has(isPod)) {
                    name
                    type
                    replicasetSimplePodCpu as cpuRequest
                    replicasetSimplePodMemory as memoryRequest
					replicasetSimplePvcStorage as replicasetRequest
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
			}
			namespaceCpu as sum(val(namespaceChildCpu))
			namespaceMemory as sum(val(namespaceChildMemory))
			namespaceStorage as sum(val(namespaceChildStorage))
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
				cpuCost: math(namespaceChildCpu * ` + defaultCPUCostPerCPUPerHour + `)
				memoryCost: math(namespaceChildMemory * ` + defaultMemCostPerGBPerHour + `)
				storageCost: math(namespaceChildStorage * ` + defaultStorageCostPerGBPerHour + `)
			}
			cpu: val(namespaceCpu)
			memory: val(namespaceMemory)
			storage: val(namespaceStorage)
			cpuCost: math(namespaceCpu * ` + defaultCPUCostPerCPUPerHour + `)
			memoryCost: math(namespaceMemory * ` + defaultMemCostPerGBPerHour + `)
			storageCost: math(namespaceStorage * ` + defaultStorageCostPerGBPerHour + `)
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
