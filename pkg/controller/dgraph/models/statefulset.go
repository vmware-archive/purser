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
)

// Dgraph Model Constants
const (
	IsStatefulset = "isStatefulset"
)

// Statefulset schema in dgraph
type Statefulset struct {
	dgraph.ID
	IsStatefulset bool      `json:"isStatefulset,omitempty"`
	Name          string    `json:"name,omitempty"`
	StartTime     time.Time `json:"startTime,omitempty"`
	EndTime       time.Time `json:"endTime,omitempty"`
	Pods          []*Pod    `json:"pods,omitempty"`
}

func createStatefulsetObject(statefulset apps_v1beta1.StatefulSet) Statefulset {
	newStatefulset := Statefulset{
		Name:         statefulset.Name,
		IsStatefulset: true,
		ID:           dgraph.ID{Xid: statefulset.Namespace + ":" + statefulset.Name},
		StartTime:    statefulset.GetCreationTimestamp().Time,
	}
	statefulsetDeletionTimestamp := statefulset.GetDeletionTimestamp()
	if !statefulsetDeletionTimestamp.IsZero() {
		newStatefulset.EndTime = statefulsetDeletionTimestamp.Time
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
