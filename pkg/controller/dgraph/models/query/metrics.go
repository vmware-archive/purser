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
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

// DaemonsetMetrics query
func getQueryForDaemonsetMetrics(name string) string {
	return `query {
		parent(func: has(isDaemonset)) @filter(eq(name, "` + name + `")) {
			children: ~daemonset @filter(has(isPod)) {
				` + getQueryForMetricsComputationWithAliasAndVariables("Pod") + `
			}
			` + getQueryForAggregatingChildMetricsWithAlias("Pod") + `
		}
	}`
}

// JobMetrics query
func getQueryForJobMetrics(name string) string {
	return `query {
		parent(func: has(isJob)) @filter(eq(name, "` + name + `")) {
			children: ~job @filter(has(isPod)) {
				` + getQueryForMetricsComputationWithAliasAndVariables("Pod") + `
			}
			` + getQueryForAggregatingChildMetricsWithAlias("Pod") + `
		}
	}`
}

// ReplicasetMetrics query
func getQueryForReplicasetMetrics(name string) string {
	return `query {
		parent(func: has(isReplicaset)) @filter(eq(name, "` + name + `")) {
			children: ~replicaset @filter(has(isPod)) {
				` + getQueryForMetricsComputationWithAliasAndVariables("Pod") + `
			}
			` + getQueryForAggregatingChildMetricsWithAlias("Pod") + `
		}
	}`
}

// StatefulsetMetrics query
func getQueryForStatefulsetMetrics(name string) string {
	return `query {
		parent(func: has(isStatefulset)) @filter(eq(name, "` + name + `")) {
			children: ~statefulset @filter(has(isPod)) {
				` + getQueryForMetricsComputationWithAliasAndVariables("Pod") + `
			}
			` + getQueryForAggregatingChildMetricsWithAlias("Pod") + `
		}
	}`
}

// DeploymentMetrics query
func getQueryForDeploymentMetrics(name string) string {
	return `query {
		dep as var(func: has(isDeployment)) @filter(eq(name, "` + name + `")) {
			~deployment @filter(has(isReplicaset)) {
				~replicaset @filter(has(isPod)) {
					` + getQueryForMetricsComputation("ReplicasetPod") + `
				}
				` + getQueryForAggregatingChildMetrics("DeploymentReplicaset", "ReplicasetPod") + `
			}
			` + getQueryForAggregatingChildMetrics("Deployment", "DeploymentReplicaset") + `
		}

		parent(func: uid(dep)) {
			children: ~deployment @filter(has(isReplicaset)) {
				` + getQueryFromSubQueryWithAlias("DeploymentReplicaset") + `
			}
			` + getQueryFromSubQueryWithAlias("Deployment") + `
		}
	}`
}

// PodMetrics query
func getQueryForPodMetrics(name, cpuPrice, memoryPrice string) string {
	return `query {
		parent(func: has(isPod)) @filter(eq(name, "` + name + `")) {
			children: ~pod @filter(has(isContainer)) {
				name
				type
				` + getQueryForTimeComputation("Container") + `
				cpu: cpu as cpuRequest
				memory: memory as memoryRequest
				cpuCost: math(cpu * durationInHoursContainer * ` + cpuPrice + `)
				memoryCost: math(memory * durationInHoursContainer * ` + memoryPrice + `)
			}
			name
			type
			cpu: podCpu as cpuRequest
			memory: podMemory as memoryRequest
			storage: pvcStorage as storageRequest
			` + getQueryForTimeComputation("") + `
			cpuCost: math(podCpu * durationInHours * ` + cpuPrice + `)
			memoryCost: math(podMemory * durationInHours * ` + memoryPrice + `)
			storageCost: math(pvcStorage * durationInHours * ` + models.DefaultStorageCostPerGBPerHour + `)
		}
	}`
}

// ContainerMetrics query
func getQueryForContainerMetrics(name string) string {
	return `query {
		parent(func: has(isContainer)) @filter(eq(name, "` + name + `")) {
			name
			type
			cpu: cpu as cpuRequest
			memory: memory as memoryRequest
			` + getQueryForTimeComputation("") + `
			cpuCost: math(cpu * durationInHours * ` + models.DefaultCPUCostPerCPUPerHour + `)
			memoryCost: math(memory * durationInHours * ` + models.DefaultMemCostPerGBPerHour + `)
		}
	}`
}

// PVMetrics query
func getQueryForPVMetrics(name string) string {
	return `query {
		parent(func: has(isPersistentVolume)) @filter(eq(name, "` + name + `")) {
			children: ~pv @filter(has(isPersistentVolumeClaim)) {
				name
				type
				storage: pvcStorage as storageCapacity
				` + getQueryForTimeComputation("PVC") + `
				storageCost: math(pvcStorage * durationInHoursPVC * ` + models.DefaultStorageCostPerGBPerHour + `)
			}
			name
			type
			storage: storage as storageCapacity
			` + getQueryForTimeComputation("") + `
			storageCost: math(storage * durationInHours * ` + models.DefaultStorageCostPerGBPerHour + `)
        }
    }`
}

// PVCMetrics query
func getQueryForPVCMetrics(name string) string {
	return `query {
		parent(func: has(isPersistentVolumeClaim)) @filter(eq(name, "` + name + `")) {
			name
			type
			storage: storage as storageCapacity
			` + getQueryForTimeComputation("") + `
			storageCost: math(storage * durationInHours * ` + models.DefaultStorageCostPerGBPerHour + `)
        }
    }`
}

// NodeMetrics query
func getQueryForNodeMetrics(name string) string {
	return `query {
		parent(func: has(isNode)) @filter(eq(name, "` + name + `")) {
			children: ~node @filter(has(isPod)) {
				` + getQueryForMetricsComputationWithAlias("Pod") + `
			}
			name
			type
			cpu: cpu as cpuCapacity
			memory: memory as memoryCapacity
			storage: storage as sum(val(storagePod))
			` + getQueryForTimeComputation("") + `
			` + getQueryForCostWithPriceWithAlias("") + `
		}
	}`
}

// NamespaceMetrics query
func getQueryForNamespaceMetrics(name string) string {
	return `query {
		ns as var(func: has(isNamespace)) @filter(eq(name, "` + name + `")) {
			childs as ~namespace @filter(has(isDeployment) OR has(isStatefulset) OR has(isJob) OR has(isDaemonset) OR (has(isReplicaset) AND (NOT has(deployment)))) {
				name
				type
				~deployment @filter(has(isReplicaset)) {
					name
					type
					~replicaset @filter(has(isPod)) {
						` + getQueryForMetricsComputation("ReplicasetPod") + `
			        }
					` + getQueryForAggregatingChildMetrics("DeploymentReplicaset", "ReplicasetPod") + `
                }
				~statefulset @filter(has(isPod)) {
					` + getQueryForMetricsComputation("StatefulsetPod") + `
                }
				~job @filter(has(isPod)) {
					` + getQueryForMetricsComputation("JobPod") + `
                }
				~daemonset @filter(has(isPod)) {
					` + getQueryForMetricsComputation("DaemonsetPod") + `
                }
				~replicaset @filter(has(isPod)) {
					` + getQueryForMetricsComputation("ReplicasetSimplePod") + `
                }
				` + getQueryForAggregatingChildMetrics("SumReplicasetSimplePod", "ReplicasetSimplePod") + `
				` + getQueryForAggregatingChildMetrics("SumDaemonsetPod", "DaemonsetPod") + `
				` + getQueryForAggregatingChildMetrics("SumJobPod", "JobPod") + `
				` + getQueryForAggregatingChildMetrics("SumStatefulsetPod", "StatefulsetPod") + `
				` + getQueryForAggregatingChildMetrics("SumDeploymentReplicaset", "DeploymentReplicaset") + `
				cpuNamespaceChild as math(cpu` + "SumReplicasetSimplePod" + ` + cpu` + "SumDaemonsetPod" + ` + cpu` + "SumJobPod" + ` + cpu` + "SumStatefulsetPod" + ` + cpu` + "SumDeploymentReplicaset" + `)
				memoryNamespaceChild as math(memory` + "SumReplicasetSimplePod" + ` + memory` + "SumDaemonsetPod" + ` + memory` + "SumJobPod" + ` + memory` + "SumStatefulsetPod" + ` + memory` + "SumDeploymentReplicaset" + `)
				storageNamespaceChild as math(storage` + "SumReplicasetSimplePod" + ` + storage` + "SumDaemonsetPod" + ` + storage` + "SumJobPod" + ` + storage` + "SumStatefulsetPod" + ` + storage` + "SumDeploymentReplicaset" + `)
				cpuCostNamespaceChild as math(cpuCost` + "SumReplicasetSimplePod" + ` + cpuCost` + "SumDaemonsetPod" + ` + cpuCost` + "SumJobPod" + ` + cpuCost` + "SumStatefulsetPod" + ` + cpuCost` + "SumDeploymentReplicaset" + `)
				memoryCostNamespaceChild as math(memoryCost` + "SumReplicasetSimplePod" + ` + memoryCost` + "SumDaemonsetPod" + ` + memoryCost` + "SumJobPod" + ` + memoryCost` + "SumStatefulsetPod" + ` + memoryCost` + "SumDeploymentReplicaset" + `)
				storageCostNamespaceChild as math(storageCost` + "SumReplicasetSimplePod" + ` + storageCost` + "SumDaemonsetPod" + ` + storageCost` + "SumJobPod" + ` + storageCost` + "SumStatefulsetPod" + ` + storageCost` + "SumDeploymentReplicaset" + `)
			}
			` + getQueryForAggregatingChildMetrics("Namespace", "NamespaceChild") + `
		}

		parent(func: uid(ns)) {
			children: ~namespace @filter(uid(childs)) {
				` + getQueryFromSubQueryWithAlias("NamespaceChild") + `
			}
			` + getQueryFromSubQueryWithAlias("Namespace") + `
        }
    }`
}

// LogicalResourcesMetrics query
func getQueryForLogicalResources() string {
	return `query {
			ns as var(func: has(isNamespace)) {
				~namespace @filter(has(isPod)){
					` + getQueryForMetricsComputation("NamespacePod") + `
				}
				` + getQueryForAggregatingChildMetrics("Namespace", "NamespacePod") + `
			}
	
			children(func: uid(ns)) {
				getQueryFromSubQueryWithAlias
				` + getQueryFromSubQueryWithAlias("Namespace") + `
			}
		}`
}

// PhysicalResourcesMetrics query
func getQueryForPhysicalResources() string {
	return `query {
			children(func: has(name)) @filter(has(isNode) OR has(isPersistentVolume)) {
				` + getQueryForMetricsComputationWithAlias("") + `
			}
		}`
}

/*
The following functions are related to Queries for Custom Groups
*/

func getQueryForAllGroupsData() string {
	return `query {
		groups(func: has(isGroup)) {
			name
			podsCount
			mtdCPU
			mtdMemory
			mtdStorage
			cpu
			memory
			storage
			mtdCPUCost
			mtdMemoryCost
			mtdStorageCost
			mtdCost
		}
	}`
}

func getQueryForGroupMetrics(podsUIDs string) string {
	return `query {
		var(func: uid(` + podsUIDs + `)) {
			podCpu as cpuRequest
			podMemory as memoryRequest
			pvcStorage as storageRequest
			podCpuLimit as cpuLimit
			podMemoryLimit as memoryLimit
			cpuRequestCount as count(cpuRequest)
			memoryRequestCount as count(memoryRequest)
			storageRequestCount as count(storageRequest)
			cpuLimitCount as count(cpuLimit)
			memoryLimitCount as count(memoryLimit)
			` + getQueryForTimeComputation("") + `
			isAlive as math(cond(isTerminated == 0, 1, 0))
			pitPodCPU as math(cond(isTerminated == 0, cond(cpuRequestCount > 0, podCpu, 0.0), 0.0))
			pitPodMemory as math(cond(isTerminated == 0, cond(memoryRequestCount > 0, podMemory, 0.0), 0.0))
			pitPvcStorage as math(cond(isTerminated == 0, cond(storageRequestCount > 0, pvcStorage, 0.0), 0.0))
			pitPodCPULimit as math(cond(isTerminated == 0, cond(cpuLimitCount > 0, podCpuLimit, 0.0), 0.0))
			pitPodMemoryLimit as math(cond(isTerminated == 0, cond(memoryLimitCount > 0, podMemoryLimit, 0.0), 0.0))
			mtdPodCPU as math(podCpu * durationInHours)
			mtdPodMemory as math(podMemory * durationInHours)
			mtdPvcStorage as math(pvcStorage * durationInHours)
			mtdPodCPULimit as math(podCpuLimit * durationInHours)
			mtdPodMemoryLimit as math(podMemoryLimit * durationInHours)
			pricePerCPU as cpuPrice
			pricePerMemory as memoryPrice
			podCpuCost as math(mtdPodCPU * pricePerCPU)
			podMemoryCost as math(mtdPodMemory * pricePerMemory)
			podStorageCost as math(mtdPvcStorage * ` + models.DefaultStorageCostPerGBPerHour + `)
		}
		
		group() {
			pitCPU: sum(val(pitPodCPU))
			pitMemory: sum(val(pitPodMemory))
			pitStorage: sum(val(pitPvcStorage))
			pitCPULimit: sum(val(pitPodCPULimit))
			pitMemoryLimit: sum(val(pitPodMemoryLimit))
			mtdCPU: sum(val(mtdPodCPU))
			mtdMemory: sum(val(mtdPodMemory))
			mtdStorage: sum(val(mtdPvcStorage))
			mtdCPULimit: sum(val(mtdPodCPULimit))
			mtdMemoryLimit: sum(val(mtdPodMemoryLimit))
			cpuCost: sum(val(podCpuCost))
			memoryCost: sum(val(podMemoryCost))
			storageCost: sum(val(podStorageCost))
			livePods: sum(val(isAlive))
		}
	}`
}
