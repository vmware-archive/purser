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
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/utils"

	api_v1 "k8s.io/api/core/v1"
)

// Dgraph Model Constants
const (
	IsPod = "isPod"
)

// Pod schema in dgraph
type Pod struct {
	dgraph.ID
	IsPod      bool         `json:"isPod,omitempty"`
	Name       string       `json:"name,omitempty"`
	StartTime  time.Time    `json:"startTime,omitempty"`
	EndTime    time.Time    `json:"endTime,omitempty"`
	Containers []*Container `json:"containers,omitempty"`
	Interacts  []*Pod       `json:"interacts,omitempty"`
	Count      float64      `json:"interacts|count,omitempty"`
}

// newPod creates a new node for the pod in the Dgraph
func newPod(k8sPod api_v1.Pod, xid string) (*api.Assigned, error) {
	pod := Pod{
		Name:      k8sPod.Name,
		IsPod:     true,
		ID:        dgraph.ID{Xid: xid},
		StartTime: k8sPod.GetCreationTimestamp().Time,
	}
	bytes := utils.JSONMarshal(pod)
	return dgraph.MutateNode(bytes)
}

// StorePod updates the pod details and create it a new node if not exists.
// It also populates Containers of a pod.
func StorePod(k8sPod api_v1.Pod) error {
	xid := k8sPod.Namespace + ":" + k8sPod.Name
	uid := dgraph.GetUID(xid, IsPod)

	var pod Pod
	if uid == "" {
		assigned, err := newPod(k8sPod, xid)
		if err != nil {
			return err
		}
		uid = assigned.Uids["blank-0"]
	}

	podDeletedTimestamp := k8sPod.GetDeletionTimestamp()
	isDeleted := !podDeletedTimestamp.IsZero()
	if isDeleted {
		pod = Pod{
			ID:         dgraph.ID{Xid: xid, UID: uid},
			Containers: StoreAndRetrieveContainers(k8sPod, uid, isDeleted),
			EndTime:    podDeletedTimestamp.Time,
		}
	} else {
		pod = Pod{
			ID:         dgraph.ID{Xid: xid, UID: uid},
			Containers: StoreAndRetrieveContainers(k8sPod, uid, isDeleted),
		}
	}

	bytes, err := json.Marshal(pod)
	if err != nil {
		return err
	}
	_, err = dgraph.MutateNode(bytes)
	return err
}

// StorePodsInteraction store the pod interactions in Dgraph
func StorePodsInteraction(sourcePodXID string, destinationPodsXIDs []string, counts []float64) error {
	uid := dgraph.GetUID(sourcePodXID, IsPod)
	if uid == "" {
		log.Println("Source Pod " + sourcePodXID + " is not persisted yet.")
		return fmt.Errorf("source pod: %s is not persisted yet", sourcePodXID)
	}

	pods := []*Pod{}
	for index, destinationPodXID := range destinationPodsXIDs {
		dstUID := dgraph.GetUID(destinationPodXID, IsPod)
		if dstUID == "" {
			log.Printf("Destination pod: %s is not persistet yet", destinationPodXID)
			continue
		}

		pod := &Pod{
			ID:    dgraph.ID{UID: dstUID, Xid: destinationPodXID},
			Count: counts[index],
		}
		pods = append(pods, pod)
	}
	source := Pod{
		ID:        dgraph.ID{UID: uid, Xid: sourcePodXID},
		Interacts: pods,
	}

	bytes := utils.JSONMarshal(source)
	_, err := dgraph.MutateNode(bytes)
	return err
}

// RetrieveAllPods returns all pods in the dgraph
func RetrieveAllPods() ([]Pod, error) {
	const q = `query {
		pods(func: has(isPod)) {
			name
			interacts @facets {
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
