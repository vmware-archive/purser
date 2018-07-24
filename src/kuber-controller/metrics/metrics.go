package metrics

import (
	api_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	log "github.com/Sirupsen/logrus"
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
	log.Printf("Pod:\t%s\n", pod.Name)
	log.Printf("\tCpu Limit = %s\n", metrics.CpuLimit.String())
	log.Printf("\tMemory Limit = %s\n", metrics.MemoryLimit.String())
	log.Printf("\tCpu Request = %s\n", metrics.CpuRequest.String())
	log.Printf("\tMemory Request = %s\n", metrics.MemoryRequest.String())
}
