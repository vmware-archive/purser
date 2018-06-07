package main

const (
	getPodsByLabel      = "kubectl --kubeconfig=%s get pods -l %s -o json"
	getPodByName        = "kubectl --kubeconfig=%s get pod %s -o json"
	getNodeByName       = "kubectl --kubeconfig=%s get node %s -o json"
	nodeDescribeCommand = "kubectl --kubeconfig=%s describe node %s"
)
