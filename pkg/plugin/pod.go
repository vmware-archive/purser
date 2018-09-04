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

package plugin

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// Cost details
type Cost struct {
	totalCost   float64
	cpuCost     float64
	memoryCost  float64
	storageCost float64
}

// Pod Information
type Pod struct {
	name               string
	nodeName           string
	nodeCostPercentage float64
	cost               *Cost
	pvcs               []*string
}

// GetClusterPods returns the list of pods in cluster.
func GetClusterPods() []v1.Pod {
	pods, err := ClientSetInstance.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pods.Items
}

func getPodDetailsFromClient(podName string) *Pod {
	pod, err := ClientSetInstance.CoreV1().Pods("default").Get(podName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Node %s not found\n", podName)
		return nil
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting Node %s : %v\n", podName, statusError.ErrStatus.Message)
		return nil
	} else if err != nil {
		panic(err.Error())
	} else {
		podVolumes := []*string{}
		for j := 0; j < len(pod.Spec.Volumes); j++ {
			vol := pod.Spec.Volumes[j]
			if vol.PersistentVolumeClaim != nil {
				podVolumes = append(podVolumes, &vol.PersistentVolumeClaim.ClaimName)
			}
		}

		return &Pod{
			name:     pod.GetObjectMeta().GetName(),
			nodeName: pod.Spec.NodeName,
			pvcs:     podVolumes,
		}
	}
}

func getPodsForLabelThroughClient(label string) []*Pod {
	vals := strings.Split(label, "=")
	if len(vals) != 2 {
		panic("Label should be of form key=val")
	}

	m := map[string]string{vals[0]: vals[1]}
	pods, err := ClientSetInstance.CoreV1().Pods("").List(metav1.ListOptions{LabelSelector: labels.SelectorFromSet(m).String()})
	if err != nil {
		panic(err.Error())
	}

	return createPodObjects(pods)
}

func createPodObjects(pods *v1.PodList) []*Pod {
	ps := []*Pod{}

	for i := 0; i < len(pods.Items); i++ {
		pod := pods.Items[i]

		podVolumes := []*string{}
		for j := 0; j < len(pod.Spec.Volumes); j++ {
			vol := pod.Spec.Volumes[j]
			if vol.PersistentVolumeClaim != nil {
				podVolumes = append(podVolumes, &vol.PersistentVolumeClaim.ClaimName)
			}
		}

		p := Pod{
			name:     pod.GetObjectMeta().GetName(),
			nodeName: pod.Spec.NodeName,
			pvcs:     podVolumes,
		}

		ps = append(ps, &p)
	}
	return ps
}

func printPodsVerbose(pods []*Pod) {
	fmt.Printf("Cost Summary\n")
	totalCost := 0.0
	totalCPUCost := 0.0
	totalMemoryCost := 0.0
	totalStorageCost := 0.0

	for i := 0; i <= len(pods)-1; i++ {
		fmt.Printf("%-30s%s\n", "Pod Name:", pods[i].name)
		fmt.Printf("%-30s%s\n", "Node:", pods[i].nodeName)
		fmt.Printf("%-30s%.2f\n", "Pod Compute Cost Percentage:", pods[i].nodeCostPercentage*100.0)
		fmt.Printf("%-30s\n", "Persistent Volume Claims:")

		for j := 0; j <= len(pods[i].pvcs)-1; j++ {
			fmt.Printf("    %s\n", *pods[i].pvcs[j])

		}
		fmt.Printf("%-30s\n", "Cost:")
		fmt.Printf("    %-21s%f$\n", "Total Cost:", pods[i].cost.totalCost)
		fmt.Printf("    %-21s%f$\n", "Compute Cost:", pods[i].cost.cpuCost+pods[i].cost.memoryCost)
		fmt.Printf("    %-21s%f$\n", "Storage Cost:", pods[i].cost.storageCost)
		fmt.Printf("\n")

		totalCost += pods[i].cost.totalCost
		totalCPUCost += pods[i].cost.cpuCost
		totalMemoryCost += pods[i].cost.memoryCost
		totalStorageCost += pods[i].cost.storageCost
	}

	fmt.Printf("%-30s\n", "Total Cost Summary:")
	fmt.Printf("    %-21s%f$\n", "Total Cost:", totalCost)
	fmt.Printf("    %-21s%f$\n", "Compute Cost:", totalCPUCost+totalMemoryCost)
	fmt.Printf("    %-21s%f$\n", "Storage Cost:", totalStorageCost)
}
