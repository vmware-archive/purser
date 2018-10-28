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
	api_v1 "k8s.io/api/core/v1"
)

// Dgraph Model Constants
const (
	IsNode = "isNode"
)

// Node schema in dgraph
type Node struct {
	dgraph.ID
	IsNode    bool      `json:"isNode,omitempty"`
	Name      string    `json:"name,omitempty"`
	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
	Pods      []*Pod    `json:"pods,omitempty"`
}

func createNodeObject(node api_v1.Node) Node {
	newNode := Node{
		Name:      node.Name,
		IsNode:    true,
		ID:        dgraph.ID{Xid: node.Name},
		StartTime: node.GetCreationTimestamp().Time,
	}
	nodeDeletionTimestamp := node.GetDeletionTimestamp()
	if !nodeDeletionTimestamp.IsZero() {
		newNode.EndTime = nodeDeletionTimestamp.Time
	}
	return newNode
}

// CreateOrGetNodeByID create and returns the node if not present, otherwise simply returns node.
func CreateOrGetNodeByID(xid string) (string, error) {
	uid := dgraph.GetUID(xid, IsNode)
	if uid != "" {
		return uid, nil
	}
	newNode := Node{
		Name:	xid,
		IsNode: true,
		ID:     dgraph.ID{Xid: xid},
	}
	assigned, err := dgraph.MutateNode(newNode, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	return assigned.Uids["blank-0"], nil
}

// StoreNode create a new node in the Dgraph  if it is not present.
func StoreNode(node api_v1.Node) (string, error) {
	xid := node.Name
	uid := dgraph.GetUID(xid, IsNode)

	newNode := createNodeObject(node)
	if uid != "" {
		newNode.UID = uid
	}
	assigned, err := dgraph.MutateNode(newNode, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	return assigned.Uids["blank-0"], nil
}
