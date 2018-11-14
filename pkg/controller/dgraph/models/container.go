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
	"github.com/vmware/purser/pkg/controller/utils"
	api_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Dgraph Model Constants
const (
	IsContainer = "isContainer"
)

// Container schema in dgraph
type Container struct {
	dgraph.ID
	IsContainer   bool       `json:"isContainer,omitempty"`
	Name          string     `json:"name,omitempty"`
	StartTime     string  `json:"startTime,omitempty"`
	EndTime       string  `json:"endTime,omitempty"`
	Pod           Pod        `json:"pod,omitempty"`
	Procs         []*Proc    `json:"procs,omitempty"`
	Namespace     *Namespace `json:"namespace,omitempty"`
	CPURequest    float64    `json:"cpuRequest,omitempty"`
	CPULimit      float64    `json:"cpuLimit,omitempty"`
	MemoryRequest float64    `json:"memoryRequest,omitempty"`
	MemoryLimit   float64    `json:"memoryLimit,omitempty"`
	Type          string     `json:"type,omitempty"`
}

func newContainer(container api_v1.Container, podUID, namespaceUID string, pod api_v1.Pod) (*api.Assigned, error) {
	containerXid := pod.Namespace + ":" + pod.Name + ":" + container.Name
	requests := container.Resources.Requests
	limits := container.Resources.Limits
	c := &Container{
		ID:            dgraph.ID{Xid: containerXid},
		Name:          "container-" + container.Name,
		IsContainer:   true,
		Type:          "container",
		StartTime:     pod.GetCreationTimestamp().Time.Format(time.RFC3339),
		Pod:           Pod{ID: dgraph.ID{UID: podUID, Xid: pod.Namespace + ":" + pod.Name}},
		CPURequest:    utils.ConvertToFloat64CPU(requests.Cpu()),
		CPULimit:      utils.ConvertToFloat64CPU(limits.Cpu()),
		MemoryRequest: utils.ConvertToFloat64GB(requests.Memory()),
		MemoryLimit:   utils.ConvertToFloat64GB(limits.Memory()),
	}
	if namespaceUID != "" {
		c.Namespace = &Namespace{ID: dgraph.ID{UID: namespaceUID, Xid: pod.Namespace}}
	}
	return dgraph.MutateNode(c, dgraph.CREATE)
}

// StoreAndRetrieveContainersAndMetrics fetchs the list of containers in given pod
// Create a new container in dgraph if container is not in it.
func StoreAndRetrieveContainersAndMetrics(pod api_v1.Pod, podUID, namespaceUID string) ([]*Container, Metrics) {
	containers := []*Container{}
	cpuRequest := &resource.Quantity{}
	memoryRequest := &resource.Quantity{}
	cpuLimit := &resource.Quantity{}
	memoryLimit := &resource.Quantity{}
	for _, c := range pod.Spec.Containers {
		container, err := storeContainerIfNotExist(c, pod, podUID, namespaceUID)
		if err == nil {
			containers = append(containers, container)
			requests := c.Resources.Requests
			limits := c.Resources.Limits
			utils.AddResourceAToResourceB(requests.Cpu(), cpuRequest)
			utils.AddResourceAToResourceB(requests.Memory(), memoryRequest)
			utils.AddResourceAToResourceB(limits.Cpu(), cpuLimit)
			utils.AddResourceAToResourceB(limits.Memory(), memoryLimit)
		}
	}
	return containers, Metrics{
		CPURequest:    utils.ConvertToFloat64CPU(cpuRequest),
		CPULimit:      utils.ConvertToFloat64CPU(cpuLimit),
		MemoryRequest: utils.ConvertToFloat64GB(memoryRequest),
		MemoryLimit:   utils.ConvertToFloat64GB(memoryLimit),
	}
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

func storeContainerIfNotExist(c api_v1.Container, pod api_v1.Pod, podUID, namespaceUID string) (*Container, error) {
	podXid := pod.Namespace + ":" + pod.Name
	containerXid := podXid + ":" + c.Name
	containerUID := dgraph.GetUID(containerXid, IsContainer)

	var container *Container
	if containerUID == "" {
		assigned, err := newContainer(c, podUID, namespaceUID, pod)
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
		container.EndTime = endTime.Format(time.RFC3339)
	}
	_, err := dgraph.MutateNode(containers, dgraph.UPDATE)
	if err != nil {
		log.Error(err)
	}
	deleteProcessesInTerminatedContainers(containers)
}

// RetrieveAllContainers ...
func RetrieveAllContainers() ([]byte, error) {
	const q = `query {
		container(func: has(isContainer)) {
			name
			type
			process: ~container @filter(has(isProc)) {
				name
				type
			}
		}
	}`

	result, err := dgraph.ExecuteQueryRaw(q)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RetrieveContainer ...
func RetrieveContainer(name string) ([]byte, error) {
	q := `query {
		container(func: has(isContainer)) @filter(eq(name, "` + name + `")) {
			name
			type
			process: ~container @filter(has(isProc)) {
				name
				type
			}
		}
	}`


	result, err := dgraph.ExecuteQueryRaw(q)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RetrieveContainerWithMetrics ...
func RetrieveContainerWithMetrics(name string) (JsonDataWrapper, error) {
	q := `query {
		parent(func: has(isContainer)) @filter(eq(name, "` + name + `")) {
			name
			type
			cpu: cpu as cpuRequest
			memory: memory as memoryRequest
			cpuCost: math(cpu * ` + defaultCPUCostPerCPUPerHour + `)
			memoryCost: math(memory * ` + defaultMemCostPerGBPerHour + `)
		}
	}`
	parentRoot := ParentWrapper{}
	err := dgraph.ExecuteQuery(q, &parentRoot)
	root := JsonDataWrapper{}
	if len(parentRoot.Parent) == 0 {
		return root, err
	}
	root.Data = ParentWrapper{
		Name: parentRoot.Parent[0].Name,
		Type: parentRoot.Parent[0].Type,
		Children: parentRoot.Parent[0].Children,
		CPU: parentRoot.Parent[0].CPU,
		Memory: parentRoot.Parent[0].Memory,
		CPUCost: parentRoot.Parent[0].CPUCost,
		MemoryCost: parentRoot.Parent[0].MemoryCost,
	}
	return root, err
}

