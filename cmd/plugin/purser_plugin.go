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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"strings"

	"github.com/vmware/purser/pkg/plugin"
	"github.com/vmware/purser/pkg/plugin/client"
	"github.com/vmware/purser/pkg/plugin/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var crdclient *client.Crdclient

func main() {
	inputs := os.Args[1:]
	inputs = inputs[1:]
	if len(inputs) == 4 && inputs[0] == Get {
		computeMetricInsight(inputs)
	} else if len(inputs) == 2 {
		performAction(inputs)
	} else {
		printHelp()
	}
}

func computeMetricInsight(inputs []string) {
	switch inputs[1] {
	case Cost:
		computeCost(inputs)
	case Resources:
		fetchResource(inputs)
	}
}

func performAction(inputs []string) {
	switch inputs[0] {
	case Get:
		getStats(inputs)
	case Set:
		setStats(inputs)
	default:
		printHelp()
	}
}

func setStats(inputs []string) {
	if inputs[1] == "user-costs" {
		fmt.Printf("Enter CPU cost per cpu per hour:\t ")
		var cpuCostPerCPUPerHour string
		_, err := fmt.Scan(&cpuCostPerCPUPerHour)
		logError(err)

		fmt.Printf("Enter Memory cost per GB per hour:\t ")
		var memCostPerGBPerHour string
		_, err = fmt.Scan(&memCostPerGBPerHour)
		logError(err)

		fmt.Printf("Enter Storage cost per GB per hour:\t ")
		var storageCostPerGBPerHour string
		_, err = fmt.Scan(&storageCostPerGBPerHour)
		logError(err)

		plugin.SaveUserCosts(cpuCostPerCPUPerHour, memCostPerGBPerHour, storageCostPerGBPerHour)
	} else {
		printHelp()
	}
}

func getStats(inputs []string) {
	if inputs[1] == "summary" {
		plugin.GetClusterSummary()
	} else if inputs[1] == "savings" {
		plugin.GetSavings()
	} else if inputs[1] == "user-costs" {
		price := plugin.GetUserCosts()
		fmt.Printf("cpu cost per CPU per hour:\t %f$\nmem cost per GB per hour:\t %f$\nstorage cost per GB per hour:\t %f$\n",
			price.CPU,
			price.Memory,
			price.Storage)
	} else {
		printHelp()
	}
}

func fetchResource(inputs []string) {
	if inputs[2] == Namespace {
		group := controller.GetCrdByName(crdclient, inputs[3])
		if group != nil {
			controller.PrintGroup(group)
		} else {
			fmt.Printf("Group %s is not present\n", inputs[3])
		}
	} else if inputs[2] == Label {
		if !strings.Contains(inputs[3], "=") {
			printHelp()
		}
		group := controller.GetCrdByName(crdclient, createGroupNameFromLabel(inputs[3]))
		if group != nil {
			controller.PrintGroup(group)
		} else {
			fmt.Printf("Group %s is not present\n", inputs[3])
		}
	} else if inputs[2] == Group {
		group := controller.GetCrdByName(crdclient, inputs[3])
		if group != nil {
			controller.PrintGroup(group)
		} else {
			fmt.Printf("No group with name: %s\n", inputs[3])
		}
	} else {
		printHelp()
	}
}

func computeCost(inputs []string) {
	if inputs[2] == Label {
		plugin.GetPodsCostForLabel(inputs[3])
	} else if inputs[2] == Pod {
		plugin.GetPodCost(inputs[3])
	} else if inputs[2] == Node {
		plugin.GetAllNodesCost()
	} else {
		printHelp()
	}
}

func createGroupNameFromLabel(input string) string {
	inp := strings.Split(input, "=")
	key := inp[0]
	val := inp[1]
	groupName := key + "." + val
	if strings.Contains(groupName, "/") {
		groupName = strings.Replace(groupName, "/", "-", -1)
	}
	groupName = strings.ToLower(groupName)
	return groupName
}

func init() {
	kubeconfig := flag.String("kubeconfig", os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"))
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	plugin.ProvideClientSetInstance(clientset)

	// Crd client
	crdclient = controller.GetAPIExtensionClient()
}

func printHelp() {
	fmt.Printf("Try one of the following commands...\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get summary\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get resources group <group-name>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get cost label <key=val>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get cost pod <pod name>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get cost node all\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser set user-costs\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get user-costs\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get savings\n")
}

func logError(err error) {
	if err != nil {
		log.Printf("%+v", err)
	}
}
