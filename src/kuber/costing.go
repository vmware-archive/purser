package main

func calculateCost(pods []Pod, nodes map[string]*Node) []Pod {
	i := 0
	for i <= len(pods)-1 {
		node := nodes[pods[i].nodeName]
		pods[i].nodeCostPercentage = (float64)(node.getPodResourcePercentage(pods[i].name))
		totalCost, cpuCost, memoryCost := getMonthToDateCostForInstanceType(node.instanceType)
		podCost := Cost{}
		podCost.totalCost = pods[i].nodeCostPercentage * totalCost
		podCost.cpuCost = pods[i].nodeCostPercentage * cpuCost
		podCost.memoryCost = pods[i].nodeCostPercentage * memoryCost
		pods[i].cost = podCost
		i++
	}
	return pods
}

func getPodsCostForLabel(label string) {
	pods := getPodsForLabel(label)
	nodes := map[string]*Node{}
	for _, val := range pods {
		nodes[val.nodeName] = nil
	}
	nodes = collectNodes(nodes)
	pods = calculateCost(pods, nodes)
	printPodsVerbose(pods)
}
