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

package utils

import (
	log "github.com/Sirupsen/logrus"
	storagev1 "k8s.io/api/storage/v1"

	groupsv1 "github.com/vmware/purser/pkg/apis/groups/v1"
	groups "github.com/vmware/purser/pkg/client/clientset/typed/groups/v1"

	api_v1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	StorageDefault = "purser-default"
)

// RetrievePodList returns list of pods in the given namespace.
func RetrievePodList(client *kubernetes.Clientset, options metav1.ListOptions) *corev1.PodList {
	pods, err := client.CoreV1().Pods(metav1.NamespaceAll).List(options)
	if err != nil {
		log.Errorf("failed to retrieve pods: %v", err)
		return nil
	}
	return pods
}

// RetrieveServiceList returns list of services in the given namespace.
func RetrieveServiceList(client *kubernetes.Clientset, options metav1.ListOptions) *corev1.ServiceList {
	services, err := client.CoreV1().Services(metav1.NamespaceAll).List(options)
	if err != nil {
		log.Errorf("failed to retrieve services: %v", err)
		return nil
	}
	return services
}

// RetrieveGroupList returns list of group CRDs in the given namespace.
func RetrieveGroupList(groupClient *groups.GroupClient, options metav1.ListOptions) *groupsv1.GroupList {
	crdGroups, err := groupClient.List(options)
	if err != nil {
		log.Errorf("failed to retrieve group list: %v ", err)
		return nil
	}
	return crdGroups
}

// RetrieveStorageClass returns storage class with the given name. Nil if error is encountered
func RetrieveStorageClass(client *kubernetes.Clientset, options metav1.GetOptions, name string) (*storagev1.StorageClass, error) {
	storageClass, err := client.StorageV1().StorageClasses().Get(name, options)
	if err != nil {
		log.Errorf("failed to retrieve storage class: %s, err: %v", name, err)
		return nil, err
	}
	return storageClass, err
}

// GetStorageType ...
// input: persistent volume
// output: the type(final) of PV's storage class
// i.e., if PV has storage class A, A is of type B(storage class) and so on..
// until a storage class X is of its own type X. Then this function returns the final type of PV's storage as X
//
// "purser-default" is returned in special cases:
// 1. if A is of type B, if B is of type A (i.e., if a cycle is found)
// 2. an error is encountered
// 3. if A is not having any type i.e., "" (empty string case)
func GetFinalStorageTypeOfPV(pv api_v1.PersistentVolume, client *kubernetes.Clientset) string {
	cycleChecker := make(map[string]bool)
	log.Debugf("PV: %s, storageClass: %s", pv.Name, pv.Spec.StorageClassName)
	return getFinalTypeOfStorageClass(client, pv.Spec.StorageClassName, cycleChecker)
}

// getFinalTypeOfStorageClass
// this is helper function for func getStorageType
func getFinalTypeOfStorageClass(client *kubernetes.Clientset, storageClassName string, cycleChecker map[string]bool) string {
	if _, isVisited := cycleChecker[storageClassName]; isVisited {
		return StorageDefault
	} else {
		cycleChecker[storageClassName] = true
	}

	storageClass, err := RetrieveStorageClass(client, metav1.GetOptions{}, storageClassName)
	if err != nil {
		return StorageDefault
	}

	storageType := storageClass.Parameters["type"]
	if storageType == "" {
		return StorageDefault
	} else if storageType == storageClassName {
		return storageClassName
	}
	return getFinalTypeOfStorageClass(client, storageType, cycleChecker)
}
