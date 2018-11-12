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
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	log "github.com/Sirupsen/logrus"
)

// Dgraph Model Constants
const (
	IsReplicaset = "isReplicaset"
)

// Replicaset schema in dgraph
type Replicaset struct {
	dgraph.ID
	IsReplicaset bool        `json:"isReplicaset,omitempty"`
	Name         string      `json:"name,omitempty"`
	StartTime    string   `json:"startTime,omitempty"`
	EndTime      string   `json:"endTime,omitempty"`
	Namespace    *Namespace  `json:"namespace,omitempty"`
	Deployment   *Deployment `json:"deployment,omitempty"`
	Pods         []*Pod      `json:"pod,omitempty"`
	Type         string      `json:"type,omitempty"`
	CPU    float64    `json:"cpu,omitempty"`
	Memory float64    `json:"memory,omitempty"`
}

// ReplicasetsWithMetrics ...
type ReplicasetsWithMetrics struct {
	Replicaset []Replicaset  `json:"replicaset,omitempty"`
	CPU    float64    `json:"cpu,omitempty"`
	Memory float64    `json:"memory,omitempty"`
}

func createReplicasetObject(replicaset ext_v1beta1.ReplicaSet) Replicaset {
	newReplicaset := Replicaset{
		Name:         replicaset.Name,
		IsReplicaset: true,
		Type:         "replicaset",
		ID:           dgraph.ID{Xid: replicaset.Namespace + ":" + replicaset.Name},
		StartTime:    replicaset.GetCreationTimestamp().Time.Format(time.RFC3339),
	}
	namespaceUID := CreateOrGetNamespaceByID(replicaset.Namespace)
	if namespaceUID != "" {
		newReplicaset.Namespace = &Namespace{ID: dgraph.ID{UID: namespaceUID, Xid: replicaset.Namespace}}
	}
	replicasetDeletionTimestamp := replicaset.GetDeletionTimestamp()
	if !replicasetDeletionTimestamp.IsZero() {
		newReplicaset.EndTime = replicasetDeletionTimestamp.Time.Format(time.RFC3339)
	}
	setReplicasetOwners(&newReplicaset, replicaset)
	return newReplicaset
}

// StoreReplicaset create a new replicaset in the Dgraph and updates if already present.
func StoreReplicaset(replicaset ext_v1beta1.ReplicaSet) (string, error) {
	xid := replicaset.Namespace + ":" + replicaset.Name
	uid := dgraph.GetUID(xid, IsReplicaset)

	newReplicaset := createReplicasetObject(replicaset)
	if uid != "" {
		newReplicaset.UID = uid
	}
	assigned, err := dgraph.MutateNode(newReplicaset, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	return assigned.Uids["blank-0"], nil
}

func setReplicasetOwners(r *Replicaset, replicaset ext_v1beta1.ReplicaSet) {
	owners := replicaset.GetObjectMeta().GetOwnerReferences()
	if owners == nil {
		return
	}
	for _, owner := range owners {
		if owner.Kind == "Deployment" {
			deploymentXID := replicaset.Namespace + ":" + owner.Name
			deploymentUID := CreateOrGetDeploymentByID(deploymentXID)
			if deploymentUID != "" {
				r.Deployment = &Deployment{ID: dgraph.ID{UID: deploymentUID, Xid: deploymentXID}}
			}
		} else {
			log.Error("Unknown owner type " + owner.Kind + " for replicaset.")
		}
	}
}

// CreateOrGetReplicasetByID returns the uid of namespace if exists,
// otherwise creates the replicaset and returns uid.
func CreateOrGetReplicasetByID(xid string) string {
	if xid == "" {
		return ""
	}
	uid := dgraph.GetUID(xid, IsReplicaset)

	if uid != "" {
		return uid
	}

	d := Replicaset{
		ID:           dgraph.ID{Xid: xid},
		Name:         xid,
		IsReplicaset: true,
	}
	assigned, err := dgraph.MutateNode(d, dgraph.CREATE)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return assigned.Uids["blank-0"]
}

// RetrieveAllReplicasets ...
func RetrieveAllReplicasets() ([]byte, error) {
	const q = `query {
		replicaset(func: has(isReplicaset)) {
			name
			type
			pod: "~replicaset @filter(has(isPod)) {
				name
				type
				container: ~pod @filter(has(isContainer)) {
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

// RetrieveReplicaset ...
func RetrieveReplicaset(name string) ([]byte, error) {
	q := `query {
		replicaset(func: has(isReplicaset)) @filter(eq(name, "` + name + `")) {
			name
			type
			pod: ~replicaset @filter(has(isPod)) {
				name
				type
				container: ~pod @filter(has(isContainer)) {
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

// RetrieveAllReplicasetsWithMetrics ...
func RetrieveAllReplicasetsWithMetrics() (ReplicasetsWithMetrics, error) {
	const q = `query {
		replicaset(func: has(isReplicaset)) {
			name
			type
			pod: ~replicaset @filter(has(isPod)) {
				name
				type
				container: ~pod @filter(has(isContainer)) {
					name
					type
					cpu: cpuRequest
					memory: memoryRequest
				}
				cpu: podCpu as cpuRequest
				memory: podMemory as memoryRequest
			}
			cpu: sum(val(podCpu))
			memory: sum(val(podMemory))
		}
	}`
	replicasetRoot := ReplicasetsWithMetrics{}
	err := dgraph.ExecuteQuery(q, &replicasetRoot)
	calculateTotalReplicasetMetrics(&replicasetRoot)
	return replicasetRoot, err
}

// RetrieveReplicasetWithMetrics ...
func RetrieveReplicasetWithMetrics(name string) (ReplicasetsWithMetrics, error) {
	q := `query {
		replicaset(func: has(isReplicaset)) @filter(eq(name, "` + name + `")) {
			name
			type
			pod: ~replicaset @filter(has(isPod)) {
				name
				type
				container: ~pod @filter(has(isContainer)) {
					name
					type
					cpu: cpuRequest
					memory: memoryRequest
				}
				cpu: podCpu as cpuRequest
				memory: podMemory as memoryRequest
			}
			cpu: sum(val(podCpu))
			memory: sum(val(podMemory))
		}
	}`
	replicasetRoot := ReplicasetsWithMetrics{}
	err := dgraph.ExecuteQuery(q, &replicasetRoot)
	calculateTotalReplicasetMetrics(&replicasetRoot)
	return replicasetRoot, err
}

func calculateTotalReplicasetMetrics(replicasetRoot *ReplicasetsWithMetrics) {
	for _, replicaset := range replicasetRoot.Replicaset {
		replicasetRoot.CPU += replicaset.CPU
		replicasetRoot.Memory += replicaset.Memory
	}
}