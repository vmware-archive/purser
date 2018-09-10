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

package controller

import (
	"flag"
	"fmt"
	"os"

	"github.com/vmware/purser/pkg/plugin"
	"github.com/vmware/purser/pkg/plugin/client"
	"github.com/vmware/purser/pkg/plugin/crd"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClientConfig returns rest config, if path not specified assume in cluster config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

// GetAPIExtensionClient returns the CRD client instance
func GetAPIExtensionClient() *client.Crdclient {
	kubeconf := flag.String("kubeconf", os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), "path to Kubernetes config file")
	flag.Parse()

	config, err := GetClientConfig(*kubeconf)
	if err != nil {
		panic(err.Error())
	}

	// create clientset and create our CRD, this only need to run once
	_, err = apiextcs.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a new clientset which include our CRD schema
	crdcs, scheme, err := crd.NewClient(config)
	if err != nil {
		panic(err)
	}

	// Create a CRD client interface
	crdclient := client.CrdClient(crdcs, scheme, "default")

	return crdclient
}

// ListCrdInstances displays the list of CRD instances.
func ListCrdInstances(crdclient *client.Crdclient) {
	// List all Example objects
	items, err := crdclient.List(meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("List:\n%v\n", items)
}

// GetCrdByName returns the CRD group.
func GetCrdByName(crdclient *client.Crdclient, groupName string) *crd.Group {
	group, err := crdclient.Get(groupName)

	if err == nil {
		return group
	} else if apierrors.IsNotFound(err) {
		return nil
	} else {
		panic(err)
	}
}

// PrintGroup displays the group information.
func PrintGroup(group *crd.Group) {
	pitGroupMetrics, mtdGroupMetrics, cost := plugin.GetGroupDetails(group)

	fmt.Printf("%-30s             %s\n", "Group Name:", group.Name)
	fmt.Println()
	fmt.Println("Point in Time Resource Stats:")
	fmt.Printf("             %-30s%.2f\n", "CPU Limit(vCPU):", pitGroupMetrics.CPULimit)
	fmt.Printf("             %-30s%.2f\n", "Memory Limit(GB):", pitGroupMetrics.MemoryLimit)
	fmt.Printf("             %-30s%.2f\n", "CPU Request(vCPU):", pitGroupMetrics.CPURequest)
	fmt.Printf("             %-30s%.2f\n", "Memory Request(GB):", pitGroupMetrics.MemoryRequest)

	fmt.Println()
	fmt.Printf("%-30s\n", "Month to Date Active Resource Stats:")
	fmt.Printf("             %-30s%.2f\n", "CPU Request(vCPU-hours):", mtdGroupMetrics.CPURequest)
	fmt.Printf("             %-30s%.2f\n", "Memory Request(GB-hours):", mtdGroupMetrics.MemoryRequest)
	fmt.Printf("             %-30s%.2f\n", "Storage Claimed(GB-hours):", mtdGroupMetrics.StorageClaimed)

	fmt.Println()
	fmt.Printf("%-30s\n", "Month to Date Cost Stats:")
	fmt.Printf("             %-30s%.2f\n", "CPU Cost($):", cost.CPUCost)
	fmt.Printf("             %-30s%.2f\n", "Memory Cost($):", cost.MemoryCost)
	fmt.Printf("             %-30s%.2f\n", "Storage Cost($):", cost.StorageCost)
	fmt.Printf("             %-30s%.2f\n", "Total Cost($):", cost.TotalCost)
}
