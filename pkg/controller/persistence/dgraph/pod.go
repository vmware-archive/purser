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

package dgraph

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/vmware/purser/pkg/controller/utils"

	api_v1 "k8s.io/api/core/v1"
)

// Dgraph Model Constants
const (
	IsPod = "isPod"
)

// ID maps the external ID used in Dgraph to the UID
type ID struct {
	Xid string `json:"xid,omitempty"`
	UID string `json:"uid,omitempty"`
}

// Pod schema in dgraph
type Pod struct {
	ID
	IsPod      bool         `json:"isPod,omitempty"`
	Name       string       `json:"name,omitempty"`
	StartTime  time.Time    `json:"startTime,omitempty"`
	EndTime    time.Time    `json:"endTime,omitempty"`
	Containers []*Container `json:"containers,omitempty"`
	Interacts  []*Pod       `json:"interacts,omitempty"`
	Count      float64      `json:"interacts|count,omitempty"`
}

// createPod creates a new node for the pod in the Dgraph
func createPod(pod api_v1.Pod, xid string) (*api.Assigned, error) {
	newPod := Pod{
		Name:  pod.Name,
		IsPod: true,
		ID:    ID{Xid: xid},
	}
	bytes := utils.JsonMarshal(newPod)
	return MutateNode(Client, bytes)
}

// PersistPod updates the pod details and create it a new node if not exists.
// It also populates Containers of a pod.
func PersistPod(pod api_v1.Pod) error {
	xid := pod.Namespace + ":" + pod.Name
	uid := GetUID(Client, xid, IsPod)

	var newPod Pod
	if uid == "" {
		assigned, err := createPod(pod, xid)
		if err != nil {
			return err
		}
		uid = assigned.Uids["blank-0"]
	}

	newPod = Pod{
		ID:         ID{Xid: xid, UID: uid},
		Containers: GetContainers(pod, uid),
	}
	bytes, err := json.Marshal(newPod)
	if err != nil {
		return err
	}
	_, err = MutateNode(Client, bytes)
	return err
}

// PersistPodsInteractionGraph store the pod interactions in Dgraph
func PersistPodsInteractionGraph(sourcePodXID string, destinationPodsXIDs []string, counts []float64) error {
	uid := GetUID(Client, sourcePodXID, IsPod)
	if uid == "" {
		log.Println("Source Pod " + sourcePodXID + " is not persisted yet.")
		return fmt.Errorf("source pod: %s is not persisted yet", sourcePodXID)
	}

	pods := []*Pod{}
	for index, destinationPodXID := range destinationPodsXIDs {
		dstUID := GetUID(Client, destinationPodXID, IsPod)
		if dstUID == "" {
			log.Printf("Destination pod: %s is not persistet yet", destinationPodXID)
			continue
		}

		pod := &Pod{
			ID:    ID{UID: dstUID},
			Count: counts[index],
		}
		pods = append(pods, pod)
	}
	source := Pod{
		ID:        ID{UID: uid},
		Interacts: pods,
	}

	bytes := utils.JsonMarshal(source)
	_, err := MutateNode(Client, bytes)
	return err
}

// FetchAllPods returns all pods in the dgraph
func FetchAllPods() ([]Pod, error) {
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
	err := executeQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Pods, nil
}
