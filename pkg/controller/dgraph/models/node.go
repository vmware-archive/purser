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

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/utils"
	api_v1 "k8s.io/api/core/v1"
)

// Dgraph Model Constants
const (
	IsNode = "isNode"
)

// Node schema in dgraph
type Node struct {
	dgraph.ID
	IsNode         bool      `json:"isNode,omitempty"`
	Name           string    `json:"name,omitempty"`
	StartTime      string `json:"startTime,omitempty"`
	EndTime        string `json:"endTime,omitempty"`
	Pods           []*Pod    `json:"pods,omitempty"`
	CPUCapity      float64   `json:"cpuCapacity,omitempty"`
	MemoryCapacity float64   `json:"memoryCapacity,omitempty"`
	Type           string    `json:"type,omitempty"`
}

func createNodeObject(node api_v1.Node) Node {
	newNode := Node{
		Name:           node.Name,
		IsNode:         true,
		Type:           "node",
		ID:             dgraph.ID{Xid: node.Name},
		StartTime:      node.GetCreationTimestamp().Time.Format(time.RFC3339),
		CPUCapity:      utils.ConvertToFloat64CPU(node.Status.Capacity.Cpu()),
		MemoryCapacity: utils.ConvertToFloat64GB(node.Status.Capacity.Memory()),
	}
	nodeDeletionTimestamp := node.GetDeletionTimestamp()
	if !nodeDeletionTimestamp.IsZero() {
		newNode.EndTime = nodeDeletionTimestamp.Time.Format(time.RFC3339)
	}
	return newNode
}

// createOrGetNodeByID create and returns the node if not present, otherwise simply returns node.
func createOrGetNodeByID(xid string) (string, error) {
	if xid == "" {
		return "", fmt.Errorf("Node xid is empty")
	}
	uid := dgraph.GetUID(xid, IsNode)
	if uid != "" {
		return uid, nil
	}
	newNode := Node{
		Name:   xid,
		IsNode: true,
		ID:     dgraph.ID{Xid: xid},
	}
	assigned, err := dgraph.MutateNode(newNode, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	log.Infof("Node with xid: (%s) persisted", xid)
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

	if uid == "" {
		log.Infof("Node with xid: (%s) persisted", xid)
	}
	return assigned.Uids["blank-0"], nil
}

// RetrieveAllNodes ...
func RetrieveAllNodes() ([]byte, error) {
	const q = `query {
		result(func: has(isNode)) {
			name
			type
			~node @filter(has(isPod) {
				name
				type
				~pod @filter(has(isContainer)) {
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

// RetrieveNode ...
func RetrieveNode(name string) ([]byte, error) {
	q := `query {
		result(func: has(isNode)) @filter(eq(name, "` + name + `")) {
			name
			type
			~node @filter(has(isPod)) {
				name
				type
				~pod @filter(has(isContainer)) {
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

