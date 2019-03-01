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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// SyncCluster will handle missed events
func SyncCluster(kubeClient *kubernetes.Clientset) {
	endTime := time.Now().Format(time.RFC3339)
	syncPods(kubeClient, endTime)
}

// syncPods handles missed creation and deletion of pod events
func syncPods(kubeClient *kubernetes.Clientset, endTime string) {
	logrus.Infof("[SYNC] started syncing pods")
	livePodsFromDgraph := query.RetrieveAllLivePods()
	logrus.Infof("[SYNC] number of livePodsFromDgraph: %d", len(livePodsFromDgraph))

	podsInCluster := utils.RetrievePodList(kubeClient, v1.ListOptions{})
	if podsInCluster == nil {
		logrus.Errorf("[SYNC] got no podsInCluster, aborting sync")
		return
	}
	logrus.Infof("[SYNC] number of pods in cluster: %d", len(podsInCluster.Items))

	handleDeadPodsAndNewPods(livePodsFromDgraph, podsInCluster, endTime)
	logrus.Infof("[SYNC] finished syncing of pods")
}

// if dead pods end time isn't updated in dgraph this function will update it
// if an pod creation event is missed then this function will create a new pod in dgraph
func handleDeadPodsAndNewPods(livePodsFromDgraph []models.Pod, podsInCluster *corev1.PodList, endTime string) {
	// create a map from pod xid to k8s pod pointer
	podXIDToPod := make(map[string]*corev1.Pod)
	for _, pod := range podsInCluster.Items {
		xid := pod.Namespace + ":" + pod.Name
		if _, isPresent := podXIDToPod[xid]; !isPresent {
			podXIDToPod[xid] = &pod
		}
	}

	var deadPods []models.Pod
	podsXIDs := make(map[string]bool) // create a map from pod xid to bool
	for _, pod := range livePodsFromDgraph {
		if _, isAlive := podXIDToPod[pod.Xid]; !isAlive {
			// pod is in dgraph but not in cluster -> pod got deleted but end time not updated in dgraph ->
			// missed pod deletion event -> update pod in dgraph with end time
			deadPod := models.Pod{
				ID:      dgraph.ID{Xid: pod.Xid + endTime, UID: pod.UID},
				EndTime: endTime,
				Name:    "pod-" + pod.Name + "*" + endTime,
			}
			deadPods = append(deadPods, deadPod)
		}

		if _, isPresent := podsXIDs[pod.Xid]; !isPresent {
			podsXIDs[pod.Xid] = true
		}
	}

	// update deletion time stamps for dead pods
	_, err := dgraph.MutateNode(deadPods, dgraph.UPDATE)
	if err != nil {
		logrus.Errorf("[SYNC] unable to update deleted pods with end time: # deleted pods: %d, err: %v", len(deadPods), err)
	}

	// create new pod if it isn't in dgraph
	for podXID, pod := range podXIDToPod {
		if _, isPresent := podsXIDs[podXID]; !isPresent {
			// pod is in cluster but not in dgraph -> missed pod creation event -> create new pod in dgraph
			err = models.StorePod(*pod)
			if err != nil {
				logrus.Errorf("[SYNC] Error while persisting pod: %s, err: %v", podXID, err)
			}
		}
	}
}
