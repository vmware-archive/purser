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
	"fmt"
	api_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Metrics struct {
	CpuLimit      *resource.Quantity
	MemoryLimit   *resource.Quantity
	CpuRequest    *resource.Quantity
	MemoryRequest *resource.Quantity
}

func CalculatePodStatsFromContainers(pods []v1.Pod) *Metrics {
	cpuLimit := &resource.Quantity{}
	memoryLimit := &resource.Quantity{}
	cpuRequest := &resource.Quantity{}
	memoryRequest := &resource.Quantity{}
	for _, pod := range pods {
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
	}
	return &Metrics{
		CpuLimit:      cpuLimit,
		MemoryLimit:   memoryLimit,
		CpuRequest:    cpuRequest,
		MemoryRequest: memoryRequest,
	}
}

func CalculateNodeStats(nodes []v1.Node) *Metrics {
	cpuLimit := &resource.Quantity{}
	memoryLimit := &resource.Quantity{}
	cpuRequest := &resource.Quantity{}
	memoryRequest := &resource.Quantity{}
	for _, node := range nodes {
		cpuLimit.Add(*node.Status.Capacity.Cpu())
		memoryLimit.Add(*node.Status.Capacity.Memory())
	}
	return &Metrics{
		CpuLimit:      cpuLimit,
		MemoryLimit:   memoryLimit,
		CpuRequest:    cpuRequest,
		MemoryRequest: memoryRequest,
	}
}

func PrintPodStats(pod *api_v1.Pod, metrics *Metrics) {
	fmt.Printf("Pod:\t%s\n", pod.Name)
	fmt.Printf("\tCpu Limit = %s\n", metrics.CpuLimit.String())
	fmt.Printf("\tMemory Limit = %s\n", metrics.MemoryLimit.String())
	fmt.Printf("\tCpu Request = %s\n", metrics.CpuRequest.String())
	fmt.Printf("\tMemory Request = %s\n", metrics.MemoryRequest.String())
}
