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
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
)

// RateCard constants
const (
	IsRateCard     = "isRateCard"
	IsNodePrice    = "isNodePrice"
	IsStoragePrice = "isStoragePrice"
	RateCardXID    = "purser-rateCard"
)

// RateCard structure
type RateCard struct {
	dgraph.ID
	IsRateCard    bool            `json:"isRateCard,omitempty"`
	CloudProvider string          `json:"cloudProvider,omitempty"`
	Region        string          `json:"region,omitempty"`
	NodePrices    []*NodePrice    `json:"nodePrices,omitempty"`
	StoragePrices []*StoragePrice `json:"storagePrices,omitempty"`
}

// NodePrice structure
// Unit of Node Price should be USD($)-(per Hour)
type NodePrice struct {
	dgraph.ID
	IsNodePrice     bool    `json:"isNodePrice,omitempty"`
	InstanceType    string  `json:"instanceType,omitempty"`
	InstanceFamily  string  `json:"instanceFamily,omitempty"`
	OperatingSystem string  `json:"operatingSystem,omitempty"`
	Price           float64 `json:"price,omitempty"`
	PricePerCPU     float64 `json:"cpuPrice,omitempty"`
	PricePerMemory  float64 `json:"memoryPrice,omitempty"`
}

// StoragePrice structure
// Unit of Storage Price should be USD($)-(per GB)-(per Hour)
type StoragePrice struct {
	dgraph.ID
	IsStoragePrice bool    `json:"isStoragePrice,omitempty"`
	VolumeType     string  `json:"volumeType,omitempty"`
	UsageType      string  `json:"usageType,omitempty"`
	Price          float64 `json:"price,omitempty"`
}

// StoreRateCard given a cloudProvider and region it gets rate card and stores(create/update) in dgraph
func StoreRateCard(rateCard *RateCard) {
	logrus.Debugf("IsRateCardNil: %v", rateCard == nil)
	if rateCard != nil {
		uid := dgraph.GetUID(RateCardXID, IsRateCard)
		if uid != "" {
			rateCard.ID = dgraph.ID{UID: uid, Xid: RateCardXID}
		}
		logrus.Debugf("RateCard: (%v)", rateCard)
		_, err := dgraph.MutateNode(rateCard, dgraph.CREATE)
		if err != nil {
			logrus.Errorf("Unable to store rateCard reason: %v", err)
			return
		}
		logrus.Infof("Successfully stored/updated rateCard")
	}
}

// StoreNodePrice given nodePrice and its XID it stores(create/update) in dgraph
func StoreNodePrice(nodePrice *NodePrice, productXID string) string {
	uid := dgraph.GetUID(productXID, IsNodePrice)
	if uid != "" {
		nodePrice.ID = dgraph.ID{Xid: productXID, UID: uid}
	}
	logrus.Debugf("nodePrice: %v, productXID: %v\n", *nodePrice, productXID)
	assigned, err := dgraph.MutateNode(nodePrice, dgraph.CREATE)
	if err != nil {
		logrus.Errorf("Unable to store nodePrice: (%v), reason: %v", nodePrice, err)
		return ""
	}
	logrus.Debugf("Successfully stored/updated nodePrice: %v", productXID)

	if uid != "" {
		return uid
	}
	return assigned.Uids["blank-0"]
}

// StoreStoragePrice given storagePrice and its XID it stores(create/update) in dgraph
func StoreStoragePrice(storagePrice *StoragePrice, productXID string) string {
	uid := dgraph.GetUID(productXID, IsStoragePrice)
	if uid != "" {
		storagePrice.ID = dgraph.ID{Xid: productXID, UID: uid}
	}
	logrus.Debugf("storagePrice: %v, productXID: %v\n", *storagePrice, productXID)
	assigned, err := dgraph.MutateNode(storagePrice, dgraph.CREATE)
	if err != nil {
		logrus.Errorf("Unable to store storagePrice: (%v), reason: %v", storagePrice, err)
		return ""
	}
	logrus.Debugf("Successfully stored/updated storagePrice: %v", productXID)

	if uid != "" {
		return uid
	}
	return assigned.Uids["blank-0"]
}

// retrieveNode given a node name it returns pointer to models.Node - nil in case of error
func retrieveNode(name string) (*Node, error) {
	query := `query {
		nodes(func: has(isNode)) @filter(eq(name, "` + name + `")) {
			name
			type
			startTime
			endTime
			cpuCapacity
			memoryCapacity
			instanceType
			os
        }
    }`
	type root struct {
		Nodes []Node `json:"nodes"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(query, &newRoot)
	if err != nil {
		return nil, err
	} else if len(newRoot.Nodes) < 1 {
		return nil, fmt.Errorf("no node with name: %v", name)
	}

	return &newRoot.Nodes[0], nil
}

// retrieveNodePrice given a node name it returns pointer to models.Node - nil in case of error
func retrieveNodePrice(xid string) (*NodePrice, error) {
	query := `query {
		nodePrices(func: has(isNodePrice)) @filter(eq(xid, "` + xid + `")) {
			instanceType
			instanceFamily
			operatingSystem
			price
			cpuPrice
			memoryPrice
        }
    }`
	type root struct {
		NodePrices []NodePrice `json:"nodePrices"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(query, &newRoot)
	if err != nil {
		return nil, err
	} else if len(newRoot.NodePrices) < 1 {
		return nil, fmt.Errorf("no node with xid: %v", xid)
	}

	return &newRoot.NodePrices[0], nil
}

// getPerUnitResourcePriceForNode returns price per cpu and price per memory
func getPerUnitResourcePriceForNode(nodeName string) (float64, float64) {
	node, err := retrieveNode(nodeName)
	if err == nil {
		return getPricePerUnitResourceFromNodePrice(*node)
	}
	return DefaultCPUCostInFloat64, DefaultMemCostInFloat64
}

func getPricePerUnitResourceFromNodePrice(node Node) (float64, float64) {
	xidsToTry := []string{
		node.InstanceType + "-" + node.OS,
		node.InstanceType + "-linux",
		node.InstanceType + "-ANY",
		node.InstanceType,
	}
	for _, xid := range xidsToTry {
		nodePrice, err := retrieveNodePrice(xid)
		if err == nil {
			return nodePrice.PricePerCPU, nodePrice.PricePerMemory
		}
	}
	return DefaultCPUCostInFloat64, DefaultMemCostInFloat64
}
