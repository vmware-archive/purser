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
	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
)

// RateCard constants
const (
	IsRateCard     = "isRateCard"
	IsNodePrice    = "isNodePrice"
	IsStoragePrice = "isStoragePrice"
	RateCardXID    = "purser-rateCard"

	AWS = "aws"
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
	PricePerCPU     string  `json:"cpuPrice,omitempty"`
	PricePerMemory  string  `json:"memoryPrice,omitempty"`
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
