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

package main

import (
	"fmt"
	"kuber/metrics"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/api/core/v1"
)

func calculateCost(pods []*Pod, nodes map[string]*Node, pvcs map[string]*PersistentVolumeClaim) []*Pod {
	i := 0
	for i <= len(pods)-1 {
		node := nodes[pods[i].nodeName]
		pods[i].nodeCostPercentage = (float64)(node.getPodResourcePercentage(pods[i].name))
		totalComputeCost, cpuCost, memoryCost := getMonthToDateCostForInstanceType(node.instanceType)
		totalStorageCost := 0.0
		for _, pvc := range pods[i].pvcs {
			if pvcs[*pvc] != nil {
				storagePrice := getMonthToDateCostForStorageClass(*pvcs[*pvc].storageClass)
				totalStorageCost = totalStorageCost + storagePrice*pvcs[*pvc].capacityAllotedInGB
			} else {
				fmt.Printf("Persistent volume claim is not present for %s\n", *pvc)
			}
		}
		podCost := Cost{}
		podCost.totalCost = pods[i].nodeCostPercentage*totalComputeCost + totalStorageCost
		podCost.cpuCost = pods[i].nodeCostPercentage * cpuCost
		podCost.memoryCost = pods[i].nodeCostPercentage * memoryCost
		podCost.storageCost = totalStorageCost
		pods[i].cost = &podCost
		i++
	}
	return pods
}

func getPodsCostForLabel(label string) {
	pods := getPodsForLabelThroughClient(label)
	pods = getPodsCost(pods)
	printPodsVerbose(pods)
}

/*func getClusterPods() {
	pods := GetClusterPods()
	fmt.Printf("Total number of pods = %d\n", len(pods))
	//printPodsVerbose(pods)
}*/

func getClusterSummary() {
	pods := GetClusterPods()
	podMetrics := metrics.CalculatePodStatsFromContainers(pods)
	//fmt.Printf("%-30s%s\n", "Pod Name:", pods[i].name)
	fmt.Printf("Cluster Summary\n")

	fmt.Println("Compute:")
	nodes := GetClusterNodes()
	fmt.Printf("   %-25s   %d\n", "Node count:", len(nodes))
	var computeCost = 0.0
	for _, node := range nodes {
		instanceType := GetNodeType(node)
		total, _, _ := getMonthToDateCostForInstanceType(instanceType)
		computeCost = computeCost + total
	}
	fmt.Printf("   %-25s   %.2f$\n", "Cost:", computeCost)
	nodeMetrics := metrics.CalculateNodeStats(nodes)
	fmt.Printf("   Total Capacity:\n")
	fmt.Printf("      %-25s%d\n","Cpu(vCPU):", nodeMetrics.CpuLimit.Value())
	fmt.Printf("      %-25s%.2f\n","Memory(GB):", bytesToGB(nodeMetrics.MemoryLimit.Value()))


	fmt.Printf("   Provisioned Resources:\n")
	//fmt.Printf("      %-25s%d\n", "Cpu Limit(vCPU):", podMetrics.CpuLimit.Value())
	//fmt.Printf("      %-25s%.2f\n", "Memory Limit(GB):", bytesToGB(podMetrics.MemoryLimit.Value()))
	fmt.Printf("      %-25s%d\n", "Cpu Request(vCPU):", podMetrics.CpuRequest.Value())
	fmt.Printf("      %-25s%.2f\n", "Memory Request(GB):", bytesToGB(podMetrics.MemoryRequest.Value()))

	fmt.Printf("Storage:\n")

	pvs := GetClusterVolumes()
	storageCost, storageCapacity := getPvCostAndCapacity(pvs)

	fmt.Printf("   %-25s   %d\n", "Persistent Volume count:", len(pvs))
	fmt.Printf("   %-25s   %.2f\n", "Capacity(GB):", bytesToGB(storageCapacity))
	fmt.Printf("   %-25s   %.2f$\n", "Cost:", storageCost)

	pvcs := GetClusterPersistentVolumeClaims()
	_, pvcCapacity := getPvcCostAndCapacity(pvcs)
	fmt.Printf("   %-25s   %d\n", "PV Claim count:", len(pvcs))

	//fmt.Printf("   %-25s   %.2f$\n", "PVC Cost:", pvcCost)
	fmt.Printf("   %-25s   %.2f\n", "PV Claim Capacity(GB):", bytesToGB(pvcCapacity))


	fmt.Printf("Cost:\n")
	fmt.Printf("   %-25s   %.2f$\n", "Compute cost:", computeCost)
	fmt.Printf("   %-25s   %.2f$\n", "Storage cost:", storageCost)
	fmt.Printf("   %-25s   %.2f$\n", "Total cost:", computeCost + storageCost)
}

func bytesToGB(val int64) float64 {
	return float64(val) / (1024.0 * 1024.0 * 1024.0)
}

func getPvCostAndCapacity(pvs []v1.PersistentVolume) (float64, int64) {
	var storageCost = 0.0
	var storageCapacity =  resource.Quantity{}
	for _, pv := range pvs {
		storageClass := pv.Spec.StorageClassName
		total := getMonthToDateCostForStorageClass(storageClass)
		var cur = resource.Quantity{}
		cur.Add(pv.Spec.Capacity["storage"])
		storageCost += total * bytesToGB(cur.Value())
		storageCapacity.Add(pv.Spec.Capacity["storage"])
	}
	return storageCost, storageCapacity.Value()
}

func getPvcCostAndCapacity(pvcs []v1.PersistentVolumeClaim) (float64, int64) {
	var pvcCapacity = resource.Quantity{}
	var pvcCost = 0.0
	for _, pvc := range pvcs {
		pvcCapacity.Add(pvc.Spec.Resources.Requests["storage"])

		storageClass := pvc.Spec.StorageClassName
		total := getMonthToDateCostForStorageClass(*storageClass)
		var cur = resource.Quantity{}
		cur.Add(pvc.Spec.Resources.Requests["storage"])
		pvcCost += total * bytesToGB(cur.Value())
	}
	return pvcCost, pvcCapacity.Value()
}

func getSavings() {
	fmt.Printf("Savings Summary\n")

	fmt.Printf("Storage:\n")
	pvs := GetClusterVolumes()
	storageCost, storageCapacity := getPvCostAndCapacity(pvs)

	pvcs := GetClusterPersistentVolumeClaims()
	pvcCost, pvcCapacity := getPvcCostAndCapacity(pvcs)
	mtdSaving := storageCost - pvcCost
	projectedSaving := projectToMonth(mtdSaving)

	fmt.Printf("   %-25s   %d\n", "Unused Volumes:", len(pvs) - len(pvcs))
	fmt.Printf("   %-25s   %.2f\n", "Unused Capacity(GB):", bytesToGB(storageCapacity - pvcCapacity))
	fmt.Printf("   %-25s   %.2f$\n", "Month To Date Savings:", mtdSaving)
	fmt.Printf("   %-25s   %.2f$\n", "Projected Monthly Savings:", projectedSaving)

	/*fmt.Printf("Compute:\n")
	pods := GetClusterPods()
	var ps = []*Pod{}

	for _, p := range pods {
		var temp Pod
		temp.name = p.GetObjectMeta().GetName()
		temp.nodeName = p.Spec.NodeName
		j := 0
		podVolumes := []*string{}
		for j < len(p.Spec.Volumes) {
			vol := p.Spec.Volumes[j]
			if vol.PersistentVolumeClaim != nil {
				podVolumes = append(podVolumes, &vol.PersistentVolumeClaim.ClaimName)
			}
			j++
		}
		temp.pvcs = podVolumes
		ps = append(ps, &temp)
	}

	ps = getPodsComputeCost(ps)
	totalPodCost := 0.0
	for _, temp := range ps {
		totalPodCost += temp.cost.cpuCost + temp.cost.memoryCost
	}

	/*podMetrics := metrics.CalculatePodStatsFromContainers(pods)
	//fmt.Println(podMetrics)

	var computeCost = 0.0
	nodes := GetClusterNodes()
	for _, node := range nodes {
		instanceType := GetNodeType(node)
		total, _, _ := getMonthToDateCostForInstanceType(instanceType)
		computeCost = computeCost + total
	}
	nodeMetrics := metrics.CalculateNodeStats(nodes)
	fmt.Printf("   %-25s   %d\n", "Unused CPU(vCPU):", nodeMetrics.CpuLimit.Value() - podMetrics.CpuRequest.Value())
	fmt.Printf("   %-25s   %.2f\n", "Unused Memory(GB):", bytesToGB(nodeMetrics.MemoryLimit.Value()) -
		bytesToGB(podMetrics.MemoryRequest.Value()))
	mtdComputeSaving := computeCost - totalPodCost
	projectedComputeSaving := projectToMonth(mtdSaving)
	fmt.Printf("   %-25s   %.2f$\n", "Month To Date Savings:", mtdComputeSaving)
	fmt.Printf("   %-25s   %.2f$\n", "Projected Monthly Savings:", projectedComputeSaving)*/
}

func getPodCost(podName string) {
	pod := getPodDetailsFromClient(podName)
	pods := getPodsCost([]*Pod{pod})
	printPodsVerbose(pods)
}

func getPodsCost(pods []*Pod) []*Pod {
	nodes := map[string]*Node{}
	pvcs := map[string]*PersistentVolumeClaim{}
	for _, pod := range pods {
		nodes[pod.nodeName] = nil
		for _, pvc := range pod.pvcs {
			pvcs[*pvc] = nil
		}
	}
	nodes = collectNodes(nodes)
	pvcs = collectPersistentVolumeClaims(pvcs)
	pods = calculateCost(pods, nodes, pvcs)
	//printPodsVerbose(pods)
	return pods
}

func getPodsComputeCost(pods []*Pod) []*Pod  {
	nodes := map[string]*Node{}
	nodes = collectNodes(nodes)

	i := 0
	for i <= len(pods)-1 {
		node := nodes[pods[i].nodeName]
		pods[i].nodeCostPercentage = (float64)(node.getPodResourcePercentage(pods[i].name))
		totalComputeCost, cpuCost, memoryCost := getMonthToDateCostForInstanceType(node.instanceType)

		podCost := Cost{}
		podCost.totalCost = pods[i].nodeCostPercentage*totalComputeCost
		podCost.cpuCost = pods[i].nodeCostPercentage * cpuCost
		podCost.memoryCost = pods[i].nodeCostPercentage * memoryCost
		pods[i].cost = &podCost
		i++
	}
	return pods
}

func getAllNodesCost() {
	nodes := getAllNodeDetailsFromClient()
	i := 0
	for i < len(nodes) {
		totalComputeCost, cpuCost, memoryCost := getMonthToDateCostForInstanceType(nodes[i].instanceType)
		nodeCost := Cost{}
		nodeCost.totalCost = totalComputeCost
		nodeCost.cpuCost = cpuCost
		nodeCost.memoryCost = memoryCost
		nodes[i].cost = &nodeCost
		i++
	}
	printNodeDetails(nodes)
	//fmt.Println(nodes)
}