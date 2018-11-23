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

package processor

import (
	"sync"

	"github.com/vmware/purser/pkg/controller"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/discovery/linker"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const maxGoRoutines = 20

var wg sync.WaitGroup

// ProcessPodInteractions fetches details of all the running processes in each container of
// each pod in a given namespace and generates a 1:1 mapping between the communicating pods.
func ProcessPodInteractions(conf controller.Config) {
	k8sPods := RetrievePodList(conf.Kubeclient, metav1.ListOptions{})

	linker.PopulatePodIPTable(k8sPods)
	processPodDetails(conf, k8sPods)

	linker.GenerateAndStorePodInteractions()
	log.Infof("Successfully generated Pod To Pod mapping.")
}

func processPodDetails(conf controller.Config, pods *corev1.PodList) {
	podsCount := len(pods.Items)
	log.Infof("Processing total of (%d) Pods.", podsCount)

	freeRoutines := maxGoRoutines
	numChannelsReceived := 0
	ch := make(chan int, 1)

	wg.Add(podsCount)
	{
		for index, pod := range pods.Items {
			log.Debugf("Processing Pod: (%s), (%d/%d) ... ", pod.Name, index+1, podsCount)

			// wait for a free goroutine
			if freeRoutines < 1 {
				// wait for a go routine to send to channel i.e, it will wait until a go routine finishes.
				<-ch
				numChannelsReceived++
			}

			// decrease 1 from freeRoutines before starting a new go routine
			freeRoutines--
			go func(pod corev1.Pod, index int) {
				defer wg.Done()

				containers := pod.Spec.Containers
				interactions := processContainerDetails(conf, pod, containers)
				linker.UpdatePodToPodTable(interactions.PodInteractions)
				linker.StoreProcessInteractions(interactions.ContainerProcessInteraction, interactions.ProcessToPodInteraction,
					pod.GetCreationTimestamp().Time)
				log.Debugf("Finished processing Pod: (%s), (%d/%d)", pod.Name, index+1, podsCount)

				// increase 1 from freeRoutines after processing a pod.
				freeRoutines++
				// send 1 to channel
				ch <- 1
			}(pod, index)
		}
	}
	// receive channel from remaining go routines
	for i := 0; i < podsCount-numChannelsReceived; i++ {
		<-ch
	}
	wg.Wait()
	close(ch)
}
