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

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/purser/pkg/controller/dgraph"

	api_v1 "k8s.io/api/core/v1"
)

// Dgraph Model Constants
const (
	IsNamespace = "isNamespace"
)

// Namespace schema in dgraph
type Namespace struct {
	dgraph.ID
	IsNamespace bool   `json:"isNamespace,omitempty"`
	Name        string `json:"name,omitempty"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
	Type        string `json:"type,omitempty"`
}

func newNamespace(namespace api_v1.Namespace) Namespace {
	ns := Namespace{
		ID:          dgraph.ID{Xid: namespace.Name},
		Name:        "namespace-" + namespace.Name,
		IsNamespace: true,
		Type:        "namespace",
		StartTime:   namespace.GetCreationTimestamp().Time.Format(time.RFC3339),
	}
	nsDeletionTimestamp := namespace.GetDeletionTimestamp()
	if !nsDeletionTimestamp.IsZero() {
		ns.EndTime = nsDeletionTimestamp.Time.Format(time.RFC3339)
		ns.Xid += ns.EndTime
		ns.Name += "*" + ns.EndTime
	}
	return ns
}

// CreateOrGetNamespaceByID returns the uid of namespace if exists,
// otherwise creates the namespace and returns uid.
func CreateOrGetNamespaceByID(xid string) string {
	if xid == "" {
		log.Error("Namespace is empty")
		return ""
	}
	uid := dgraph.GetUID(xid, IsNamespace)

	if uid != "" {
		return uid
	}

	ns := Namespace{
		ID:          dgraph.ID{Xid: xid},
		Name:        xid,
		IsNamespace: true,
	}
	assigned, err := dgraph.MutateNode(ns, dgraph.CREATE)
	if err != nil {
		log.Error(err)
		return ""
	}
	log.Infof("Namespace with xid: (%s) persisted", xid)
	return assigned.Uids["blank-0"]
}

// StoreNamespace create a new namespace in the Dgraph  if it is not present.
func StoreNamespace(namespace api_v1.Namespace) (string, error) {
	xid := namespace.Name
	uid := dgraph.GetUID(xid, IsNamespace)

	ns := newNamespace(namespace)
	if uid != "" {
		ns.UID = uid
	}
	assigned, err := dgraph.MutateNode(ns, dgraph.CREATE)
	if err != nil {
		return "", err
	}

	if uid == "" {
		log.Infof("Namespace with xid: (%s) persisted", xid)
	}
	return assigned.Uids["blank-0"], nil
}
