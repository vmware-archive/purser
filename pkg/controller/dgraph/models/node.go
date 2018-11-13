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
	CPU    float64    `json:"cpu,omitempty"`
	Memory float64    `json:"memory,omitempty"`
	Children []*Children `json:"children,omitempty"`
}

// NodesWithMetrics ...
type NodesWithMetrics struct {
	Node []Node  `json:"node,omitempty"`
	CPU    float64    `json:"cpu,omitempty"`
	Memory float64    `json:"memory,omitempty"`
}

type PhysicalCluster struct {
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	Children []Node `json:"children,omitempty"`
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
	log.Infof("cpu capacity of node: (%v)", newNode.CPUCapity)
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
		node(func: has(isNode)) {
			name
			type
			pod: ~node @filter(has(isPod)) {
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

// RetrieveNode ...
func RetrieveNode(name string) ([]byte, error) {
	q := `query {
		node(func: has(isNode)) @filter(eq(name, "` + name + `")) {
			name
			type
			pod: ~node @filter(has(isPod)) {
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

// RetrieveAllNodesWithMetrics ...
func RetrieveAllNodesWithMetrics() (ClusterWithMetrics, error) {
	const q = `query {
		nd as var(func: has(isNode)) {
			~node @filter(has(isPod)){
				nodePodCpu as cpuRequest
				nodePodMem as memoryRequest
			}
			nodeCpu as sum(val(nodePodCpu))
			nodeMem as sum(val(nodePodMem))
        }

		node(func: uid(nd)) {
			name
            type
			cpu: val(nodeCpu)
			memory: val(nodeMem)
        }
	}`
	nodeRoot := NodesWithMetrics{}
	err := dgraph.ExecuteQuery(q, &nodeRoot)
	calculateTotalNodeMetrics(&nodeRoot)
	clusterRoot := ClusterWithMetrics{}
	clusterRoot.CPU = nodeRoot.CPU
	clusterRoot.Memory = nodeRoot.Memory
	clusterRoot.PhysicalCluster = append(clusterRoot.PhysicalCluster, PhysicalCluster{
		Name: "cluster",
		Type: "cluster",
		Children: nodeRoot.Node,
	})
	return clusterRoot, err
}

// RetrieveNodeWithMetrics ...
func RetrieveNodeWithMetrics(name string) (NodesWithMetrics, error) {
	q := `query {
		node(func: has(isNode)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~node @filter(has(isPod)) {
				name
				type
				cpu: puRequest
				memory: memoryRequest
			}
			cpu: cpuCapacity
			memory: memoryCapacity
		}
	}`
	nodeRoot := NodesWithMetrics{}
	err := dgraph.ExecuteQuery(q, &nodeRoot)
	calculateTotalNodeMetrics(&nodeRoot)
	return nodeRoot, err
}

func calculateTotalNodeMetrics(objRoot *NodesWithMetrics) {
	for _, obj := range objRoot.Node {
		objRoot.CPU += obj.CPU
		objRoot.Memory += obj.Memory
	}
}
