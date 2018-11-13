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
	IsStatefulset = "isStatefulset"
)

// Statefulset schema in dgraph
type Statefulset struct {
	dgraph.ID
	IsStatefulset bool       `json:"isStatefulset,omitempty"`
	Name          string     `json:"name,omitempty"`
	StartTime     string  `json:"startTime,omitempty"`
	EndTime       string  `json:"endTime,omitempty"`
	Namespace     *Namespace `json:"namespace,omitempty"`
	Pods          []*Pod     `json:"pods,omitempty"`
	Type          string     `json:"type,omitempty"`
}

func createStatefulsetObject(statefulset apps_v1beta1.StatefulSet) Statefulset {
	newStatefulset := Statefulset{
		Name:          statefulset.Name,
		IsStatefulset: true,
		Type:          "statefulset",
		ID:            dgraph.ID{Xid: statefulset.Namespace + ":" + statefulset.Name},
		StartTime:     statefulset.GetCreationTimestamp().Time.Format(time.RFC3339),
	}
	namespaceUID := CreateOrGetNamespaceByID(statefulset.Namespace)
	if namespaceUID != "" {
		newStatefulset.Namespace = &Namespace{ID: dgraph.ID{UID: namespaceUID, Xid: statefulset.Namespace}}
	}
	statefulsetDeletionTimestamp := statefulset.GetDeletionTimestamp()
	if !statefulsetDeletionTimestamp.IsZero() {
		newStatefulset.EndTime = statefulsetDeletionTimestamp.Time.Format(time.RFC3339)
	}
	return newStatefulset
}

// StoreStatefulset create a new statefulset in the Dgraph and updates if already present.
func StoreStatefulset(statefulset apps_v1beta1.StatefulSet) (string, error) {
	xid := statefulset.Namespace + ":" + statefulset.Name
	uid := dgraph.GetUID(xid, IsStatefulset)

	newStatefulset := createStatefulsetObject(statefulset)
	if uid != "" {
		newStatefulset.UID = uid
	}
	assigned, err := dgraph.MutateNode(newStatefulset, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	return assigned.Uids["blank-0"], nil
}

// CreateOrGetStatefulsetByID returns the uid of namespace if exists,
// otherwise creates the stateful and returns uid.
func CreateOrGetStatefulsetByID(xid string) string {
	if xid == "" {
		return ""
	}
	uid := dgraph.GetUID(xid, IsStatefulset)

	if uid != "" {
		return uid
	}

	d := Statefulset{
		ID:            dgraph.ID{Xid: xid},
		Name:          xid,
		IsStatefulset: true,
	}
	assigned, err := dgraph.MutateNode(d, dgraph.CREATE)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return assigned.Uids["blank-0"]
}

// RetrieveAllStatefulsets ...
func RetrieveAllStatefulsets() ([]byte, error) {
	const q = `query {
		statefulset(func: has(isStatefulset)) {
			name
			type
			pod: ~statefulset @filter(has(isPod)) {
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

// RetrieveStatefulset ...
func RetrieveStatefulset(name string) ([]byte, error) {
	q := `query {
		statefulset(func: has(isStatefulset)) @filter(eq(name, "` + name + `")) {
			name
			type
			pod: ~statefulset @filter(has(isPod)) {
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

// RetrieveStatefulsetWithMetrics ...
func RetrieveStatefulsetWithMetrics(name string) (JsonDataWrapper, error) {
	q := `query {
		parent(func: has(isStatefulset)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~statefulset @filter(has(isPod)) {
				name
				type
				cpu: podCpu as cpuRequest
				memory: podMemory as memoryRequest
			}
			cpu: sum(val(podCpu))
			memory: sum(val(podMemory))
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