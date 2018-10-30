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
)

// Dgraph Model Constants
const (
	IsReplicaset = "isReplicaset"
)

// Replicaset schema in dgraph
type Replicaset struct {
	dgraph.ID
	IsReplicaset bool      `json:"isReplicaset,omitempty"`
	Name         string    `json:"name,omitempty"`
	StartTime    time.Time `json:"startTime,omitempty"`
	EndTime      time.Time `json:"endTime,omitempty"`
	Pods         []*Pod    `json:"pods,omitempty"`
	Type         string    `json:"type,omitempty"`
}

func createReplicasetObject(replicaset ext_v1beta1.ReplicaSet) Replicaset {
	newReplicaset := Replicaset{
		Name:         replicaset.Name,
		IsReplicaset: true,
		Type:         "replicaset",
		ID:           dgraph.ID{Xid: replicaset.Namespace + ":" + replicaset.Name},
		StartTime:    replicaset.GetCreationTimestamp().Time,
	}
	replicasetDeletionTimestamp := replicaset.GetDeletionTimestamp()
	if !replicasetDeletionTimestamp.IsZero() {
		newReplicaset.EndTime = replicasetDeletionTimestamp.Time
	}
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
