package main

import (
	"fmt"
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
		pods[i].cost = podCost
		i++
	}
	return pods
}

func getPodsCostForLabel(label string) {
	//pods := getPodsForLabel(label)
	pods := getPodsForLabelThroughClient(label)
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
	printPodsVerbose(pods)
}
