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

	wg.Add(podsCount)
	{
		for index, pod := range pods.Items {
			log.Debugf("Processing Pod (%d/%d) ... ", index+1, podsCount)

			go func(pod corev1.Pod, index int) {
				defer wg.Done()

				containers := pod.Spec.Containers
				podInteractions := processContainerDetails(conf, pod, containers)
				linker.UpdatePodToPodTable(podInteractions)
				log.Debugf("Finished processing Pod (%d/%d)", index+1, podsCount)
			}(pod, index)
		}
	}
	wg.Wait()
}
