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

	log "github.com/Sirupsen/logrus"

	groups_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	groups "github.com/vmware/purser/pkg/client/clientset/typed/groups/v1"
	"github.com/vmware/purser/pkg/plugin/metrics"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetGroupByName return group CRD by name.
func GetGroupByName(groupClient *groups.GroupClient, groupName string) *groups_v1.Group {
	group, err := groupClient.Get(groupName)
	if err != nil {
		log.Errorf("failed to get custom group by name %s, %v", groupName, err)
	}
	return group
}

// GetGroupDetails returns aggregated metrics (cpu, memory, storage) and cost (total, cpu, memory and storage) of a Group
func GetGroupDetails(group *groups_v1.Group) (*metrics.GroupMetrics, *metrics.GroupMetrics, *Cost) {
	// TODO: include storage in group details
	price := GetUserCosts()

	currentTime := getCurrentTime()
	monthStartTime := getCurrentMonthStartTime()

	podsDetails := group.Spec.PodsDetails
	var totalCPURequest, totalCPULimit, totalMemoryRequest, totalMemoryLimit, totalStorageClaimed float64
	// [PIT] Point In Time metrics for the group
	var pitCPURequest, pitCPULimit, pitMemoryRequest, pitMemoryLimit, pitStorageClaimed float64
	for _, podDetails := range podsDetails {
		startTime := podDetails.StartTime
		endTime := podDetails.EndTime

		podActiveHours := currentMonthActiveTimeInHoursMulti(startTime, endTime, currentTime, monthStartTime)

		podComputeMetrics := getPodComputeMetrics(podDetails)
		podStorageAllocatedInGBHours, podActiveStorageAllocated := getPodStorageMetrics(podDetails)

		podCPURequest := float64(podComputeMetrics.CPURequest.Value())
		podMemRequest := bytesToGB(podComputeMetrics.MemoryRequest.Value())
		podCPULimit := float64(podComputeMetrics.CPULimit.Value())
		podMemLimit := bytesToGB(podComputeMetrics.MemoryLimit.Value())

		// Pod is alive
		if endTime.IsZero() {
			pitCPULimit += podCPULimit
			pitCPURequest += podCPURequest
			pitMemoryLimit += podMemLimit
			pitMemoryRequest += podMemRequest
			pitStorageClaimed += podActiveStorageAllocated
		}

		totalCPURequest += podCPURequest * podActiveHours
		totalMemoryRequest += podMemRequest * podActiveHours
		totalCPULimit += podCPULimit * podActiveHours
		totalMemoryLimit += podMemLimit * podActiveHours

		totalStorageClaimed += podStorageAllocatedInGBHours
	}

	totalCPUCost := price.CPU * totalCPURequest
	totalMemoryCost := price.Memory * totalMemoryRequest
	totalStorageCost := price.Storage * totalStorageClaimed

	totalCumulativeCost := totalCPUCost + totalMemoryCost + totalStorageCost

	mtdGroupMetrics := &metrics.GroupMetrics{
		CPULimit:       totalCPULimit,
		MemoryLimit:    totalMemoryLimit,
		CPURequest:     totalCPURequest,
		MemoryRequest:  totalMemoryRequest,
		StorageClaimed: totalStorageClaimed,
	}

	pitGroupMetrics := &metrics.GroupMetrics{
		CPULimit:       pitCPULimit,
		MemoryLimit:    pitMemoryLimit,
		CPURequest:     pitCPURequest,
		MemoryRequest:  pitMemoryRequest,
		StorageClaimed: pitStorageClaimed,
	}

	cost := &Cost{
		TotalCost:   totalCumulativeCost,
		CPUCost:     totalCPUCost,
		MemoryCost:  totalMemoryCost,
		StorageCost: totalStorageCost,
	}

	return pitGroupMetrics, mtdGroupMetrics, cost
}

// getPodComputeMetrics returns PodMetrics given its details(container metrics)
func getPodComputeMetrics(podDetails *groups_v1.PodDetails) *metrics.Metrics {
	CPURequest := resource.Quantity{}
	MemoryRequest := resource.Quantity{}
	CPULimit := resource.Quantity{}
	MemoryLimit := resource.Quantity{}
	for _, container := range podDetails.Containers {
		addResourceAToResourceB(container.Metrics.CPURequest, &CPURequest)
		addResourceAToResourceB(container.Metrics.MemoryRequest, &MemoryRequest)
		addResourceAToResourceB(container.Metrics.MemoryLimit, &MemoryLimit)
		addResourceAToResourceB(container.Metrics.CPULimit, &CPULimit)
	}
	return &metrics.Metrics{
		CPULimit:      &CPULimit,
		MemoryLimit:   &MemoryLimit,
		CPURequest:    &CPURequest,
		MemoryRequest: &MemoryRequest,
	}
}

func getPodStorageMetrics(podDetails *groups_v1.PodDetails) (float64, float64) {
	var podStorageAllocatedInGBHours, podActiveStorageAllocated float64
	currentTime := getCurrentTime()
	monthStart := getCurrentMonthStartTime()

	podVolumeClaims := podDetails.PodVolumeClaims
	for _, pvc := range podVolumeClaims {
		for i := 0; i < len(pvc.CapacityAllocatedInGB); i++ {
			allocation := pvc.CapacityAllocatedInGB[i]
			boundTime := pvc.BoundTimes[i]
			unboundTime := getUnboundTime(pvc, i)
			activeTime := currentMonthActiveTimeInHoursMulti(boundTime, unboundTime, currentTime, monthStart)
			podStorageAllocatedInGBHours += allocation * activeTime

			// pvc is active(i.e, unboundTime is nil if pvc is still bounded)
			if unboundTime.IsZero() {
				podActiveStorageAllocated += allocation
			}
		}
	}
	return podStorageAllocatedInGBHours, podActiveStorageAllocated
}

func getUnboundTime(pvc *groups_v1.PersistentVolumeClaim, i int) metav1.Time {
	var unboundTime metav1.Time
	if i < len(pvc.UnboundTimes) {
		return pvc.UnboundTimes[i]
	}
	// unboundTime is nil if pvc is still bounded (i.e, length of BoundTimes - length of UnboundTime is 1.)
	return unboundTime
}

// PrintGroup displays the group information.
func PrintGroup(group *groups_v1.Group) {
	pitGroupMetrics, mtdGroupMetrics, cost := GetGroupDetails(group)

	fmt.Printf("%-30s             %s\n", "Group Name:", group.Name)
	fmt.Println()
	fmt.Println("Point in Time Resource Stats:")
	fmt.Printf("             %-30s%.2f\n", "CPU Limit(vCPU):", pitGroupMetrics.CPULimit)
	fmt.Printf("             %-30s%.2f\n", "Memory Limit(GB):", pitGroupMetrics.MemoryLimit)
	fmt.Printf("             %-30s%.2f\n", "CPU Request(vCPU):", pitGroupMetrics.CPURequest)
	fmt.Printf("             %-30s%.2f\n", "Memory Request(GB):", pitGroupMetrics.MemoryRequest)
	fmt.Printf("             %-30s%.2f\n", "Storage Claimed(GB):", pitGroupMetrics.StorageClaimed)

	fmt.Println()
	fmt.Printf("%-30s\n", "Month to Date Active Resource Stats:")
	fmt.Printf("             %-30s%.2f\n", "CPU Request(vCPU-hours):", mtdGroupMetrics.CPURequest)
	fmt.Printf("             %-30s%.2f\n", "Memory Request(GB-hours):", mtdGroupMetrics.MemoryRequest)
	fmt.Printf("             %-30s%.2f\n", "Storage Claimed(GB-hours):", mtdGroupMetrics.StorageClaimed)

	fmt.Println()
	fmt.Printf("%-30s\n", "Month to Date Cost Stats:")
	fmt.Printf("             %-30s%.2f\n", "CPU Cost($):", cost.CPUCost)
	fmt.Printf("             %-30s%.2f\n", "Memory Cost($):", cost.MemoryCost)
	fmt.Printf("             %-30s%.2f\n", "Storage Cost($):", cost.StorageCost)
	fmt.Printf("             %-30s%.2f\n", "Total Cost($):", cost.TotalCost)
}
