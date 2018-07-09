package main

import (
	"fmt"
	"kuber/metrics"
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

	fmt.Printf("===Cluster Details===\n")
	fmt.Printf("Provision:\n")
	fmt.Printf("\tCpu Limit = %s\n", podMetrics.CpuLimit.String())
	fmt.Printf("\tMemory Limit = %s\n", podMetrics.MemoryLimit.String())
	fmt.Printf("\tCpu Request = %s\n", podMetrics.CpuRequest.String())
	fmt.Printf("\tMemory Request = %s\n", podMetrics.MemoryRequest.String())

	nodes := GetClusterNodes()
	nodeMetrics := metrics.CalculateNodeStats(nodes)
	fmt.Printf("Capacity:\n")
	fmt.Printf("\tCpu Limit = %s\n", nodeMetrics.CpuLimit.String())
	fmt.Printf("\tMemory Limit = %s\n", nodeMetrics.MemoryLimit.String())

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
