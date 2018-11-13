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
	"github.com/vmware/purser/pkg/controller/utils"
	api_v1 "k8s.io/api/core/v1"
	"log"
)

// Dgraph Model Constants
const (
	IsPersistentVolume = "isPersistentVolume"
)

// PersistentVolume schema in dgraph
type PersistentVolume struct {
	dgraph.ID
	IsPersistentVolume bool       `json:"isPersistentVolume,omitempty"`
	Name               string     `json:"name,omitempty"`
	StartTime          time.Time  `json:"startTime,omitempty"`
	EndTime            time.Time  `json:"endTime,omitempty"`
	Type               string     `json:"type,omitempty"`
	StorageCapacity    float64    `json:"storageCapacity,omitempty"`
}

func createPersistentVolumeObject(pv api_v1.PersistentVolume) PersistentVolume {
	newPv := PersistentVolume{
		Name:               pv.Name,
		IsPersistentVolume: true,
		Type:               "pv",
		ID:                 dgraph.ID{Xid: pv.Name},
		StartTime:          pv.GetCreationTimestamp().Time,
	}
	capacity := pv.Spec.Capacity["storage"]
	newPv.StorageCapacity = utils.ConvertToFloat64GB(&capacity)

	deletionTimestamp := pv.GetDeletionTimestamp()
	if !deletionTimestamp.IsZero() {
		newPv.EndTime = deletionTimestamp.Time
	}
	return newPv
}

// StorePersistentVolume create a new persistent volume in the Dgraph and updates if already present.
func StorePersistentVolume(pv api_v1.PersistentVolume) (string, error) {
	xid := pv.Name
	uid := dgraph.GetUID(xid, IsPersistentVolume)

	newPv := createPersistentVolumeObject(pv)
	if uid != "" {
		newPv.UID = uid
	}
	assigned, err := dgraph.MutateNode(newPv, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	return assigned.Uids["blank-0"], nil
}

// CreateOrGetPersistentVolumeByID returns the uid of persistent volume if exists,
// otherwise creates the persistent volume and returns uid.
func CreateOrGetPersistentVolumeByID(xid string) string {
	if xid == "" {
		return ""
	}
	uid := dgraph.GetUID(xid, IsPersistentVolume)

	if uid != "" {
		return uid
	}

	d := PersistentVolume{
		ID:                 dgraph.ID{Xid: xid},
		Name:               xid,
		IsPersistentVolume: true,
	}
	assigned, err := dgraph.MutateNode(d, dgraph.CREATE)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return assigned.Uids["blank-0"]
}
