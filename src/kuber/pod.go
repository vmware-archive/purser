package main

import (
	"fmt"
	"os"

	"github.com/tidwall/gjson"
)

// Cost details
type Cost struct {
	totalCost  float64
	cpuCost    float64
	memoryCost float64
}

// Pod Information
type Pod struct {
	name               string
	nodeName           string
	nodeCostPercentage float64
	cost               Cost
	pvcs               []*string
}

func getPodsForLabel(label string) []Pod {
	pods := []Pod{}
	command := fmt.Sprintf(getPodsByLabel, os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), label)
	bytes := executeCommand(command)
	json := string(bytes)
	items := gjson.Get(json, "items")

	items.ForEach(func(key, value gjson.Result) bool {
		name := value.Get("metadata.name")
		nodeName := value.Get("spec.nodeName")
		pod := Pod{name: name.Str, nodeName: nodeName.Str}

		podVolumes := []*string{}
		volumes := value.Get("spec.volumes")
		volumes.ForEach(func(volKey, volume gjson.Result) bool {
			pvc := volume.Get("persistentVolumeClaim.claimName")
			if pvc.Exists() {
				podVolumes = append(podVolumes, &pvc.Str)
			}
			return true
		})
		pod.pvcs = podVolumes
		pods = append(pods, pod)
		return true
	})
	return pods
}

func printPodsVerbose(pods []Pod) {
	i := 0
	fmt.Printf("==Pods Cost Details==\n")
	for i <= len(pods)-1 {
		fmt.Printf("%-25s%s\n", "Pod Name:", pods[i].name)
		fmt.Printf("%-25s%s\n", "Node:", pods[i].nodeName)
		fmt.Printf("%-25s%.2f\n", "Pod Cost Percentage:", pods[i].nodeCostPercentage*100.0)
		fmt.Printf("%-25s\n", "Persistent Volume Claims:")

		j := 0
		for j <= len(pods[i].pvcs)-1 {
			fmt.Printf("    %s\n", *pods[i].pvcs[j])
			j++
		}
		fmt.Printf("%-25s\n", "Cost:")
		fmt.Printf("    %-21s%f$\n", "Total Cost:", pods[i].cost.totalCost)
		fmt.Printf("    %-21s%f$\n", "CPU Cost:", pods[i].cost.cpuCost)
		fmt.Printf("    %-21s%f$\n", "Memory Cost:", pods[i].cost.memoryCost)
		fmt.Printf("\n")
		i++
	}
}

func printPodDetails(pods []Pod) {
	fmt.Println("===POD Details===")
	fmt.Println("POD Name \t\t\t\t\t Node Name")
	for _, value := range pods {
		fmt.Println(value.name + " \t" + value.nodeName)
	}
}
