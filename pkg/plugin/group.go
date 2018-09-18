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
	"github.com/vmware/purser/pkg/plugin/crd"
	"github.com/vmware/purser/pkg/plugin/metrics"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetGroupDetails returns aggregated metrics (cpu, memory, storage) and cost (total, cpu, memory and storage) of a Group
func GetGroupDetails(group *crd.Group) (*metrics.GroupMetrics, *metrics.GroupMetrics, *Cost) {
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
func getPodComputeMetrics(podDetails *crd.PodDetails) *metrics.Metrics {
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

func getPodStorageMetrics(podDetails *crd.PodDetails) (float64, float64) {
	var podStorageAllocatedInGBHours, podActiveStorageAllocated float64
	currentTime := getCurrentTime()
	monthStart := getCurrentMonthStartTime()

	podVolumeClaims := podDetails.PodVolumeClaims
	for _, pvc := range podVolumeClaims {
		for i := 0; i < len(pvc.CapacityAllotedInGB); i++ {
			allocation := pvc.CapacityAllotedInGB[i]
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

func getUnboundTime(pvc *crd.PersistentVolumeClaim, i int) metav1.Time {
	var unboundTime metav1.Time
	if i < len(pvc.UnboundTimes) {
		return pvc.UnboundTimes[i]
	}
	// unboundTime is nil if pvc is still bounded (i.e, length of BoundTimes - length of UnboundTime is 1.)
	return unboundTime
}
