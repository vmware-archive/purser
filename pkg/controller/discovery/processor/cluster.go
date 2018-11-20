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
	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RetrievePodList returns list of pods in the given namespace.
func RetrievePodList(client *kubernetes.Clientset, options metav1.ListOptions) *corev1.PodList {
	pods, err := client.CoreV1().Pods(metav1.NamespaceAll).List(options)
	if err != nil {
		log.Errorf("failed to retrieve pods: %v", err)
	}
	return pods
}

// RetrieveServiceList returns list of services in the given namespace.
func RetrieveServiceList(client *kubernetes.Clientset, options metav1.ListOptions) *corev1.ServiceList {
	services, err := client.CoreV1().Services(metav1.NamespaceAll).List(options)
	if err != nil {
		log.Errorf("failed to retrieve services: %v", err)
	}
	return services
}
