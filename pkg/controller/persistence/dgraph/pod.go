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
	"time"

	api_v1 "k8s.io/api/core/v1"
	"github.com/dgraph-io/dgo/protos/api"
	"log"
)

const (
	IsPod       = "isPod"
	IsContainer = "isContainer"
)

type ID struct {
	Xid string `json:"xid,omitempty"`
	UID string `json:"uid,omitempty"`
}

type Pod struct {
	ID
	IsPod      bool         `json:"isPod,omitempty"`
	Name       string       `json:"name,omitempty"`
	StartTime  time.Time    `json:"startTime,omitempty"`
	EndTime    time.Time    `json:"endTime,omitempty"`
	Containers []*Container `json:"containers,omitempty"`
	Interacts  []*Pod       `json:interacts,omitempty`
}

type Container struct {
	ID
	IsContainer bool      `json:"isContainer,omitempty"`
	Name        string    `json:"name,omitempty"`
	StartTime   time.Time `json:"startTime,omitempty"`
	EndTime     time.Time `json:"endTime,omitempty"`
	Pod         Pod       `json:"pod,omitempty"`
}

func CreatePod(pod api_v1.Pod, xid string) (*api.Assigned, error) {
	newPod := Pod{
		Name:  pod.Name,
		IsPod: true,
		ID:    ID{Xid: xid},
	}
	bytes, err := json.Marshal(newPod)
	if err != nil {
		return nil, err
	}
	return MutateNode(Client, bytes)
}

func PersistPod(pod api_v1.Pod) error {
	xid := pod.Namespace + ":" + pod.Name
	uid, _ := GetUId(Client, xid, IsPod)

	var newPod Pod
	if uid == "" {
		assigned, err := CreatePod(pod, xid)
		if err != nil {
			return err
		}
		uid = assigned.Uids["blank-0"]
		newPod = Pod{
			ID:         ID{Xid: xid, UID: uid},
			Containers: GetContainers(pod, uid),
		}
	} else {
		newPod = Pod{
			ID:         ID{Xid: xid, UID: uid},
			Containers: GetContainers(pod, uid),
		}
	}
	bytes, err := json.Marshal(newPod)
	if err != nil {
		return err
	}
	_, err = MutateNode(Client, bytes)
	return err
}

func GetContainers(pod api_v1.Pod, podUid string) []*Container {
	podXid := pod.Namespace + ":" + pod.Name

	containers := []*Container{}
	for _, c := range pod.Spec.Containers {
		containerXid := podXid + ":" + c.Name
		uid, _ := GetUId(Client, containerXid, IsContainer)

		var container *Container
		if uid == "" {
			container = &Container{
				ID:          ID{Xid: containerXid, UID: uid},
				Name:        c.Name,
				IsContainer: true,
				Pod:         Pod{ID: ID{Xid: podXid, UID: podUid}},
			}
		} else {
			container = &Container{
				ID:  ID{Xid: containerXid, UID: uid},
				Pod: Pod{ID: ID{Xid: podXid, UID: podUid}},
			}
		}

		containers = append(containers, container)
	}
	return containers
}

func PersistPodsInteractionGraph(sourcePod string, destinationPods []string) error {
	uid, err := GetUId(Client, sourcePod, IsPod)
	if err != nil {
		return err
	}
	if uid == "" {
		log.Println("Source Pod " + sourcePod + " is not persisted yet.")
	}

	pods := []*Pod{}
	for _, destinationPod := range destinationPods {
		uid, err := GetUId(Client, destinationPod, IsPod)
		if err != nil {
			return err
		}

		if uid == "" {
			continue
		}

		pod := &Pod{
			ID: ID{UID: uid},
		}
		pods = append(pods, pod)
	}
	source := Pod{
		ID: ID{UID: uid},
		Interacts:	pods,
	}

	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	_, err = MutateNode(Client, bytes)
	return err

}
