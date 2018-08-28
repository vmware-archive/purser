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
	"os"

	"strings"

	"github.com/vmware/purser/pkg/purser_plugin"
	"github.com/vmware/purser/pkg/purser_plugin/client"
	"github.com/vmware/purser/pkg/purser_plugin/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var crdclient *client.Crdclient

func main() {
	inputs := os.Args[1:]
	inputs = inputs[1:]
	if len(inputs) == 4 && inputs[0] == "get" && inputs[1] == "cost" {
		if inputs[2] == "label" {
			purser_plugin.GetPodsCostForLabel(inputs[3])
		} else if inputs[2] == "pod" {
			purser_plugin.GetPodCost(inputs[3])
		} else if inputs[2] == "node" {
			purser_plugin.GetAllNodesCost()
		} else {
			printHelp()
		}
	} else if len(inputs) == 4 && inputs[0] == "get" && inputs[1] == "resources" {
		if inputs[2] == "namespace" {
			group := controller.GetCrdByName(crdclient, inputs[3])
			if group != nil {
				controller.PrintGroup(group)
			} else {
				fmt.Printf("Group %s is not present\n", inputs[3])
			}
		} else if inputs[2] == "label" {
			if !strings.Contains(inputs[3], "=") {
				printHelp()
			}
			group := controller.GetCrdByName(crdclient, createGroupNameFromLabel(inputs[3]))
			if group != nil {
				controller.PrintGroup(group)
			} else {
				fmt.Printf("Group %s is not present\n", inputs[3])
			}
		} else {
			printHelp()
		}
	} else if len(inputs) == 2 && inputs[0] == "get" {
		if inputs[1] == "summary" {
			purser_plugin.GetClusterSummary()
		} else if inputs[1] == "savings" {
			purser_plugin.GetSavings()
		} else if inputs[1] == "user-costs" {
			cpuCostPerCPUPerHour, memCostPerGBPerHour, storageCostPerGBPerHour := purser_plugin.GetUserCosts()
			fmt.Printf("cpu cost per CPU per hour:\t %f$\nmem cost per GB per hour:\t %f$\nstorage cost per GB per hour:\t %f$\n",
				cpuCostPerCPUPerHour,
				memCostPerGBPerHour,
				storageCostPerGBPerHour)
		} else {
			printHelp()
		}
	} else if len(inputs) == 2 && inputs[0] == "set" {
		if inputs[1] == "user-costs" {
			fmt.Printf("Enter CPU cost per cpu per hour:\t ")
			var cpuCostPerCPUPerHour string
			fmt.Scan(&cpuCostPerCPUPerHour)

			fmt.Printf("Enter Memory cost per GB per hour:\t ")
			var memCostPerGBPerHour string
			fmt.Scan(&memCostPerGBPerHour)

			fmt.Printf("Enter Storage cost per GB per hour:\t ")
			var storageCostPerGBPerHour string
			fmt.Scan(&storageCostPerGBPerHour)

			purser_plugin.SaveUserCosts(cpuCostPerCPUPerHour, memCostPerGBPerHour, storageCostPerGBPerHour)
		} else {
			printHelp()
		}
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

func main2() {
	//controller.ListCrdInstances(crdclient)
	groupName := "apundlik1"
	group := controller.GetCrdByName(crdclient, groupName)
	//fmt.Println(group)
	if group != nil {
		controller.PrintGroup(group)
	} else {
		fmt.Printf("Group %s is not present\n", groupName)
	}
}

func init2() {
	crdclient = controller.GetApiExtensionClient()
}

func init() {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"))
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	purser_plugin.ProvideClientSetInstance(clientset)

	// Crd client
	crdclient = controller.GetApiExtensionClient()
}

func printHelp() {
	fmt.Printf("Try one of the following commands...\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get summary\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get resources namespace <Namespace>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get resources label <key=val>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get cost label <key=val>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get cost pod <pod name>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get cost node <node name>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser set user-costs\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin purser get user-costs\n")
}
