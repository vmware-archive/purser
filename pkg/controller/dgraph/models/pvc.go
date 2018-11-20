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
	IsPersistentVolumeClaim = "isPersistentVolumeClaim"
)

// PersistentVolumeClaim schema in dgraph
type PersistentVolumeClaim struct {
	dgraph.ID
	IsPersistentVolumeClaim bool              `json:"isPersistentVolumeClaim,omitempty"`
	Name                    string            `json:"name,omitempty"`
	StartTime               string            `json:"startTime,omitempty"`
	EndTime                 string            `json:"endTime,omitempty"`
	Namespace               *Namespace        `json:"namespace,omitempty"`
	Type                    string            `json:"type,omitempty"`
	StorageCapacity         float64           `json:"storageCapacity,omitempty"`
	PersistentVolume        *PersistentVolume `json:"persistentvolume,omitempty"`
}

func createPvcObject(pvc api_v1.PersistentVolumeClaim) PersistentVolumeClaim {
	newPvc := PersistentVolumeClaim{
		Name:                    pvc.Name,
		IsPersistentVolumeClaim: true,
		Type:                    "pvc",
		ID:                      dgraph.ID{Xid: pvc.Namespace + ":" + pvc.Name},
		StartTime:               pvc.GetCreationTimestamp().Time.Format(time.RFC3339),
	}
	capacity := pvc.Status.Capacity["storage"]
	newPvc.StorageCapacity = utils.ConvertToFloat64GB(&capacity)

	volume := pvc.Spec.VolumeName
	pvUID := CreateOrGetPersistentVolumeByID(volume)
	if volume != "" {
		newPvc.PersistentVolume = &PersistentVolume{ID: dgraph.ID{UID: pvUID}}
	}

	namespaceUID := CreateOrGetNamespaceByID(pvc.Namespace)
	if namespaceUID != "" {
		newPvc.Namespace = &Namespace{ID: dgraph.ID{UID: namespaceUID, Xid: pvc.Namespace}}
	}
	deletionTimestamp := pvc.GetDeletionTimestamp()
	if !deletionTimestamp.IsZero() {
		newPvc.EndTime = deletionTimestamp.Time.Format(time.RFC3339)
	}
	return newPvc
}

// StorePersistentVolumeClaim create a new pvc in the Dgraph and updates if already present.
func StorePersistentVolumeClaim(pvc api_v1.PersistentVolumeClaim) (string, error) {
	xid := pvc.Namespace + ":" + pvc.Name
	uid := dgraph.GetUID(xid, IsPersistentVolumeClaim)

	newPvc := createPvcObject(pvc)
	if uid != "" {
		newPvc.UID = uid
	}
	assigned, err := dgraph.MutateNode(newPvc, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	return assigned.Uids["blank-0"], nil
}

// CreateOrGetPersistentVolumeClaimByID returns the uid of pvc if exists,
// otherwise creates the pvc and returns uid.
func CreateOrGetPersistentVolumeClaimByID(xid string) string {
	if xid == "" {
		return ""
	}
	uid := dgraph.GetUID(xid, IsPersistentVolumeClaim)

	if uid != "" {
		return uid
	}

	d := PersistentVolumeClaim{
		ID:                      dgraph.ID{Xid: xid},
		Name:                    xid,
		IsPersistentVolumeClaim: true,
	}
	assigned, err := dgraph.MutateNode(d, dgraph.CREATE)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return assigned.Uids["blank-0"]
}

func getPVCFromUID(uid string) (PersistentVolumeClaim, error) {
	q := `query {
		pvcs(func: uid(` + uid + `)) {
			name
			type
			storageCapacity
		}
	}`

	type root struct {
		Pvcs []PersistentVolumeClaim `json:"pvcs"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return PersistentVolumeClaim{}, err
	}
	return newRoot.Pvcs[0], nil
}
