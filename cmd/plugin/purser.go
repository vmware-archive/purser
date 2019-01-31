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

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/purser/pkg/client"
	groups_client_v1 "github.com/vmware/purser/pkg/client/clientset/typed/groups/v1"
	"github.com/vmware/purser/pkg/plugin"
	"github.com/vmware/purser/pkg/utils"
)

const (
	pluginVersion = "version v1.0.0"
)

var (
	groupClient *groups_client_v1.GroupClient

	// Variables used for cmd interface
	kubeconfig string
	info       string
	version    string

	description   = fmt.Sprintf("Purser gives cost insights of kubernetes deployments.\n\n")
	usage         = fmt.Sprintf("Usage:\n  kubectl plugin purser [options] <command> <args>\n\n")
	supportedCmds = fmt.Sprintf("The supported commands are:\n  get  Get resource information.\n  set  Set resource information.\n\n")

	optionHelp       = fmt.Sprintf("\n  --info            Show more details about the plugin.")
	optionKubeConfig = fmt.Sprintf("\n  --kubeconfig      Absolute path for the kube config file.")
	optionVersion    = fmt.Sprintf("\n  --version         Show plugin version.")
	options          = fmt.Sprintf("options:%s%s%s\n\n", optionHelp, optionKubeConfig, optionVersion)

	kubecltOption = fmt.Sprintf("\nUse \"kubectl options\" for a list of global command-line options (applies to all commands).\n\n")
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), "path to Kubernetes config file")

	flag.StringVar(&info, "info", os.Getenv("KUBECTL_PLUGINS_LOCAL_FLAG_INFO"), "Show help documentation")
	flag.StringVar(&version, "version", os.Getenv("KUBECTL_PLUGINS_LOCAL_FLAG_VERSION"), "Show version number")

	flag.Usage = func() {
		_, err := fmt.Fprint(flag.CommandLine.Output(), description)
		if err != nil {
			log.Fatal(err)
		}

		_, err = fmt.Fprint(flag.CommandLine.Output(), usage)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fmt.Fprint(flag.CommandLine.Output(), supportedCmds)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fmt.Fprint(flag.CommandLine.Output(), options)
		if err != nil {
			log.Fatal(err)
		}

		_, err = fmt.Fprint(flag.CommandLine.Output(), "Example(s):\n\n")
		if err != nil {
			log.Fatal(err)
		}

		printHelp()
		_, err = fmt.Fprint(flag.CommandLine.Output(), kubecltOption)
		if err != nil {
			log.Fatal(err)
		}
	}

	if version != "" {
		fmt.Println(pluginVersion)
		os.Exit(0)
	}

	if info != "" {
		flag.Usage()
		os.Exit(0)
	}

	config, err := utils.GetKubeconfig(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	plugin.ProvideClientSetInstance(utils.GetKubeclient(config))

	client, clusterConfig := client.GetAPIExtensionClient(kubeconfig)
	groupClient = groups_client_v1.NewGroupClient(client, clusterConfig)
}

func main() {
	inputs := os.Args[2:] // index 1 is empty
	if len(inputs) == 4 && inputs[0] == Get {
		computeMetricInsight(inputs)
	} else if len(inputs) == 2 {
		computeStats(inputs)
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

func computeCost(inputs []string) {
	switch inputs[2] {
	case Label:
		plugin.GetPodsCostForLabel(inputs[3])
	case Pod:
		plugin.GetPodCost(inputs[3])
	case Node:
		plugin.GetAllNodesCost()
	default:
		printHelp()
	}
}

func fetchResource(inputs []string) {
	switch inputs[2] {
	case Namespace:
		group := plugin.GetGroupByName(groupClient, inputs[3])
		if group != nil {
			plugin.PrintGroup(group)
		} else {
			fmt.Printf("Group %s is not present\n", inputs[3])
		}
	case Label:
		if !strings.Contains(inputs[3], "=") {
			printHelp()
		}
		group := plugin.GetGroupByName(groupClient, createGroupNameFromLabel(inputs[3]))
		if group != nil {
			plugin.PrintGroup(group)
		} else {
			fmt.Printf("Group %s is not present\n", inputs[3])
		}
	case Group:
		group := plugin.GetGroupByName(groupClient, inputs[3])
		if group != nil {
			plugin.PrintGroup(group)
		} else {
			fmt.Printf("No group with name: %s\n", inputs[3])
		}
	default:
		printHelp()
	}
}

func createGroupNameFromLabel(input string) string {
	inp := strings.Split(input, "=")
	key, val := inp[0], inp[1]
	groupName := key + "." + val
	if strings.Contains(groupName, "/") {
		groupName = strings.Replace(groupName, "/", "-", -1)
	}
	return strings.ToLower(groupName)
}

func computeStats(inputs []string) {
	switch inputs[0] {
	case Get:
		getStats(inputs)
	case Set:
		inputUserCosts(inputs)
	default:
		printHelp()
	}
}

func getStats(inputs []string) {
	switch inputs[1] {
	case "summary":
		plugin.GetClusterSummary()
	case "savings":
		plugin.GetSavings()
	case "user-costs":
		price := plugin.GetUserCosts()
		fmt.Printf("cpu cost per CPU per hour:\t %f$\nmem cost per GB per hour:\t %f$\nstorage cost per GB per hour:\t %f$\n",
			price.CPU,
			price.Memory,
			price.Storage)
	default:
		printHelp()
	}
}

func inputUserCosts(inputs []string) {
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

func printHelp() {
	pluginExt := "kubectl --kubeconfig=<absolute path to config> plugin purser "

	fmt.Println("Try one of the following commands...")
	fmt.Println(pluginExt + "get summary")
	fmt.Println(pluginExt + "get resources group <group-name>")
	fmt.Println(pluginExt + "get cost label <key=val>")
	fmt.Println(pluginExt + "get cost pod <pod name>")
	fmt.Println(pluginExt + "get cost node all")
	fmt.Println(pluginExt + "set user-costs")
	fmt.Println(pluginExt + "get user-costs")
	fmt.Println(pluginExt + "get savings")
}

func logError(err error) {
	if err != nil {
		log.Printf("failed to read user input %+v", err)
	}
}
