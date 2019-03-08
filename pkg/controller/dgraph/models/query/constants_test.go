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

const (
	testSecondsSinceMonthStart = "1.45"
	testNodeName               = "node-minikube"
	testNamespaceName          = "namespace-default"
	testPodUIDList             = "0x3e283, 0x3e288"
	testDeploymentName         = "deployment-purser"
	testPodName                = "pod-purser-dgraph-0"
	testPVName                 = "pv-datadir-purser-dgraph"
	testPVCName                = "pvc-datadir-purser-dgraph"
	testContainerName          = "container-purser-controller"
	testJobName                = "job-purser"
	testDaemonsetName          = "daemonset-purser"
	testPodUID                 = "0x3e283"
	testPodXID                 = "purser:pod-purser-dgraph-0"
	testCPUPrice               = 0.24
	testMemoryPrice            = 0.1

	testHierarchy            = "hierarchy"
	testMetrics              = "metrics"
	testRetrieveAllGroups    = "retrieveAllGroups"
	testRetrieveGroupMetrics = "retrieveGroupMetrics"
	testRetrieveSubscribers  = "retrieveSubscribers"
	testLabelFilterPods      = "labelFilterPods"
	testAlivePods            = "alivePods"
	testPodInteractions      = "podInteractions"
	testWrongQuery           = "wrongQuery"
)

const deploymentMetricTestQuery = `query {
		dep as var(func: has(isDeployment)) @filter(eq(name, "deployment-purser")) {
			~deployment @filter(has(isReplicaset)) {
				~replicaset @filter(has(isPod)) {
					cpuReplicasetPod as cpuRequest
			memoryReplicasetPod as memoryRequest
			storageReplicasetPod as storageRequest
			stReplicasetPod as startTime
			stSecondsReplicasetPod as math(since(stReplicasetPod))
			secondsSinceStartReplicasetPod as math(cond(stSecondsReplicasetPod > 1.45, 1.45, stSecondsReplicasetPod))
			etReplicasetPod as endTime
			isTerminatedReplicasetPod as count(endTime)
			secondsSinceEndReplicasetPod as math(cond(isTerminatedReplicasetPod == 0, 0.0, since(etReplicasetPod)))
			durationInHoursReplicasetPod as math(cond(secondsSinceStartReplicasetPod > secondsSinceEndReplicasetPod, (secondsSinceStartReplicasetPod - secondsSinceEndReplicasetPod) / 3600, 0.0))
			pricePerCPUReplicasetPod as cpuPrice
			pricePerMemoryReplicasetPod as memoryPrice
			cpuCostReplicasetPod as math(cpuReplicasetPod * durationInHoursReplicasetPod * pricePerCPUReplicasetPod)
			memoryCostReplicasetPod as math(memoryReplicasetPod * durationInHoursReplicasetPod * pricePerMemoryReplicasetPod)
			storageCostReplicasetPod as math(storageReplicasetPod * durationInHoursReplicasetPod * 0.00013888888)
				}
				cpuDeploymentReplicaset as sum(val(cpuReplicasetPod))
			memoryDeploymentReplicaset as sum(val(memoryReplicasetPod))
			storageDeploymentReplicaset as sum(val(storageReplicasetPod))
			cpuCostDeploymentReplicaset as sum(val(cpuCostReplicasetPod))
			memoryCostDeploymentReplicaset as sum(val(memoryCostReplicasetPod))
			storageCostDeploymentReplicaset as sum(val(storageCostReplicasetPod))
			}
			cpuDeployment as sum(val(cpuDeploymentReplicaset))
			memoryDeployment as sum(val(memoryDeploymentReplicaset))
			storageDeployment as sum(val(storageDeploymentReplicaset))
			cpuCostDeployment as sum(val(cpuCostDeploymentReplicaset))
			memoryCostDeployment as sum(val(memoryCostDeploymentReplicaset))
			storageCostDeployment as sum(val(storageCostDeploymentReplicaset))
		}

		parent(func: uid(dep)) {
			children: ~deployment @filter(has(isReplicaset)) {
				name
			type
			cpu: val(cpuDeploymentReplicaset)
			memory: val(memoryDeploymentReplicaset)
			storage: val(storageDeploymentReplicaset)
			cpuCost: val(cpuCostDeploymentReplicaset)
			memoryCost: val(memoryCostDeploymentReplicaset)
			storageCost: val(storageCostDeploymentReplicaset)
			}
			name
			type
			cpu: val(cpuDeployment)
			memory: val(memoryDeployment)
			storage: val(storageDeployment)
			cpuCost: val(cpuCostDeployment)
			memoryCost: val(memoryCostDeployment)
			storageCost: val(storageCostDeployment)
		}
	}`

const groupMetricTestQuery = `query {
		var(func: uid(0x3e283, 0x3e288)) {
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
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > 1.45, 1.45, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math(cond(secondsSinceStart > secondsSinceEnd, (secondsSinceStart - secondsSinceEnd) / 3600, 0.0))
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
			podStorageCost as math(mtdPvcStorage * 0.00013888888)
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

const allGroupsDataTestQuery = `query {
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

const namespaceMetricTestQuery = `query {
		ns as var(func: has(isNamespace)) @filter(eq(name, "namespace-default")) {
			childs as ~namespace @filter(has(isDeployment) OR has(isStatefulset) OR has(isJob) OR has(isDaemonset) OR (has(isReplicaset) AND (NOT has(deployment)))) {
				name
				type
				~deployment @filter(has(isReplicaset)) {
					name
					type
					~replicaset @filter(has(isPod)) {
						cpuReplicasetPod as cpuRequest
			memoryReplicasetPod as memoryRequest
			storageReplicasetPod as storageRequest
			stReplicasetPod as startTime
			stSecondsReplicasetPod as math(since(stReplicasetPod))
			secondsSinceStartReplicasetPod as math(cond(stSecondsReplicasetPod > 1.45, 1.45, stSecondsReplicasetPod))
			etReplicasetPod as endTime
			isTerminatedReplicasetPod as count(endTime)
			secondsSinceEndReplicasetPod as math(cond(isTerminatedReplicasetPod == 0, 0.0, since(etReplicasetPod)))
			durationInHoursReplicasetPod as math(cond(secondsSinceStartReplicasetPod > secondsSinceEndReplicasetPod, (secondsSinceStartReplicasetPod - secondsSinceEndReplicasetPod) / 3600, 0.0))
			pricePerCPUReplicasetPod as cpuPrice
			pricePerMemoryReplicasetPod as memoryPrice
			cpuCostReplicasetPod as math(cpuReplicasetPod * durationInHoursReplicasetPod * pricePerCPUReplicasetPod)
			memoryCostReplicasetPod as math(memoryReplicasetPod * durationInHoursReplicasetPod * pricePerMemoryReplicasetPod)
			storageCostReplicasetPod as math(storageReplicasetPod * durationInHoursReplicasetPod * 0.00013888888)
			        }
					cpuDeploymentReplicaset as sum(val(cpuReplicasetPod))
			memoryDeploymentReplicaset as sum(val(memoryReplicasetPod))
			storageDeploymentReplicaset as sum(val(storageReplicasetPod))
			cpuCostDeploymentReplicaset as sum(val(cpuCostReplicasetPod))
			memoryCostDeploymentReplicaset as sum(val(memoryCostReplicasetPod))
			storageCostDeploymentReplicaset as sum(val(storageCostReplicasetPod))
                }
				~statefulset @filter(has(isPod)) {
					cpuStatefulsetPod as cpuRequest
			memoryStatefulsetPod as memoryRequest
			storageStatefulsetPod as storageRequest
			stStatefulsetPod as startTime
			stSecondsStatefulsetPod as math(since(stStatefulsetPod))
			secondsSinceStartStatefulsetPod as math(cond(stSecondsStatefulsetPod > 1.45, 1.45, stSecondsStatefulsetPod))
			etStatefulsetPod as endTime
			isTerminatedStatefulsetPod as count(endTime)
			secondsSinceEndStatefulsetPod as math(cond(isTerminatedStatefulsetPod == 0, 0.0, since(etStatefulsetPod)))
			durationInHoursStatefulsetPod as math(cond(secondsSinceStartStatefulsetPod > secondsSinceEndStatefulsetPod, (secondsSinceStartStatefulsetPod - secondsSinceEndStatefulsetPod) / 3600, 0.0))
			pricePerCPUStatefulsetPod as cpuPrice
			pricePerMemoryStatefulsetPod as memoryPrice
			cpuCostStatefulsetPod as math(cpuStatefulsetPod * durationInHoursStatefulsetPod * pricePerCPUStatefulsetPod)
			memoryCostStatefulsetPod as math(memoryStatefulsetPod * durationInHoursStatefulsetPod * pricePerMemoryStatefulsetPod)
			storageCostStatefulsetPod as math(storageStatefulsetPod * durationInHoursStatefulsetPod * 0.00013888888)
                }
				~job @filter(has(isPod)) {
					cpuJobPod as cpuRequest
			memoryJobPod as memoryRequest
			storageJobPod as storageRequest
			stJobPod as startTime
			stSecondsJobPod as math(since(stJobPod))
			secondsSinceStartJobPod as math(cond(stSecondsJobPod > 1.45, 1.45, stSecondsJobPod))
			etJobPod as endTime
			isTerminatedJobPod as count(endTime)
			secondsSinceEndJobPod as math(cond(isTerminatedJobPod == 0, 0.0, since(etJobPod)))
			durationInHoursJobPod as math(cond(secondsSinceStartJobPod > secondsSinceEndJobPod, (secondsSinceStartJobPod - secondsSinceEndJobPod) / 3600, 0.0))
			pricePerCPUJobPod as cpuPrice
			pricePerMemoryJobPod as memoryPrice
			cpuCostJobPod as math(cpuJobPod * durationInHoursJobPod * pricePerCPUJobPod)
			memoryCostJobPod as math(memoryJobPod * durationInHoursJobPod * pricePerMemoryJobPod)
			storageCostJobPod as math(storageJobPod * durationInHoursJobPod * 0.00013888888)
                }
				~daemonset @filter(has(isPod)) {
					cpuDaemonsetPod as cpuRequest
			memoryDaemonsetPod as memoryRequest
			storageDaemonsetPod as storageRequest
			stDaemonsetPod as startTime
			stSecondsDaemonsetPod as math(since(stDaemonsetPod))
			secondsSinceStartDaemonsetPod as math(cond(stSecondsDaemonsetPod > 1.45, 1.45, stSecondsDaemonsetPod))
			etDaemonsetPod as endTime
			isTerminatedDaemonsetPod as count(endTime)
			secondsSinceEndDaemonsetPod as math(cond(isTerminatedDaemonsetPod == 0, 0.0, since(etDaemonsetPod)))
			durationInHoursDaemonsetPod as math(cond(secondsSinceStartDaemonsetPod > secondsSinceEndDaemonsetPod, (secondsSinceStartDaemonsetPod - secondsSinceEndDaemonsetPod) / 3600, 0.0))
			pricePerCPUDaemonsetPod as cpuPrice
			pricePerMemoryDaemonsetPod as memoryPrice
			cpuCostDaemonsetPod as math(cpuDaemonsetPod * durationInHoursDaemonsetPod * pricePerCPUDaemonsetPod)
			memoryCostDaemonsetPod as math(memoryDaemonsetPod * durationInHoursDaemonsetPod * pricePerMemoryDaemonsetPod)
			storageCostDaemonsetPod as math(storageDaemonsetPod * durationInHoursDaemonsetPod * 0.00013888888)
                }
				~replicaset @filter(has(isPod)) {
					cpuReplicasetSimplePod as cpuRequest
			memoryReplicasetSimplePod as memoryRequest
			storageReplicasetSimplePod as storageRequest
			stReplicasetSimplePod as startTime
			stSecondsReplicasetSimplePod as math(since(stReplicasetSimplePod))
			secondsSinceStartReplicasetSimplePod as math(cond(stSecondsReplicasetSimplePod > 1.45, 1.45, stSecondsReplicasetSimplePod))
			etReplicasetSimplePod as endTime
			isTerminatedReplicasetSimplePod as count(endTime)
			secondsSinceEndReplicasetSimplePod as math(cond(isTerminatedReplicasetSimplePod == 0, 0.0, since(etReplicasetSimplePod)))
			durationInHoursReplicasetSimplePod as math(cond(secondsSinceStartReplicasetSimplePod > secondsSinceEndReplicasetSimplePod, (secondsSinceStartReplicasetSimplePod - secondsSinceEndReplicasetSimplePod) / 3600, 0.0))
			pricePerCPUReplicasetSimplePod as cpuPrice
			pricePerMemoryReplicasetSimplePod as memoryPrice
			cpuCostReplicasetSimplePod as math(cpuReplicasetSimplePod * durationInHoursReplicasetSimplePod * pricePerCPUReplicasetSimplePod)
			memoryCostReplicasetSimplePod as math(memoryReplicasetSimplePod * durationInHoursReplicasetSimplePod * pricePerMemoryReplicasetSimplePod)
			storageCostReplicasetSimplePod as math(storageReplicasetSimplePod * durationInHoursReplicasetSimplePod * 0.00013888888)
                }
				cpuSumReplicasetSimplePod as sum(val(cpuReplicasetSimplePod))
			memorySumReplicasetSimplePod as sum(val(memoryReplicasetSimplePod))
			storageSumReplicasetSimplePod as sum(val(storageReplicasetSimplePod))
			cpuCostSumReplicasetSimplePod as sum(val(cpuCostReplicasetSimplePod))
			memoryCostSumReplicasetSimplePod as sum(val(memoryCostReplicasetSimplePod))
			storageCostSumReplicasetSimplePod as sum(val(storageCostReplicasetSimplePod))
				cpuSumDaemonsetPod as sum(val(cpuDaemonsetPod))
			memorySumDaemonsetPod as sum(val(memoryDaemonsetPod))
			storageSumDaemonsetPod as sum(val(storageDaemonsetPod))
			cpuCostSumDaemonsetPod as sum(val(cpuCostDaemonsetPod))
			memoryCostSumDaemonsetPod as sum(val(memoryCostDaemonsetPod))
			storageCostSumDaemonsetPod as sum(val(storageCostDaemonsetPod))
				cpuSumJobPod as sum(val(cpuJobPod))
			memorySumJobPod as sum(val(memoryJobPod))
			storageSumJobPod as sum(val(storageJobPod))
			cpuCostSumJobPod as sum(val(cpuCostJobPod))
			memoryCostSumJobPod as sum(val(memoryCostJobPod))
			storageCostSumJobPod as sum(val(storageCostJobPod))
				cpuSumStatefulsetPod as sum(val(cpuStatefulsetPod))
			memorySumStatefulsetPod as sum(val(memoryStatefulsetPod))
			storageSumStatefulsetPod as sum(val(storageStatefulsetPod))
			cpuCostSumStatefulsetPod as sum(val(cpuCostStatefulsetPod))
			memoryCostSumStatefulsetPod as sum(val(memoryCostStatefulsetPod))
			storageCostSumStatefulsetPod as sum(val(storageCostStatefulsetPod))
				cpuSumDeploymentReplicaset as sum(val(cpuDeploymentReplicaset))
			memorySumDeploymentReplicaset as sum(val(memoryDeploymentReplicaset))
			storageSumDeploymentReplicaset as sum(val(storageDeploymentReplicaset))
			cpuCostSumDeploymentReplicaset as sum(val(cpuCostDeploymentReplicaset))
			memoryCostSumDeploymentReplicaset as sum(val(memoryCostDeploymentReplicaset))
			storageCostSumDeploymentReplicaset as sum(val(storageCostDeploymentReplicaset))
				cpuNamespaceChild as math(cpuSumReplicasetSimplePod + cpuSumDaemonsetPod + cpuSumJobPod + cpuSumStatefulsetPod + cpuSumDeploymentReplicaset)
				memoryNamespaceChild as math(memorySumReplicasetSimplePod + memorySumDaemonsetPod + memorySumJobPod + memorySumStatefulsetPod + memorySumDeploymentReplicaset)
				storageNamespaceChild as math(storageSumReplicasetSimplePod + storageSumDaemonsetPod + storageSumJobPod + storageSumStatefulsetPod + storageSumDeploymentReplicaset)
				cpuCostNamespaceChild as math(cpuCostSumReplicasetSimplePod + cpuCostSumDaemonsetPod + cpuCostSumJobPod + cpuCostSumStatefulsetPod + cpuCostSumDeploymentReplicaset)
				memoryCostNamespaceChild as math(memoryCostSumReplicasetSimplePod + memoryCostSumDaemonsetPod + memoryCostSumJobPod + memoryCostSumStatefulsetPod + memoryCostSumDeploymentReplicaset)
				storageCostNamespaceChild as math(storageCostSumReplicasetSimplePod + storageCostSumDaemonsetPod + storageCostSumJobPod + storageCostSumStatefulsetPod + storageCostSumDeploymentReplicaset)
			}
			cpuNamespace as sum(val(cpuNamespaceChild))
			memoryNamespace as sum(val(memoryNamespaceChild))
			storageNamespace as sum(val(storageNamespaceChild))
			cpuCostNamespace as sum(val(cpuCostNamespaceChild))
			memoryCostNamespace as sum(val(memoryCostNamespaceChild))
			storageCostNamespace as sum(val(storageCostNamespaceChild))
		}

		parent(func: uid(ns)) {
			children: ~namespace @filter(uid(childs)) {
				name
			type
			cpu: val(cpuNamespaceChild)
			memory: val(memoryNamespaceChild)
			storage: val(storageNamespaceChild)
			cpuCost: val(cpuCostNamespaceChild)
			memoryCost: val(memoryCostNamespaceChild)
			storageCost: val(storageCostNamespaceChild)
			}
			name
			type
			cpu: val(cpuNamespace)
			memory: val(memoryNamespace)
			storage: val(storageNamespace)
			cpuCost: val(cpuCostNamespace)
			memoryCost: val(memoryCostNamespace)
			storageCost: val(storageCostNamespace)
        }
    }`

const nodeMetricTestQuery = `query {
		parent(func: has(isNode)) @filter(eq(name, "node-minikube")) {
			children: ~node @filter(has(isPod)) {
				name
			type
			cpu: cpuPod as cpuRequest
			memory: memoryPod as memoryRequest
			storage: storagePod as storageRequest
			stPod as startTime
			stSecondsPod as math(since(stPod))
			secondsSinceStartPod as math(cond(stSecondsPod > 1.45, 1.45, stSecondsPod))
			etPod as endTime
			isTerminatedPod as count(endTime)
			secondsSinceEndPod as math(cond(isTerminatedPod == 0, 0.0, since(etPod)))
			durationInHoursPod as math(cond(secondsSinceStartPod > secondsSinceEndPod, (secondsSinceStartPod - secondsSinceEndPod) / 3600, 0.0))
			pricePerCPUPod as cpuPrice
			pricePerMemoryPod as memoryPrice
			cpuCost: math(cpuPod * durationInHoursPod * pricePerCPUPod)
			memoryCost: math(memoryPod * durationInHoursPod * pricePerMemoryPod)
			storageCost: math(storagePod * durationInHoursPod * 0.00013888888)
			}
			name
			type
			cpu: cpu as cpuCapacity
			memory: memory as memoryCapacity
			storage: storage as sum(val(storagePod))
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > 1.45, 1.45, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math(cond(secondsSinceStart > secondsSinceEnd, (secondsSinceStart - secondsSinceEnd) / 3600, 0.0))
			pricePerCPU as cpuPrice
			pricePerMemory as memoryPrice
			cpuCost: math(cpu * durationInHours * pricePerCPU)
			memoryCost: math(memory * durationInHours * pricePerMemory)
			storageCost: math(storage * durationInHours * 0.00013888888)
		}
	}`

const podMetricTestQuery = `query {
		parent(func: has(isPod)) @filter(eq(name, "pod-purser-dgraph-0")) {
			children: ~pod @filter(has(isContainer)) {
				name
				type
				stContainer as startTime
			stSecondsContainer as math(since(stContainer))
			secondsSinceStartContainer as math(cond(stSecondsContainer > 1.45, 1.45, stSecondsContainer))
			etContainer as endTime
			isTerminatedContainer as count(endTime)
			secondsSinceEndContainer as math(cond(isTerminatedContainer == 0, 0.0, since(etContainer)))
			durationInHoursContainer as math(cond(secondsSinceStartContainer > secondsSinceEndContainer, (secondsSinceStartContainer - secondsSinceEndContainer) / 3600, 0.0))
				cpu: cpu as cpuRequest
				memory: memory as memoryRequest
				cpuCost: math(cpu * durationInHoursContainer * 0.24000000000)
				memoryCost: math(memory * durationInHoursContainer * 0.10000000000)
			}
			name
			type
			cpu: cpuPod as cpuRequest
			memory: memoryPod as memoryRequest
			storage: storagePod as storageRequest
			stPod as startTime
			stSecondsPod as math(since(stPod))
			secondsSinceStartPod as math(cond(stSecondsPod > 1.45, 1.45, stSecondsPod))
			etPod as endTime
			isTerminatedPod as count(endTime)
			secondsSinceEndPod as math(cond(isTerminatedPod == 0, 0.0, since(etPod)))
			durationInHoursPod as math(cond(secondsSinceStartPod > secondsSinceEndPod, (secondsSinceStartPod - secondsSinceEndPod) / 3600, 0.0))
			pricePerCPUPod as cpuPrice
			pricePerMemoryPod as memoryPrice
			cpuCost: math(cpuPod * durationInHoursPod * pricePerCPUPod)
			memoryCost: math(memoryPod * durationInHoursPod * pricePerMemoryPod)
			storageCost: math(storagePod * durationInHoursPod * 0.00013888888)
		}
	}`

const pvMetricTestQuery = `query {
		parent(func: has(isPersistentVolume)) @filter(eq(name, "pv-datadir-purser-dgraph")) {
			children: ~pv @filter(has(isPersistentVolumeClaim)) {
				name
				type
				storage: pvcStorage as storageCapacity
				stPVC as startTime
			stSecondsPVC as math(since(stPVC))
			secondsSinceStartPVC as math(cond(stSecondsPVC > 1.45, 1.45, stSecondsPVC))
			etPVC as endTime
			isTerminatedPVC as count(endTime)
			secondsSinceEndPVC as math(cond(isTerminatedPVC == 0, 0.0, since(etPVC)))
			durationInHoursPVC as math(cond(secondsSinceStartPVC > secondsSinceEndPVC, (secondsSinceStartPVC - secondsSinceEndPVC) / 3600, 0.0))
				storageCost: math(pvcStorage * durationInHoursPVC * 0.00013888888)
			}
			name
			type
			storage: storage as storageCapacity
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > 1.45, 1.45, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math(cond(secondsSinceStart > secondsSinceEnd, (secondsSinceStart - secondsSinceEnd) / 3600, 0.0))
			storageCost: math(storage * durationInHours * 0.00013888888)
        }
    }`

const pvcMetricTestQuery = `query {
		parent(func: has(isPersistentVolumeClaim)) @filter(eq(name, "pvc-datadir-purser-dgraph")) {
			name
			type
			storage: storage as storageCapacity
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > 1.45, 1.45, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math(cond(secondsSinceStart > secondsSinceEnd, (secondsSinceStart - secondsSinceEnd) / 3600, 0.0))
			storageCost: math(storage * durationInHours * 0.00013888888)
        }
    }`

const containerMetricTestQuery = `query {
		parent(func: has(isContainer)) @filter(eq(name, "container-purser-controller")) {
			name
			type
			cpu: cpu as cpuRequest
			memory: memory as memoryRequest
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > 1.45, 1.45, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math(cond(secondsSinceStart > secondsSinceEnd, (secondsSinceStart - secondsSinceEnd) / 3600, 0.0))
			cpuCost: math(cpu * durationInHours * 0.024)
			memoryCost: math(memory * durationInHours * 0.01)
		}
	}`

const logicalResourcesMetricTestQuery = `query {
			ns as var(func: has(isNamespace)) {
				~namespace @filter(has(isPod)){
					cpuNamespacePod as cpuRequest
			memoryNamespacePod as memoryRequest
			storageNamespacePod as storageRequest
			stNamespacePod as startTime
			stSecondsNamespacePod as math(since(stNamespacePod))
			secondsSinceStartNamespacePod as math(cond(stSecondsNamespacePod > 1.45, 1.45, stSecondsNamespacePod))
			etNamespacePod as endTime
			isTerminatedNamespacePod as count(endTime)
			secondsSinceEndNamespacePod as math(cond(isTerminatedNamespacePod == 0, 0.0, since(etNamespacePod)))
			durationInHoursNamespacePod as math(cond(secondsSinceStartNamespacePod > secondsSinceEndNamespacePod, (secondsSinceStartNamespacePod - secondsSinceEndNamespacePod) / 3600, 0.0))
			pricePerCPUNamespacePod as cpuPrice
			pricePerMemoryNamespacePod as memoryPrice
			cpuCostNamespacePod as math(cpuNamespacePod * durationInHoursNamespacePod * pricePerCPUNamespacePod)
			memoryCostNamespacePod as math(memoryNamespacePod * durationInHoursNamespacePod * pricePerMemoryNamespacePod)
			storageCostNamespacePod as math(storageNamespacePod * durationInHoursNamespacePod * 0.00013888888)
				}
				cpuNamespace as sum(val(cpuNamespacePod))
			memoryNamespace as sum(val(memoryNamespacePod))
			storageNamespace as sum(val(storageNamespacePod))
			cpuCostNamespace as sum(val(cpuCostNamespacePod))
			memoryCostNamespace as sum(val(memoryCostNamespacePod))
			storageCostNamespace as sum(val(storageCostNamespacePod))
			}
	
			children(func: uid(ns)) {
				name
			type
			cpu: val(cpuNamespace)
			memory: val(memoryNamespace)
			storage: val(storageNamespace)
			cpuCost: val(cpuCostNamespace)
			memoryCost: val(memoryCostNamespace)
			storageCost: val(storageCostNamespace)
			}
		}`

const phycialResourcesMetricTestQuery = `query {
			children(func: has(name)) @filter(has(isNode) OR has(isPersistentVolume)) {
				name
			type
			cpu: cpu as cpuRequest
			memory: memory as memoryRequest
			storage: storage as storageRequest
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > 1.45, 1.45, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math(cond(secondsSinceStart > secondsSinceEnd, (secondsSinceStart - secondsSinceEnd) / 3600, 0.0))
			pricePerCPU as cpuPrice
			pricePerMemory as memoryPrice
			cpuCost: math(cpu * durationInHours * pricePerCPU)
			memoryCost: math(memory * durationInHours * pricePerMemory)
			storageCost: math(storage * durationInHours * 0.00013888888)
			}
		}`

const testQueryForMetricsComputationWithAliasAndVariables = `name
			type
			cpu: cpuPod as cpuRequest
			memory: memoryPod as memoryRequest
			storage: storagePod as storageRequest
			stPod as startTime
			stSecondsPod as math(since(stPod))
			secondsSinceStartPod as math(cond(stSecondsPod > 1.45, 1.45, stSecondsPod))
			etPod as endTime
			isTerminatedPod as count(endTime)
			secondsSinceEndPod as math(cond(isTerminatedPod == 0, 0.0, since(etPod)))
			durationInHoursPod as math(cond(secondsSinceStartPod > secondsSinceEndPod, (secondsSinceStartPod - secondsSinceEndPod) / 3600, 0.0))
			pricePerCPUPod as cpuPrice
			pricePerMemoryPod as memoryPrice
			cpuCost: cpuCostPod as math(cpuPod * durationInHoursPod * pricePerCPUPod)
			memoryCost: memoryCostPod as math(memoryPod * durationInHoursPod * pricePerMemoryPod)
			storageCost: storageCostPod as math(storagePod * durationInHoursPod * 0.00013888888)`

const testQueryForAggregatingChildMetricsWithAlias = `name
			type
			cpu: sum(val(cpuPod))
			memory: sum(val(memoryPod))
			storage: sum(val(storagePod))
			cpuCost: sum(val(cpuCostPod))
			memoryCost: sum(val(memoryCostPod))
			storageCost: sum(val(storageCostPod))`

const testQueryForPodParentMetrics = `query {
		parent(func: has(isJob)) @filter(eq(name, "job-purser")) {
			children: ~job @filter(has(isPod)) {
				name
			type
			cpu: cpuPod as cpuRequest
			memory: memoryPod as memoryRequest
			storage: storagePod as storageRequest
			stPod as startTime
			stSecondsPod as math(since(stPod))
			secondsSinceStartPod as math(cond(stSecondsPod > 1.45, 1.45, stSecondsPod))
			etPod as endTime
			isTerminatedPod as count(endTime)
			secondsSinceEndPod as math(cond(isTerminatedPod == 0, 0.0, since(etPod)))
			durationInHoursPod as math(cond(secondsSinceStartPod > secondsSinceEndPod, (secondsSinceStartPod - secondsSinceEndPod) / 3600, 0.0))
			pricePerCPUPod as cpuPrice
			pricePerMemoryPod as memoryPrice
			cpuCost: cpuCostPod as math(cpuPod * durationInHoursPod * pricePerCPUPod)
			memoryCost: memoryCostPod as math(memoryPod * durationInHoursPod * pricePerMemoryPod)
			storageCost: storageCostPod as math(storagePod * durationInHoursPod * 0.00013888888)
			}
			name
			type
			cpu: sum(val(cpuPod))
			memory: sum(val(memoryPod))
			storage: sum(val(storagePod))
			cpuCost: sum(val(cpuCostPod))
			memoryCost: sum(val(memoryCostPod))
			storageCost: sum(val(storageCostPod))
		}
	}`

const testQueryForHierarchy = `query {
		parent(func: has(isNode)) @filter(eq(name, "node-minikube")) {
			name
			type
			children: ~node @filter(has(isPod)) {
				name
				type
			}
		}
	}`
