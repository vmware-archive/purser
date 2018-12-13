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
	"fmt"
	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph/models"

	corev1 "k8s.io/api/core/v1"
)

var (
	podToSvcTable = make(map[string][]string)
	serviceMu     sync.Mutex
)

// PopulatePodToServiceTable populates the pod<->service map
func PopulatePodToServiceTable(svc corev1.Service, pods *corev1.PodList) {
	var podsXIDsInService []string
	serviceKey := svc.Namespace + KeySpliter + svc.Name

	serviceMu.Lock()
	for _, pod := range pods.Items {
		podKey := pod.Namespace + KeySpliter + pod.Name
		podToSvcTable[podKey] = append(podToSvcTable[podKey], serviceKey)
		podsXIDsInService = append(podsXIDsInService, podKey)
	}
	serviceMu.Unlock()

	err := models.StorePodServiceEdges(serviceKey, podsXIDsInService)
	if err != nil {
		log.Errorf("failed to store pod services edges for service: (%s), error: %s\n", svc.Name, err)
	}
}

// GenerateAndStoreSvcInteractions parses through pod interactions and generates a source to // destination service interaction.
func GenerateAndStoreSvcInteractions() {
	services, err := query.RetrieveAllServicesWithDstPods()
	if err != nil {
		log.Errorf("failed to fetch services: %s\n", err)
		return
	}

	for _, service := range services {
		destinationPods := getDestinationPods(service.Pod)
		destinationServices := getServicesXIDsFromPods(destinationPods)
		fmt.Printf("service: %v, xid: %s", service, service.Xid)
		err = models.StoreServicesInteraction(service.Xid, destinationServices)
		if err != nil {
			log.Errorf("failed to store services interactions, error %s\n", err)
		}
	}
}

func getDestinationPods(podsInService []*models.Pod) []*models.Pod {
	var destinationPods []*models.Pod
	for _, pod := range podsInService {
		destinationPods = append(destinationPods, pod.Pods...)
	}
	return destinationPods
}

func getServicesXIDsFromPods(pods []*models.Pod) []string {
	var servicesXIDs []string
	duplicateChecker := make(map[string]bool)
	for _, pod := range pods {
		svcsXIDs := podToSvcTable[pod.Xid]
		for _, svcXID := range svcsXIDs {
			if _, isPresent := duplicateChecker[svcXID]; !isPresent {
				duplicateChecker[svcXID] = true
				servicesXIDs = append(servicesXIDs, svcXID)
			}
		}
	}
	return servicesXIDs
}
