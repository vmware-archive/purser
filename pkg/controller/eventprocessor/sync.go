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

package eventprocessor

import (
	"time"

	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
	"github.com/vmware/purser/pkg/controller/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

// SyncCluster will handle missed events
func SyncCluster(kubeClient *kubernetes.Clientset) {
	endTime := time.Now().Format(time.RFC3339)
	syncPodsAndContainers(kubeClient, endTime)
}

// syncPodsAndContainers gets pods through k8s api and from dgraph.
// Compares between these two sets and add endTime to those pods in dgraph but not in cluster.
// Creates pod in dgraph which are in cluster but not in dgraph.
func syncPodsAndContainers(kubeClient *kubernetes.Clientset, endTime string) {
	logrus.Infof("[SYNC] started syncing pods")
	livePodsFromDgraph := query.RetrieveAllLivePods()
	logrus.Infof("[SYNC] number of livePodsFromDgraph: %d", len(livePodsFromDgraph))

	podsInCluster := utils.RetrievePodList(kubeClient, v1.ListOptions{})
	if podsInCluster == nil {
		logrus.Errorf("[SYNC] got no podsInCluster, aborting sync")
		return
	}
	logrus.Debugf("[SYNC] number of pods in cluster: %d", len(podsInCluster.Items))

	var updatedPods []models.Pod
	for _, pod := range livePodsFromDgraph {
		logrus.Debugf("[SYNC] Name: %s, XID: %s, UID: %v", pod.Name, pod.Xid, pod.UID)
		pod.EndTime = endTime
		updatedPods = append(updatedPods, pod)

		podData := models.RetrievePodWithContainers(pod.Xid)
		models.SoftDeleteContainersInTerminatedPod(podData.Containers, endTime)
	}
	_, err := dgraph.MutateNode(updatedPods, dgraph.UPDATE)
	if err != nil {
		logrus.Errorf("[SYNC] unable to update pods with end time: %v", err)
	}

	// stores new pod in dgraph if not persisted in it.
	// updates (deletes existing pod data and creates a new pod) if it is already persisted in dgraph.
	for _, pod := range podsInCluster.Items {
		logrus.Debugf("[SYNC] storing/updating pod: (name: %s, namespace: %s)", pod.Name, pod.Namespace)
		err := models.StorePod(pod)
		if err != nil {
			logrus.Errorf("[SYNC] unable to store/update pod: (name: %s, namespace: %s), err: %v", pod.Name, pod.Namespace, err)
		}
	}
}
