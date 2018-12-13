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

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/purser/pkg/controller"
	"github.com/vmware/purser/pkg/controller/discovery/linker"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

var svcwg sync.WaitGroup

// ProcessServiceInteractions parses through the list of services and it's associated pods to
// generate a 1:1 mapping between the communicating services.
func ProcessServiceInteractions(conf controller.Config) {
	services := RetrieveServiceList(conf.Kubeclient, metav1.ListOptions{})
	log.Debugf("service list retrieved: %v", services.Items)
	if services == nil {
		return
	}

	processServiceDetails(conf.Kubeclient, services)
	linker.GenerateAndStoreSvcInteractions()

	log.Infof("Successfully generated Services To Services mapping.")
}

func processServiceDetails(client *kubernetes.Clientset, services *corev1.ServiceList) {
	svcCount := len(services.Items)
	log.Infof("Processing total of (%d) Services.", svcCount)

	// TODO: restrict number of go routines, reason: decrease load on Kubernetes api server
	svcwg.Add(svcCount)
	{
		for index, svc := range services.Items {
			log.Debugf("Processing Services (%d/%d): %s ", index+1, svcCount, svc.GetName())

			go func(svc corev1.Service, index int) {
				defer svcwg.Done()

				selectorSet := labels.Set(svc.Spec.Selector)
				log.Debugf("service: %s, selectorSet: (%v)", svc.Name, selectorSet)
				if selectorSet != nil {
					options := metav1.ListOptions{
						LabelSelector: selectorSet.AsSelector().String(),
					}
					pods := RetrievePodList(client, options)
					if pods != nil {
						log.Debugf("service: %s, podsCount: %d", svc.Name, len(pods.Items))
						linker.PopulatePodToServiceTable(svc, pods)
					}
				}

				log.Debugf("Finished processing Services (%d/%d)", index+1, svcCount)
			}(svc, index)
		}
	}
	svcwg.Wait()
}
