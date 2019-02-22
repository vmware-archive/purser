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

package models

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/purser/pkg/controller/dgraph"

	api_v1 "k8s.io/api/core/v1"
)

// Dgraph Model Constants
const (
	IsPod = "isPod"
)

// Pod schema in dgraph
type Pod struct {
	dgraph.ID
	IsPod          bool                     `json:"isPod,omitempty"`
	Name           string                   `json:"name,omitempty"`
	StartTime      string                   `json:"startTime,omitempty"`
	EndTime        string                   `json:"endTime,omitempty"`
	Containers     []*Container             `json:"containers,omitempty"`
	Pods           []*Pod                   `json:"pod,omitempty"`
	Count          float64                  `json:"pod|count,omitempty"`
	Node           *Node                    `json:"node,omitempty"`
	Namespace      *Namespace               `json:"namespace,omitempty"`
	Deployment     *Deployment              `json:"deployment,omitempty"`
	Replicaset     *Replicaset              `json:"replicaset,omitempty"`
	Statefulset    *Statefulset             `json:"statefulset,omitempty"`
	Daemonset      *Daemonset               `json:"daemonset,omitempty"`
	Job            *Job                     `json:"job,omitempty"`
	Pvcs           []*PersistentVolumeClaim `json:"pvc,omitempty"`
	CPURequest     float64                  `json:"cpuRequest,omitempty"`
	CPULimit       float64                  `json:"cpuLimit,omitempty"`
	MemoryRequest  float64                  `json:"memoryRequest,omitempty"`
	MemoryLimit    float64                  `json:"memoryLimit,omitempty"`
	StorageRequest float64                  `json:"storageRequest,omitempty"`
	Type           string                   `json:"type,omitempty"`
	Cid            []Service                `json:"cid,omitempty"`
	Labels         []*Label                 `json:"label,omitempty"`
	CPUPrice       float64                  `json:"cpuPrice,omitempty"`
	MemoryPrice    float64                  `json:"memoryPrice,omitempty"`
}

// Metrics ...
type Metrics struct {
	CPURequest    float64
	CPULimit      float64
	MemoryRequest float64
	MemoryLimit   float64
}

// newPod creates a new node for the pod in the Dgraph
func newPod(k8sPod api_v1.Pod) Pod {
	pod := Pod{
		Name:      "pod-" + k8sPod.Name,
		IsPod:     true,
		Type:      "pod",
		ID:        dgraph.ID{Xid: k8sPod.Namespace + ":" + k8sPod.Name},
		StartTime: k8sPod.GetCreationTimestamp().Time.Format(time.RFC3339),
	}
	nodeUID, err := createOrGetNodeByID(k8sPod.Spec.NodeName)
	if err == nil {
		pod.Node = &Node{ID: dgraph.ID{UID: nodeUID, Xid: k8sPod.Spec.NodeName}}
	}
	namespaceUID := CreateOrGetNamespaceByID(k8sPod.Namespace)
	if namespaceUID != "" {
		pod.Namespace = &Namespace{ID: dgraph.ID{UID: namespaceUID, Xid: k8sPod.Namespace}}
	}
	pod.Pvcs, pod.StorageRequest = getPodVolumes(k8sPod)
	setPodOwners(&pod, k8sPod)
	return pod
}

// StorePod updates the pod details and create it a new node if not exists.
// It also populates Containers of a pod.
func StorePod(k8sPod api_v1.Pod) error {
	xid := k8sPod.Namespace + ":" + k8sPod.Name
	uid := dgraph.GetUID(xid, IsPod)

	pod := newPod(k8sPod)
	if uid == "" {
		assigned, err := dgraph.MutateNode(pod, dgraph.CREATE)
		if err != nil {
			return err
		}
		log.Infof("Pod with xid: (%s) persisted in dgraph", xid)
		uid = assigned.Uids["blank-0"]
	} else {
		// case when pod is already persisted in dgraph
		// delete end time for pod
		log.Debugf("[SYNC] deleting endTime for pod in dgraph, xid: %s", xid)

		oldPod := Pod{ID: dgraph.ID{UID: uid}, EndTime: ""}
		_, err := dgraph.MutateNode(oldPod, dgraph.DELETE)
		if err != nil {
			return fmt.Errorf("unable to delete endTime for pod, xid: %s, uid: %v, err: %v", xid, uid, err)
		}
		log.Debugf("[SYNC] deleted endTime for pod, xid: %s, uid: %s", xid, uid)
	}

	podDeletedTimestamp := k8sPod.GetDeletionTimestamp()
	if !podDeletedTimestamp.IsZero() {
		pod = Pod{
			ID:      dgraph.ID{Xid: xid, UID: uid},
			EndTime: podDeletedTimestamp.Time.Format(time.RFC3339),
		}
		podData := RetrievePodWithContainers(xid)
		SoftDeleteContainersInTerminatedPod(podData.Containers, podDeletedTimestamp.Time.Format(time.RFC3339))
	} else {
		namespaceUID := CreateOrGetNamespaceByID(k8sPod.Namespace)
		StoreContainersAndMetricsInPod(k8sPod, uid, namespaceUID, &pod)
		populatePodLabels(&pod, k8sPod.Labels)
	}

	// store/update CPUPrice, MemoryPrice
	pod.CPUPrice, pod.MemoryPrice = getPerUnitResourcePriceForNode(k8sPod.Spec.NodeName)

	pod.ID = dgraph.ID{UID: uid, Xid: xid}
	_, err := dgraph.MutateNode(pod, dgraph.UPDATE)
	return err
}

// StorePodsInteraction store the pod interactions in Dgraph
func StorePodsInteraction(sourcePodXID string, destinationPodsXIDs []string, counts []float64) error {
	uid := dgraph.GetUID(sourcePodXID, IsPod)
	if uid == "" {
		log.Println("Source Pod " + sourcePodXID + " is not persisted yet.")
		return fmt.Errorf("source pod: %s is not persisted yet", sourcePodXID)
	}

	pods := retrievePodsWithCountAsEdgeWeightFromPodsXIDs(destinationPodsXIDs, counts)
	source := Pod{
		ID:   dgraph.ID{UID: uid, Xid: sourcePodXID},
		Pods: pods,
	}
	_, err := dgraph.MutateNode(source, dgraph.UPDATE)
	return err
}

func retrievePodsFromPodsXIDs(podsXIDs []string) []*Pod {
	pods := []*Pod{}
	for _, podXID := range podsXIDs {
		podUID := dgraph.GetUID(podXID, IsPod)
		if podUID == "" {
			log.Debugf("Pod uid is empty for pod xid: %s", podXID)
			continue
		}
		pod := &Pod{
			ID: dgraph.ID{UID: podUID, Xid: podXID},
		}
		pods = append(pods, pod)
	}
	return pods
}

func retrievePodsWithCountAsEdgeWeightFromPodsXIDs(podsXIDs []string, counts []float64) []*Pod {
	pods := []*Pod{}
	for index, podXID := range podsXIDs {
		podUID := dgraph.GetUID(podXID, IsPod)
		if podUID == "" {
			log.Printf("Destination pod: %s is not persisted yet", podXID)
			continue
		}

		pod := &Pod{
			ID:    dgraph.ID{UID: podUID, Xid: podXID},
			Count: counts[index],
		}
		pods = append(pods, pod)
	}
	return pods
}

// nolint: gocyclo
func setPodOwners(pod *Pod, k8sPod api_v1.Pod) {
	owners := k8sPod.GetObjectMeta().GetOwnerReferences()
	if owners == nil {
		return
	}
	for _, owner := range owners {
		if owner.Kind == "Deployment" {
			deploymentXID := k8sPod.Namespace + ":" + owner.Name
			deploymentUID := CreateOrGetDeploymentByID(deploymentXID)
			if deploymentUID != "" {
				pod.Deployment = &Deployment{ID: dgraph.ID{UID: deploymentUID, Xid: deploymentXID}}
			}
		} else if owner.Kind == "ReplicaSet" {
			replicasetXID := k8sPod.Namespace + ":" + owner.Name
			replicasetUID := CreateOrGetReplicasetByID(replicasetXID)
			if replicasetUID != "" {
				pod.Replicaset = &Replicaset{ID: dgraph.ID{UID: replicasetUID, Xid: replicasetXID}}
			}
		} else if owner.Kind == "StatefulSet" {
			statefulsetXID := k8sPod.Namespace + ":" + owner.Name
			statefulsetUID := CreateOrGetStatefulsetByID(statefulsetXID)
			if statefulsetUID != "" {
				pod.Statefulset = &Statefulset{ID: dgraph.ID{UID: statefulsetUID, Xid: statefulsetXID}}
			}
		} else if owner.Kind == "Job" {
			jobXID := k8sPod.Namespace + ":" + owner.Name
			jobUID := CreateOrGetJobByID(jobXID)
			if jobUID != "" {
				pod.Job = &Job{ID: dgraph.ID{UID: jobUID, Xid: jobXID}}
			}
		} else if owner.Kind == "DaemonSet" {
			daemonsetXID := k8sPod.Namespace + ":" + owner.Name
			daemonsetUID := CreateOrGetDaemonsetByID(daemonsetXID)
			if daemonsetUID != "" {
				pod.Daemonset = &Daemonset{ID: dgraph.ID{UID: daemonsetUID, Xid: daemonsetXID}}
			}
		} else {
			log.Error("Unknown owner type " + owner.Kind + " for pod.")
		}
	}
}

func getPodVolumes(k8sPod api_v1.Pod) ([]*PersistentVolumeClaim, float64) {
	podVolumes := []*PersistentVolumeClaim{}
	storage := 0.0
	for j := 0; j < len(k8sPod.Spec.Volumes); j++ {
		vol := k8sPod.Spec.Volumes[j]
		if vol.PersistentVolumeClaim != nil {
			pvcXID := k8sPod.Namespace + ":" + vol.PersistentVolumeClaim.ClaimName
			pvcUID := CreateOrGetPersistentVolumeClaimByID(pvcXID)
			if pvcUID != "" {
				podVolumes = append(podVolumes, &PersistentVolumeClaim{ID: dgraph.ID{UID: pvcUID, Xid: pvcXID}})
				pvc, err := getPVCFromUID(pvcUID)
				if err == nil {
					storage += pvc.StorageCapacity
				} else {
					log.Errorf("error while getting pvc from uid: (%v), error: (%v)", pvcUID, err)
				}
			}
		}
	}
	return podVolumes, storage
}

func populatePodLabels(pod *Pod, podLabels map[string]string) {
	log.Debugf("k8s pod: (%v), labels: (%v)", pod.Name, podLabels)
	var labels []*Label
	for key, value := range podLabels {
		labels = append(labels, GetLabel(key, value))
	}
	pod.Labels = labels
}

// RetrievePodWithContainers given a name of pod it retrieves its containers
func RetrievePodWithContainers(xid string) Pod {
	query := `query {
		pods(func: has(isPod)) @filter(eq(xid, "` + xid + `")) {
			name
			containers: ~pod @filter(has(isContainer)) {
				uid
			}
		}
	}`
	type root struct {
		Pods []Pod `json:"pods"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(query, &newRoot)
	if err != nil || len(newRoot.Pods) < 1 {
		log.Errorf("unable to retrieve pod with containers: %v", err)
		return Pod{}
	}
	return newRoot.Pods[0]
}
