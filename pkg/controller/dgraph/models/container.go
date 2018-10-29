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
	IsContainer = "isContainer"
)

// Container schema in dgraph
type Container struct {
	dgraph.ID
	IsContainer bool       `json:"isContainer,omitempty"`
	Name        string     `json:"name,omitempty"`
	StartTime   time.Time  `json:"startTime,omitempty"`
	EndTime     time.Time  `json:"endTime,omitempty"`
	Pod         Pod        `json:"pod,omitempty"`
	Procs       []*Proc    `json:"procs,omitempty"`
	Namespace   *Namespace `json:"namespace,omitempty"`
}

func newContainer(containerXid, containerName, podUID string, pod api_v1.Pod) (*api.Assigned, error) {
	container := &Container{
		ID:          dgraph.ID{Xid: containerXid},
		Name:        containerName,
		IsContainer: true,
		StartTime:   pod.GetCreationTimestamp().Time,
		Pod:         Pod{ID: dgraph.ID{UID: podUID, Xid: pod.Namespace + ":" + pod.Name}},
	}
	namespaceUID, err := createOrGetNamespaceByID(pod.Namespace)
	if err == nil {
		container.Namespace = &Namespace{ID: dgraph.ID{UID: namespaceUID, Xid: pod.Namespace}}
	}
	return dgraph.MutateNode(container, dgraph.CREATE)
}

// StoreAndRetrieveContainers fetchs the list of containers in given pod
// Create a new container in dgraph if container is not in it.
func StoreAndRetrieveContainers(pod api_v1.Pod, podUID string) []*Container {
	containers := []*Container{}
	for _, c := range pod.Spec.Containers {
		container, err := storeContainerIfNotExist(c, pod, podUID)
		if err == nil {
			containers = append(containers, container)
		}
	}
	return containers
}

// StoreContainerProcessEdge ...
func StoreContainerProcessEdge(containerXID string, procsXIDs []string) error {
	containerUID := dgraph.GetUID(containerXID, IsContainer)
	if containerUID == "" {
		return fmt.Errorf("Container: %s not persisted in dgraph", containerXID)
	}

	procs := retrieveProcessesFromProcessesXIDs(procsXIDs)
	container := Container{
		ID:    dgraph.ID{UID: containerUID, Xid: containerXID},
		Procs: procs,
	}
	_, err := dgraph.MutateNode(container, dgraph.UPDATE)
	return err
}

func storeContainerIfNotExist(c api_v1.Container, pod api_v1.Pod, podUID string) (*Container, error) {
	podXid := pod.Namespace + ":" + pod.Name
	containerXid := podXid + ":" + c.Name
	containerUID := dgraph.GetUID(containerXid, IsContainer)

	var container *Container
	if containerUID == "" {
		assigned, err := newContainer(containerXid, c.Name, podUID, pod)
		if err != nil {
			log.Errorf("Unable to create container: %s", containerXid)
			return container, err
		}
		log.Infof("Container with xid: (%s) persisted in dgraph", containerXid)
		containerUID = assigned.Uids["blank-0"]
	}

	container = &Container{
		ID: dgraph.ID{UID: containerUID, Xid: containerXid},
	}
	return container, nil
}

func deleteContainersInTerminatedPod(containers []*Container, endTime time.Time) {
	for _, container := range containers {
		container.EndTime = endTime
	}
	_, err := dgraph.MutateNode(containers, dgraph.UPDATE)
	if err != nil {
		log.Error(err)
	}
	deleteProcessesInTerminatedContainers(containers)
}
