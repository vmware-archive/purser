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

package linker

import (
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"

	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

// InteractionsWrapper ...
type InteractionsWrapper struct {
	PodInteractions             map[string](map[string]float64)
	ProcessToPodInteraction     map[string](map[string]bool)
	ContainerProcessInteraction map[string][]string
}

// podIPTable: maps pod name with pod IP address
// podToPodTable: maps src pod to the interacting dest pod along with the interaction frequency count.
var (
	podIPTable    = make(map[string]string)
	podToPodTable = make(map[string](map[string]float64))
)

var (
	mu sync.Mutex
)

const (
	// KeySpliter splits the key into resource namespace and name used for processing Xids
	KeySpliter = ":"
)

// PopulatePodIPTable populates the podIP<->podName map
func PopulatePodIPTable(pods *corev1.PodList) {
	for _, pod := range pods.Items {
		podIP := pod.Status.PodIP
		podIPTable[podIP] = pod.Namespace + KeySpliter + pod.Name
	}
}

// GenerateAndStorePodInteractions generates source to destination Pod mapping and stores it in Dgraph.
func GenerateAndStorePodInteractions() {
	log.Info("Storing Pod Interactions ....")
	for srcPodName, communication := range podToPodTable {
		dstPods := []string{}
		counts := []float64{}
		for dstPodName, count := range communication {
			dstPods = append(dstPods, dstPodName)
			counts = append(counts, count)
		}
		err := models.StorePodsInteraction(srcPodName, dstPods, counts)
		if err != nil {
			log.Errorf("failed to store pod interaction in Dgraph %v", err)
		}
	}
	log.Info("Finished storing pod interactions.")
}

// PopulateMappingTables updates PodToPodTable
func PopulateMappingTables(tcpDump []string, pod corev1.Pod, process Process, containerName string, interactions *InteractionsWrapper) {
	podXID := pod.Namespace + KeySpliter + pod.Name
	containerXID := podXID + KeySpliter + containerName
	procXID := containerXID + KeySpliter + process.ID + KeySpliter + process.Name
	populateContainerProcessTable(containerXID, procXID, interactions)
	for _, address := range tcpDump {
		address := strings.Split(address, KeySpliter)
		srcIP, dstIP := address[0], address[2]
		srcName, dstName := podIPTable[srcIP], podIPTable[dstIP]
		updatePodInteractions(srcName, dstName, interactions)
		updatePodProcessInteractions(procXID, dstName, interactions)
	}
}

func updatePodInteractions(srcName, dstName string, interactions *InteractionsWrapper) {
	if dstName != "" && srcName != "" {
		log.Debugf("pod interactions srcName: (%s), dstName: (%s)", srcName, dstName)
		if _, ok := interactions.PodInteractions[srcName]; !ok {
			interactions.PodInteractions[srcName] = make(map[string]float64)
		}

		if _, isPresent := interactions.PodInteractions[srcName][dstName]; !isPresent {
			interactions.PodInteractions[srcName][dstName] = 1
		} else {
			interactions.PodInteractions[srcName][dstName]++
		}
	}
}

// UpdatePodToPodTable ...
func UpdatePodToPodTable(podInteractions map[string](map[string]float64)) {
	mu.Lock()
	for srcPod, interaction := range podInteractions {
		if _, ok := podToPodTable[srcPod]; !ok {
			podToPodTable[srcPod] = interaction
		} else {
			for dstPod, count := range interaction {
				if _, isPresent := podToPodTable[srcPod][dstPod]; !isPresent {
					podToPodTable[srcPod][dstPod] = count
				} else {
					podToPodTable[srcPod][dstPod] += count
				}
			}
		}
	}
	mu.Unlock()
}
