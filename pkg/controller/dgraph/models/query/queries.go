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
			` + getQueryForMetricsComputationWithAlias("Pod") + `
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
			storageCapacity
			` + getQueryForTimeComputation("") + `
			storageCost: math(storage * durationInHours * ` + models.DefaultStorageCostPerGBPerHour + `)
			storageAllocated: sum(val(pvcStorage))
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
			cpuAllocated: sum(val(cpuPod))
			memoryAllocated: sum(val(memoryPod))
			cpuCapacity
			memoryCapacity
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
func getMetricsQueryForLogicalResources() string {
	return `query {
			ns as var(func: has(isNamespace)) {
				~namespace @filter(has(isPod) AND (NOT has(endTime))) {
					` + getQueryForMetricsComputation("NamespacePod") + `
				}
				` + getQueryForAggregatingChildMetrics("Namespace", "NamespacePod") + `
			}
	
			children(func: uid(ns)) {
				` + getQueryFromSubQueryWithAlias("Namespace") + `
			}
		}`
}

// PhysicalResourcesMetrics query
func getMetricsQueryForPhysicalResources() string {
	return `query {
			children(func: has(name)) @filter((has(isNode) OR has(isPersistentVolume)) AND (NOT has(endTime))) {
				name
			type
			cpu: cpu as cpuCapacity
			memory: memory as memoryCapacity
			storage: storage as storageCapacity
			` + getQueryForTimeComputation("") + `
			` + getQueryForCostWithPriceWithAlias("") + `
			}
		}`
}

// LogicalResourcesHierarchy query
func getHierarchyQueryForLogicalResource() string {
	return `query {
			children(func: has(isNamespace)) {
				name
				type
			}
		}`
}

// PhysicalResourcesHierarchy query
func getHierarchyQueryForPhysicalResource() string {
	return `query {
			children(func: has(name)) @filter(has(isNode) OR has(isPersistentVolume)) {
				name
				type
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
			projectedCPUCost
			projectedMemoryCost
			projectedStorageCost
			projectedCost
			lastMonthCPUCost
			lastMonthMemoryCost
			lastMonthStorageCost
			lastMonthCost
			lastLastMonthCPUCost
			lastLastMonthMemoryCost
			lastLastMonthStorageCost
			lastLastMonthCost
		}
	}`
}

func getQueryForGroupMetrics(podsUIDs string) string {
	secondsSince := getSecondsSinceForOtherMonths()
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
			podEndTime as endTime
			isTerminated as count(endTime)
			secondsSincePodEndTime as math(cond(isTerminated == 0, 0.0, since(podEndTime)))
			podStartTime as startTime
			secondsSincePodStartTime as math(since(podStartTime))
			secondsSinceCurrentMonthTrueStart as math(cond(secondsSincePodStartTime > ` + secondsSince["currentMonthStart"] + `, ` + secondsSince["currentMonthStart"] + `, secondsSincePodStartTime))
			currentMonthTrueDurationInHours as math(cond(secondsSinceCurrentMonthTrueStart > secondsSincePodEndTime, (secondsSinceCurrentMonthTrueStart - secondsSincePodEndTime)/3600, 0.0))
			secondsSinceLastMonthTrueStart as math(cond(secondsSincePodStartTime > ` + secondsSince["lastMonthStart"] + `, ` + secondsSince["lastMonthStart"] + `, secondsSincePodStartTime))
			secondsSinceLastMonthTrueEnd as math(cond(secondsSincePodEndTime > ` + secondsSince["lastMonthEnd"] + `, secondsSincePodEndTime, ` + secondsSince["lastMonthEnd"] + `))
			lastMonthTrueDurationInHours as math(cond(secondsSinceLastMonthTrueStart > secondsSinceLastMonthTrueEnd, (secondsSinceLastMonthTrueStart - secondsSinceLastMonthTrueEnd)/3600, 0.0))
			secondsSinceLastLastMonthTrueStart as math(cond(secondsSincePodStartTime > ` + secondsSince["lastLastMonthStart"] + `, ` + secondsSince["lastLastMonthStart"] + `, secondsSincePodStartTime))
			secondsSinceLastLastMonthTrueEnd as math(cond(secondsSincePodEndTime > ` + secondsSince["lastLastMonthEnd"] + `, secondsSincePodEndTime, ` + secondsSince["lastLastMonthEnd"] + `))
			lastLastMonthTrueDurationInHours as math(cond(secondsSinceLastLastMonthTrueStart > secondsSinceLastLastMonthTrueEnd, (secondsSinceLastLastMonthTrueStart - secondsSinceLastLastMonthTrueEnd)/3600, 0.0))
			isAlive as math(cond(isTerminated == 0, 1, 0))
			pitPodCPU as math(cond(isTerminated == 0, cond(cpuRequestCount > 0, podCpu, 0.0), 0.0))
			pitPodMemory as math(cond(isTerminated == 0, cond(memoryRequestCount > 0, podMemory, 0.0), 0.0))
			pitPvcStorage as math(cond(isTerminated == 0, cond(storageRequestCount > 0, pvcStorage, 0.0), 0.0))
			pitPodCPULimit as math(cond(isTerminated == 0, cond(cpuLimitCount > 0, podCpuLimit, 0.0), 0.0))
			pitPodMemoryLimit as math(cond(isTerminated == 0, cond(memoryLimitCount > 0, podMemoryLimit, 0.0), 0.0))
			mtdPodCPU as math(podCpu * currentMonthTrueDurationInHours)
			mtdPodMemory as math(podMemory * currentMonthTrueDurationInHours)
			mtdPvcStorage as math(pvcStorage * currentMonthTrueDurationInHours)
			mtdPodCPULimit as math(podCpuLimit * currentMonthTrueDurationInHours)
			mtdPodMemoryLimit as math(podMemoryLimit * currentMonthTrueDurationInHours)
			pricePerCPU as cpuPrice
			pricePerMemory as memoryPrice
			podCpuCost as math(mtdPodCPU * pricePerCPU)
			podMemoryCost as math(mtdPodMemory * pricePerMemory)
			podStorageCost as math(mtdPvcStorage * ` + models.DefaultStorageCostPerGBPerHour + `)
			podLiveCPUCostPerHour as math(pitPodCPU * pricePerCPU)
			podLiveMemoryCostPerHour as math(pitPodMemory * pricePerMemory)
			podLiveStorageCostPerHour as math(pitPvcStorage * ` + models.DefaultStorageCostPerGBPerHour + `)
			podCPUCostPerHour as math(podCpu * pricePerCPU)
			podMemoryCostPerHour as math(podMemory * pricePerMemory)
			podStorageCostPerHour as math(pvcStorage * ` + models.DefaultStorageCostPerGBPerHour + `)
			podCPUCostLastMonth as math(podCPUCostPerHour * lastMonthTrueDurationInHours)
			podMemoryCostLastMonth as math(podMemoryCostPerHour * lastMonthTrueDurationInHours)
			podStorageCostLastMonth as math(podStorageCostPerHour * lastMonthTrueDurationInHours)
			podCPUCostLastLastMonth as math(podCPUCostPerHour * lastLastMonthTrueDurationInHours)
			podMemoryCostLastLastMonth as math(podMemoryCostPerHour * lastLastMonthTrueDurationInHours)
			podStorageCostLastLastMonth as math(podStorageCostPerHour * lastLastMonthTrueDurationInHours)
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
			cpuCostPerHour: sum(val(podLiveCPUCostPerHour))
			memoryCostPerHour: sum(val(podLiveMemoryCostPerHour))
			storageCostPerHour: sum(val(podLiveStorageCostPerHour))
			lastMonthCPUCost: sum(val(podCPUCostLastMonth))
			lastMonthMemoryCost: sum(val(podMemoryCostLastMonth))
			lastMonthStorageCost: sum(val(podStorageCostLastMonth))
			lastLastMonthCPUCost: sum(val(podCPUCostLastLastMonth))
			lastLastMonthMemoryCost: sum(val(podMemoryCostLastLastMonth))
			lastLastMonthStorageCost: sum(val(podStorageCostLastLastMonth))
			livePods: sum(val(isAlive))
		}
	}`
}

func getQueryForSubscribersRetrieval() string {
	return `query {
		subscribers(func: has(isSubscriber)) @filter(NOT(has(endTime))) {
			name
			spec {
				headers
				url
			}
		}
	}`
}

func getAllLivePodsQuery() string {
	return `query {
		pods(func: has(isPod)) @filter(NOT has(endTime)) {
			uid
			xid
			name
		}
	}`
}

func getQueryForPodsWithLabelFilter(labelFilter string) string {
	return `query {
		var(func: has(isLabel)) @filter(` + labelFilter + `) {
            podUIDs as ~label @filter(has(isPod)) {
				name
			}
		}
		pods(func: uid(podUIDs)) {
			uid
			name
		}
	}`
}
