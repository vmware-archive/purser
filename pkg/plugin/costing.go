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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Cost details
type Cost struct {
	totalCost   float64
	cpuCost     float64
	memoryCost  float64
	storageCost float64
}

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

	nodeMetrics := metrics.CalculateNodeStats(nodes)
	fmt.Printf("   Total Capacity:\n")
	fmt.Printf("      %-25s%d\n", "CPU(vCPU):", nodeMetrics.CPULimit.Value())

	fmt.Printf("      %-25s%.2f\n", "Memory(GB):", bytesToGB(nodeMetrics.MemoryLimit.Value()))

	fmt.Printf("   Provisioned Resources:\n")
	fmt.Printf("      %-25s%d\n", "CPU Request(vCPU):", podMetrics.CPURequest.Value())
	fmt.Printf("      %-25s%.2f\n", "Memory Request(GB):", bytesToGB(podMetrics.MemoryRequest.Value()))

	price := GetUserCosts()
	hoursInMonthTillNow := totalHoursTillNow()

	cpuCost := float64(nodeMetrics.CPULimit.Value()) * hoursInMonthTillNow * price.CPU
	memCost := bytesToGB(nodeMetrics.MemoryLimit.Value()) * hoursInMonthTillNow * price.Memory
	computeCost := cpuCost + memCost
	fmt.Printf("   %-25s   %.2f$\n", "Cost:", computeCost)

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
	nodes := GetClusterNodes()

	price := GetUserCosts()
	hoursInMonthTillNow := totalHoursTillNow()

	fmt.Println("Node name\tNode cpu-cost\tNode mem-cost\tNode total-cost")
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		nodeMetrics := metrics.CalculateNodeStats([]v1.Node{node})

		cpuCost := float64(nodeMetrics.CPULimit.Value()) * hoursInMonthTillNow * price.CPU
		memoryCost := bytesToGB(nodeMetrics.MemoryLimit.Value()) * hoursInMonthTillNow * price.Memory
		totalComputeCost := cpuCost + memoryCost

		fmt.Printf("%s\t%f\t%f\t%f\n", node.Name, cpuCost, memoryCost, totalComputeCost)
	}
}

func calculateCost(pods []*Pod, pvcs map[string]*PersistentVolumeClaim) []*Pod {
	price := GetUserCosts()
	for i := 0; i <= len(pods)-1; i++ {
		pod := pods[i]
		pods[i].cost = calculateCostOfPod(*pod, pvcs, price)
	}
	return pods
}

func calculateCostOfPod(pod Pod, pvcs map[string]*PersistentVolumeClaim, price *Price) *Cost {
	podDurationInHours := currentMonthActiveTimeInHours(pod.startTime, metav1.Now())

	podCPUCost := float64(pod.podMetrics.CPURequest.Value()) * podDurationInHours * (price.CPU)
	podMemoryCost := bytesToGB(pod.podMetrics.MemoryRequest.Value()) * podDurationInHours * (price.Memory)

	podStorageCost := 0.0
	for _, pvc := range pod.pvcs {
		if pvcs[*pvc] != nil {
			podStorageCost += pvcs[*pvc].capacityAllotedInGB * podDurationInHours * (price.Storage)
		} else {
			fmt.Printf("Persistent volume claim is not present for %s\n", *pvc)
		}
	}
	podTotalCost := podCPUCost + podMemoryCost + podStorageCost
	return &Cost{
		totalCost:   podTotalCost,
		cpuCost:     podCPUCost,
		memoryCost:  podMemoryCost,
		storageCost: podStorageCost,
	}
}

func getPvCostAndCapacity(pvs []v1.PersistentVolume) (float64, int64) {
	price := GetUserCosts()
	hoursInMonthTillNow := totalHoursTillNow()
	var storageCapacity = resource.Quantity{}
	for _, pv := range pvs {
		storageCapacity.Add(pv.Spec.Capacity["storage"])
	}
	storageCost := bytesToGB(storageCapacity.Value()) * hoursInMonthTillNow * price.Storage
	return storageCost, storageCapacity.Value()
}

func getPvcCostAndCapacity(pvcs []v1.PersistentVolumeClaim) (float64, int64) {
	price := GetUserCosts()
	hoursInMonthTillNow := totalHoursTillNow()
	var pvcCapacity = resource.Quantity{}
	for _, pvc := range pvcs {
		pvcCapacity.Add(pvc.Spec.Resources.Requests["storage"])
	}
	pvcCost := bytesToGB(pvcCapacity.Value()) * hoursInMonthTillNow * price.Storage
	return pvcCost, pvcCapacity.Value()
}

func getPodsCost(pods []*Pod) []*Pod {
	pvcs := map[string]*PersistentVolumeClaim{}
	for _, pod := range pods {
		for _, pvc := range pod.pvcs {
			pvcs[*pvc] = nil
		}
	}
	pvcs = collectPersistentVolumeClaims(pvcs)
	pods = calculateCost(pods, pvcs)
	return pods
}

// GetGroupDetails returns aggregated metrics (cpu, memory, storage) and cost (total, cpu, memory and storage) of a Group
func GetGroupDetails(group *crd.Group) (*metrics.GroupMetrics, *Cost) {
	// TODO: include storage in group details
	price := GetUserCosts()

	currentTime := getCurrentTime()
	monthStartTime := getCurrentMonthStartTime()

	podsDetails := group.Spec.PodsDetails
	var totalCPURequest, totalCPULimit, totalMemoryRequest, totalMemoryLimit, totalStorageClaimed float64
	for podName, podDetails := range podsDetails {
		startTime := podDetails.StartTime
		endTime := podDetails.EndTime

		podActiveHours := currentMonthActiveTimeInHoursMulti(startTime, endTime, currentTime, monthStartTime)

		podMetrics := group.Spec.PodsMetrics[podName]

		podCPURequest := float64(podMetrics.CPURequest.Value())
		podMemRequest := bytesToGB(podMetrics.MemoryRequest.Value())
		podCPULimit := float64(podMetrics.CPULimit.Value())
		podMemLimit := bytesToGB(podMetrics.MemoryLimit.Value())

		totalCPURequest += podCPURequest * podActiveHours
		totalMemoryRequest += podMemRequest * podActiveHours
		totalCPULimit += podCPULimit * podActiveHours
		totalMemoryLimit += podMemLimit * podActiveHours

		// TODO: find podStorage
		podStorageClaimed := 0.0
		totalStorageClaimed += podStorageClaimed * podActiveHours
	}

	totalCPUCost := price.CPU * totalCPURequest
	totalMemoryCost := price.Memory * totalMemoryRequest
	totalStorageCost := price.Storage * totalStorageClaimed

	totalCumulativeCost := totalCPUCost + totalMemoryCost + totalStorageCost

	groupMetrics := &metrics.GroupMetrics{
		ActiveCPULimit:       totalCPULimit,
		ActiveMemoryLimit:    totalMemoryLimit,
		ActiveCPURequest:     totalCPURequest,
		ActiveMemoryRequest:  totalMemoryRequest,
		ActiveStorageClaimed: totalStorageClaimed,
	}

	cost := &Cost{
		totalCost:   totalCumulativeCost,
		cpuCost:     totalCPUCost,
		memoryCost:  totalMemoryCost,
		storageCost: totalStorageCost,
	}

	return groupMetrics, cost
}
