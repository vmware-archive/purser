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

	"github.com/dgraph-io/dgo/protos/api"
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
	IsPod         bool         `json:"isPod,omitempty"`
	Name          string       `json:"name,omitempty"`
	StartTime     string    `json:"startTime,omitempty"`
	EndTime       string    `json:"endTime,omitempty"`
	Containers    []*Container `json:"containers,omitempty"`
	Interacts     []*Pod       `json:"interacts,omitempty"`
	Count         float64      `json:"interacts|count,omitempty"`
	Node          *Node        `json:"podNode,omitempty"`
	Namespace     *Namespace   `json:"namespace,omitempty"`
	Deployment    *Deployment  `json:"deployment,omitempty"`
	Replicaset    *Replicaset  `json:"replicaset,omitempty"`
	Statefulset   *Statefulset `json:"statefulset,omitempty"`
	Daemonset   *Daemonset `json:"daemonset,omitempty"`
	Job   		*Job 			`json:"job,omitempty"`
	CPURequest    float64      `json:"cpuRequest,omitempty"`
	CPULimit      float64      `json:"cpuLimit,omitempty"`
	MemoryRequest float64      `json:"memoryRequest,omitempty"`
	MemoryLimit   float64      `json:"memoryLimit,omitempty"`
	Type          string       `json:"type,omitempty"`
	Cid []Service `json:"cid,omitempty"`
}

// Metrics ...
type Metrics struct {
	CPURequest    float64
	CPULimit      float64
	MemoryRequest float64
	MemoryLimit   float64
}

// newPod creates a new node for the pod in the Dgraph
func newPod(k8sPod api_v1.Pod) (*api.Assigned, error) {
	pod := Pod{
		Name:      k8sPod.Name,
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
	setPodOwners(&pod, k8sPod)
	return dgraph.MutateNode(pod, dgraph.CREATE)
}

// StorePod updates the pod details and create it a new node if not exists.
// It also populates Containers of a pod.
func StorePod(k8sPod api_v1.Pod) error {
	xid := k8sPod.Namespace + ":" + k8sPod.Name
	uid := dgraph.GetUID(xid, IsPod)

	var pod Pod
	if uid == "" {
		assigned, err := newPod(k8sPod)
		if err != nil {
			return err
		}
		log.Infof("Pod with xid: (%s) persisted in dgraph", xid)
		uid = assigned.Uids["blank-0"]
	}

	podDeletedTimestamp := k8sPod.GetDeletionTimestamp()
	if !podDeletedTimestamp.IsZero() {
		pod = Pod{
			ID:      dgraph.ID{Xid: xid, UID: uid},
			EndTime: podDeletedTimestamp.Time.Format(time.RFC3339),
		}
		deleteContainersInTerminatedPod(pod.Containers, podDeletedTimestamp.Time)
	} else {
		namespaceUID := CreateOrGetNamespaceByID(k8sPod.Namespace)
		containers, metrics := StoreAndRetrieveContainersAndMetrics(k8sPod, uid, namespaceUID)
		pod = Pod{
			ID:            dgraph.ID{Xid: xid, UID: uid},
			Containers:    containers,
			CPURequest:    metrics.CPURequest,
			CPULimit:      metrics.CPULimit,
			MemoryRequest: metrics.MemoryRequest,
			MemoryLimit:   metrics.MemoryLimit,
		}
	}

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
		ID:        dgraph.ID{UID: uid, Xid: sourcePodXID},
		Interacts: pods,
	}
	_, err := dgraph.MutateNode(source, dgraph.UPDATE)
	return err
}

// RetrievePodsInteractionsForAllPodsOrphanedTrue returns all pods in the dgraph
func RetrievePodsInteractionsForAllPodsWithCount() ([]Pod, error) {
	const q = `query {
		pods(func: has(isPod)) {
			name
			interacts {
				name
				count
			}
			cid: ~pod @filter(has(isService)) {
				name
			}
		}
	}`

	type root struct {
		Pods []Pod `json:"pods"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Pods, nil
}

// RetrievePodsInteractionsForAllPodsOrphanedTrue returns all pods in the dgraph
func RetrievePodsInteractionsForAllPodsOrphanedTrue() ([]Pod, error) {
	const q = `query {
		pods(func: has(isPod)) {
			name
			interacts {
				name
			}
		}
	}`

	type root struct {
		Pods []Pod `json:"pods"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Pods, nil
}

func RetrievePodsInboundInteractionsForAllPods() ([]byte, error) {
	const q = `query {
		pods(func: has(isPod)) {
			name
			inbound: ~pod @filter(has(isPod)) {
				name
			}
		}
	}`

	result, err := dgraph.ExecuteQueryRaw(q)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RetrievePodsInteractionsForAllPodsOrphanedFalse returns all pods in the dgraph which has edge interacts
func RetrievePodsInteractionsForAllPodsOrphanedFalse() ([]Pod, error) {
	const q = `query {
		pods(func: has(isPod)) @filter(has(interacts)) {
			name
			interacts {
				name
			}
		}
	}`

	type root struct {
		Pods []Pod `json:"pods"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Pods, nil
}

// RetrievePodsInteractionsForGivenPod ...
func RetrievePodsInteractionsForGivenPod(name string) ([]Pod, error) {
	q := `query {
		pods(func: has(isPod)) @filter(eq(name, "` + name + `")) {
			name
			interacts {
				name
			}
		}
	}`

	type root struct {
		Pods []Pod `json:"pods"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Pods, nil
}

func RetrievePodsInboundInteractionsForGivenPod(name string) ([]byte, error) {
	q := `query {
		pods(func: has(isPod)) @filter(eq(name, "` + name + `")) {
			name
			inbound: ~pod @filter(has(isPod)) {
				name
			}
		}
	}`

	result, err := dgraph.ExecuteQueryRaw(q)
	if err != nil {
		return nil, err
	}
	return result, nil
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

// RetrieveAllPods ...
func RetrieveAllPods() ([]byte, error) {
	const q = `query {
		pod(func: has(isPod)) {
			name
			type
			container: ~pod @filter(has(isContainer)) {
				name
				type
				process: ~container @filter(has(isProc)) {
					name
					type
				}
			}
		}
	}`

	result, err := dgraph.ExecuteQueryRaw(q)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RetrievePod ...
func RetrievePod(name string) ([]byte, error) {
	q := `query {
		pod(func: has(isPod)) @filter(eq(name, "` + name + `")) {
			name
			type
			container: ~pod @filter(has(isContainer)) {
				name
				type
				process: ~container @filter(has(isProc)) {
					name
					type
				}
			}
		}
	}`


	result, err := dgraph.ExecuteQueryRaw(q)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RetrieveReplicasetsWithMetrics ...
func RetrievePodWithMetrics(name string) (JsonDataWrapper, error) {
	q := `query {
		parent(func: has(isPod)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~pod @filter(has(isContainer)) {
				name
				type
				cpu: cpuRequest
				memory: memoryRequest
			}
			cpu: cpuRequest
			memory: memoryRequest
		}
	}`
	parentRoot := ParentWrapper{}
	err := dgraph.ExecuteQuery(q, &parentRoot)
	root := JsonDataWrapper{}
	root.Data = ParentWrapper{
		Name: parentRoot.Parent[0].Name,
		Type: parentRoot.Parent[0].Type,
		Children: parentRoot.Parent[0].Children,
		CPU: parentRoot.Parent[0].CPU,
		Memory: parentRoot.Parent[0].Memory,
	}
	return root, err
}