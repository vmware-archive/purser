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
	"time"

	"github.com/vmware/purser/pkg/controller/dgraph"
	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	log "github.com/Sirupsen/logrus"
)

// Dgraph Model Constants
const (
	IsDeployment = "isDeployment"
)

// Deployment schema in dgraph
type Deployment struct {
	dgraph.ID
	IsDeployment bool       `json:"isDeployment,omitempty"`
	Name         string     `json:"name,omitempty"`
	StartTime    string  `json:"startTime,omitempty"`
	EndTime      string  `json:"endTime,omitempty"`
	Namespace    *Namespace `json:"namespace,omitempty"`
	Pods         []*Pod     `json:"pod,omitempty"`
	Type         string     `json:"type,omitempty"`
	Replicasets []*Replicaset `json:"replicaset,omitempty"`
}

func createDeploymentObject(deployment apps_v1beta1.Deployment) Deployment {
	newDeployment := Deployment{
		Name:         deployment.Name,
		IsDeployment: true,
		Type:         "deployment",
		ID:           dgraph.ID{Xid: deployment.Namespace + ":" + deployment.Name},
		StartTime:    deployment.GetCreationTimestamp().Time.Format(time.RFC3339),
	}
	namespaceUID := CreateOrGetNamespaceByID(deployment.Namespace)
	if namespaceUID != "" {
		newDeployment.Namespace = &Namespace{ID: dgraph.ID{UID: namespaceUID, Xid: deployment.Namespace}}
	}
	deploymentDeletionTimestamp := deployment.GetDeletionTimestamp()
	if !deploymentDeletionTimestamp.IsZero() {
		newDeployment.EndTime = deploymentDeletionTimestamp.Time.Format(time.RFC3339)
	}
	return newDeployment
}

// StoreDeployment create a new deployment in the Dgraph and updates if already present.
func StoreDeployment(deployment apps_v1beta1.Deployment) (string, error) {
	xid := deployment.Namespace + ":" + deployment.Name
	uid := dgraph.GetUID(xid, IsDeployment)

	newDeployment := createDeploymentObject(deployment)
	if uid != "" {
		newDeployment.UID = uid
	}
	assigned, err := dgraph.MutateNode(newDeployment, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	return assigned.Uids["blank-0"], nil
}

// CreateOrGetDeploymentByID returns the uid of namespace if exists,
// otherwise creates the deployment and returns uid.
func CreateOrGetDeploymentByID(xid string) string {
	if xid == "" {
		return ""
	}
	uid := dgraph.GetUID(xid, IsDeployment)

	if uid != "" {
		return uid
	}

	d := Deployment{
		ID:           dgraph.ID{Xid: xid},
		Name:         xid,
		IsDeployment: true,
	}
	assigned, err := dgraph.MutateNode(d, dgraph.CREATE)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return assigned.Uids["blank-0"]
}

// RetrieveAllDeployments ...
func RetrieveAllDeployments() ([]byte, error) {
	const q = `query {
		deployment(func: has(isDeployment)) {
			name
			type
			replicaset: ~deployment @filter(has(isReplicaset)) {
				name
				type
				pod: ~replicaset @filter(has(isPod)) {
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

// RetrieveDeployment ...
func RetrieveDeployment(name string) ([]byte, error) {
	q := `query {
		deployment(func: has(isDeployment)) @filter(eq(name, "` + name + `")) {
			name
			type
			replicaset: ~deployment @filter(has(isReplicaset)) {
				name
				type
				pod: ~replicaset @filter(has(isPod)) {
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

// RetrieveDeploymentWithMetrics ...
func RetrieveDeploymentWithMetrics(name string) (JsonDataWrapper, error) {
	q := `query {
		dep as var(func: has(isDeployment)) @filter(eq(name, "` + name + `")) {
			~deployment @filter(has(isReplicaset)) {
				~replicaset @filter(has(isPod)) {
					replicasetPodCpu as cpuRequest
					replicasetPodMemory as memoryRequest
				}
				deploymentReplicasetCpu as sum(val(replicasetPodCpu))
				deploymentReplicasetMemory as sum(val(replicasetPodMemory))
			}
			deploymentCpu as sum(val(deploymentReplicasetCpu))
			deploymentMemory as sum(val(deploymentReplicasetMemory))
		}

		parent(func: uid(dep)) {
			name
			type
			children: ~deployment @filter(has(isReplicaset)) {
				name
				type
				cpu: val(deploymentReplicasetCpu)
				memory: val(deploymentReplicasetMemory)
			}
			cpu: val(deploymentCpu)
			memory: val(deploymentMemory)
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