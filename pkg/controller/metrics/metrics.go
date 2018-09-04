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

package metrics

import (
	log "github.com/Sirupsen/logrus"
	api_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Metrics types
type Metrics struct {
	CPULimit      *resource.Quantity
	MemoryLimit   *resource.Quantity
	CPURequest    *resource.Quantity
	MemoryRequest *resource.Quantity
}

// CalculatePodStatsFromContainers returns the cumulative metrics from the containers.
func CalculatePodStatsFromContainers(pod *api_v1.Pod) *Metrics {
	cpuLimit := &resource.Quantity{}
	memoryLimit := &resource.Quantity{}
	cpuRequest := &resource.Quantity{}
	memoryRequest := &resource.Quantity{}
	for _, c := range pod.Spec.Containers {
		limits := c.Resources.Limits
		if limits != nil {
			cpuLimit.Add(*limits.Cpu())
			memoryLimit.Add(*limits.Memory())
		}

		requests := c.Resources.Requests
		if requests != nil {
			cpuRequest.Add(*requests.Cpu())
			memoryRequest.Add(*requests.Memory())
		}
	}
	return &Metrics{
		CPULimit:      cpuLimit,
		MemoryLimit:   memoryLimit,
		CPURequest:    cpuRequest,
		MemoryRequest: memoryRequest,
	}
}

// PrintPodStats displays the pod stats.
func PrintPodStats(pod *api_v1.Pod, metrics *Metrics) {
	log.Printf("Pod:\t%s\n", pod.Name)
	log.Printf("\tCPU Limit = %s\n", metrics.CPULimit.String())
	log.Printf("\tMemory Limit = %s\n", metrics.MemoryLimit.String())
	log.Printf("\tCPU Request = %s\n", metrics.CPURequest.String())
	log.Printf("\tMemory Request = %s\n", metrics.MemoryRequest.String())
}
