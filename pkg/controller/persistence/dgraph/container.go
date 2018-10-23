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
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/vmware/purser/pkg/controller/utils"
	api_v1 "k8s.io/api/core/v1"
)

// Dgraph Model Constants
const (
	IsContainer = "isContainer"
)

// Container schema in dgraph
type Container struct {
	ID
	IsContainer bool      `json:"isContainer,omitempty"`
	Name        string    `json:"name,omitempty"`
	StartTime   time.Time `json:"startTime,omitempty"`
	EndTime     time.Time `json:"endTime,omitempty"`
	Pod         Pod       `json:"pod,omitempty"`
	Procs       []*Proc   `json:"procs,omitempty"`
}

func createContainer(containerXid, containerName string) (*api.Assigned, error) {
	container := &Container{
		ID:          ID{Xid: containerXid},
		Name:        containerName,
		IsContainer: true,
	}
	bytes := utils.JsonMarshal(container)
	return MutateNode(Client, bytes)
}

// GetContainers fetchs the list of containers in given pod
func GetContainers(pod api_v1.Pod, podUID string) []*Container {
	podXid := pod.Namespace + ":" + pod.Name

	containers := []*Container{}
	for _, c := range pod.Spec.Containers {
		containerXid := podXid + ":" + c.Name
		containerUID := GetUID(Client, containerXid, IsContainer)

		var container *Container
		if containerUID == "" {
			assigned, err := createContainer(containerXid, c.Name)
			containerUID = assigned.Uids["blank-0"]
			if err != nil {
				log.Errorf("Unable to create container: %s", containerXid)
				continue
			}
		}
		container = &Container{
			ID:  ID{UID: containerUID},
			Pod: Pod{ID: ID{UID: podUID}},
		}
		bytes := utils.JsonMarshal(container)
		MutateNode(Client, bytes)

		containers = append(containers, container)
	}
	return containers
}

// CreateEdgeFromContainerToProcs ...
func CreateEdgeFromContainerToProcs(containerXID string, procsXIDs []string) error {
	containerUID := GetUID(Client, containerXID, IsContainer)
	if containerUID == "" {
		return fmt.Errorf("Container: %s not persisted in dgraph", containerXID)
	}

	procs := []*Proc{}
	for _, procXID := range procsXIDs {
		procUID := GetUID(Client, procXID, IsProc)
		if procUID != "" {
			procs = append(procs, &Proc{ID: ID{UID: procUID}})
		}
	}
	container := Container{
		ID:    ID{UID: containerUID},
		Procs: procs,
	}
	bytes := utils.JsonMarshal(container)
	_, err := MutateNode(Client, bytes)
	return err
}
