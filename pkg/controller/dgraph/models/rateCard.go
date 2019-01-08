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

import "github.com/vmware/purser/pkg/controller/dgraph"

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

// TODO: store/update rate card in dgraph