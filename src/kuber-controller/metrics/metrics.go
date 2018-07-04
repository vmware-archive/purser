package metrics

import (
	"fmt"
	api_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Metrics struct {
	CpuLimit      *resource.Quantity
	MemoryLimit   *resource.Quantity
	CpuRequest    *resource.Quantity
	MemoryRequest *resource.Quantity
}

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
