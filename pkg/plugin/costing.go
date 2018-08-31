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

package plugin

import (
	"fmt"

	"github.com/vmware/purser/pkg/plugin/crd"
	"github.com/vmware/purser/pkg/plugin/metrics"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
)

// ClientSetInstance helps in accessing kubernetes apis through client.
var ClientSetInstance *kubernetes.Clientset

// ProvideClientSetInstance sets the client set instance.
func ProvideClientSetInstance(clientset *kubernetes.Clientset) {
	ClientSetInstance = clientset
}

// GetPodsCostForLabel returns pods cost for given label.
func GetPodsCostForLabel(label string) {
	pods := getPodsForLabelThroughClient(label)
	pods = getPodsCost(pods)
	printPodsVerbose(pods)
}

// GetClusterSummary summarizes cluster metrics.
func GetClusterSummary() {
	pods := GetClusterPods()
	podMetrics := metrics.CalculatePodStatsFromContainers(pods)
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
	fmt.Printf("      %-25s%d\n", "CPU(vCPU):", nodeMetrics.CPULimit.Value())

	fmt.Printf("      %-25s%.2f\n", "Memory(GB):", bytesToGB(nodeMetrics.MemoryLimit.Value()))

	fmt.Printf("   Provisioned Resources:\n")
	fmt.Printf("      %-25s%d\n", "CPU Request(vCPU):", podMetrics.CPURequest.Value())
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

	fmt.Printf("   %-25s   %.2f\n", "PV Claim Capacity(GB):", bytesToGB(pvcCapacity))

	fmt.Printf("Cost:\n")
	fmt.Printf("   %-25s   %.2f$\n", "Compute cost:", computeCost)
	fmt.Printf("   %-25s   %.2f$\n", "Storage cost:", storageCost)
	fmt.Printf("   %-25s   %.2f$\n", "Total cost:", computeCost+storageCost)
}

// GetSavings returns the savings summary.
func GetSavings() {
	fmt.Printf("Savings Summary\n")

	fmt.Printf("Storage:\n")
	pvs := GetClusterVolumes()
	storageCost, storageCapacity := getPvCostAndCapacity(pvs)

	pvcs := GetClusterPersistentVolumeClaims()
	pvcCost, pvcCapacity := getPvcCostAndCapacity(pvcs)
	mtdSaving := storageCost - pvcCost
	projectedSaving := projectToMonth(mtdSaving)

	fmt.Printf("   %-25s   %d\n", "Unused Volumes:", len(pvs)-len(pvcs))
	fmt.Printf("   %-25s   %.2f\n", "Unused Capacity(GB):", bytesToGB(storageCapacity-pvcCapacity))
	fmt.Printf("   %-25s   %.2f$\n", "Month To Date Savings:", mtdSaving)
	fmt.Printf("   %-25s   %.2f$\n", "Projected Monthly Savings:", projectedSaving)
}

// GetPodCost returns the cumulative cost for the pods.
func GetPodCost(podName string) {
	pod := getPodDetailsFromClient(podName)
	pods := getPodsCost([]*Pod{pod})
	printPodsVerbose(pods)
}

// GetAllNodesCost returns the cumulative cost of all the nodes.
func GetAllNodesCost() {
	nodes := getAllNodeDetailsFromClient()

	for i := 0; i < len(nodes); i++ {
		totalComputeCost, cpuCost, memoryCost := getMonthToDateCostForInstanceType(nodes[i].instanceType)
		nodeCost := Cost{
			totalCost:  totalComputeCost,
			cpuCost:    cpuCost,
			memoryCost: memoryCost,
		}
		nodes[i].cost = &nodeCost
	}
	printNodeDetails(nodes)
}

func calculateCost(pods []*Pod, nodes map[string]*Node, pvcs map[string]*PersistentVolumeClaim) []*Pod {
	for i := 0; i <= len(pods)-1; i++ {
		node := nodes[pods[i].nodeName]
		pods[i].nodeCostPercentage = node.getPodResourcePercentage(pods[i].name)
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
		podCost := Cost{
			totalCost:   pods[i].nodeCostPercentage*totalComputeCost + totalStorageCost,
			cpuCost:     pods[i].nodeCostPercentage * cpuCost,
			memoryCost:  pods[i].nodeCostPercentage * memoryCost,
			storageCost: totalStorageCost,
		}
		pods[i].cost = &podCost
	}
	return pods
}

func bytesToGB(val int64) float64 {
	return float64(val) / (1024.0 * 1024.0 * 1024.0)
}

func getPvCostAndCapacity(pvs []v1.PersistentVolume) (float64, int64) {
	var storageCost = 0.0
	var storageCapacity = resource.Quantity{}
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
	return pods
}

// GetGroupCost returns Cost (total, cpu, memory and storage) for a Group
func GetGroupCost(group *crd.Group) *Cost {
	cpuCostPerCPUPerHour, memCostPerGBPerHour, storageCostPerGBPerHour := GetUserCosts()

	currentTime := getCurrentTime()
	monthStartTime := getCurrentMonthStartTime()

	podsDetails := group.Spec.PodsDetails
	var totalCPUCost, totalMemoryCost, totalStorageCost, totalCumulativeCost float64
	for podName, podDetails := range podsDetails {
		startTime := podDetails.StartTime
		endTime := podDetails.EndTime

		podActiveHours := currentMonthActiveTimeInHours(startTime, endTime, currentTime, monthStartTime)

		podMetrics := group.Spec.PodsMetrics[podName]
		podCPURequest := resourceQuantityToFloat64(podMetrics.CPURequest)
		podMemRequest := resourceQuantityToFloat64(podMetrics.MemoryRequest)

		// TODO: find podStorage
		podStorageClaimed := 0.0

		totalCPUCost += cpuCostPerCPUPerHour * podCPURequest * podActiveHours
		totalMemoryCost += memCostPerGBPerHour * podMemRequest * podActiveHours
		totalStorageCost += storageCostPerGBPerHour * podStorageClaimed * podActiveHours
	}

	totalCumulativeCost = totalCPUCost + totalMemoryCost + totalStorageCost

	return &Cost{
		totalCost:   totalCumulativeCost,
		cpuCost:     totalCPUCost,
		memoryCost:  totalMemoryCost,
		storageCost: totalStorageCost,
	}
}
