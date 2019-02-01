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
	"github.com/dgraph-io/dgo/protos/api"
	groups_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	"github.com/vmware/purser/pkg/controller/dgraph"
)

const (
	IsGroup   = "isGroup"
	XIDPrefix = "purser-group-"
)

// Group schema in dgraph
type Group struct {
	dgraph.ID
	IsGroup        bool    `json:"isGroup,omitempty"`
	Name           string  `json:"name,omitempty"`
	PodsCount      int     `json:"podsCount,omitempty"`
	MtdCPU         float64 `json:"mtdCPU,omitempty"`
	MtdMemory      float64 `json:"mtdMemory,omitempty"`
	MtdStorage     float64 `json:"mtdStorage,omitempty"`
	CPU            float64 `json:"cpu,omitempty"`
	Memory         float64 `json:"memory,omitempty"`
	Storage        float64 `json:"storage,omitempty"`
	MtdCPUCost     float64 `json:"mtdCPUCost,omitempty"`
	MtdMemoryCost  float64 `json:"mtdMemoryCost,omitempty"`
	MtdStorageCost float64 `json:"mtdStorageCost,omitempty"`
	MtdCost        float64 `json:"mtdCost,omitempty"`
}

// CreateOrUpdateGroup updates group if it is already present in dgraph else it creates one
func CreateOrUpdateGroup(group *groups_v1.Group, podsCount int) (*api.Assigned, error) {
	xid := XIDPrefix + group.Name
	uid := dgraph.GetUID(xid, IsGroup)

	grp := Group{
		ID:             dgraph.ID{Xid: xid},
		IsGroup:        true,
		Name:           group.Name,
		PodsCount:      podsCount,
		MtdCPU:         group.Spec.MTDMetrics.CPURequest,
		MtdMemory:      group.Spec.MTDMetrics.MemoryRequest,
		MtdStorage:     group.Spec.MTDMetrics.StorageClaim,
		CPU:            group.Spec.PITMetrics.CPURequest,
		Memory:         group.Spec.PITMetrics.MemoryRequest,
		Storage:        group.Spec.PITMetrics.StorageClaim,
		MtdCPUCost:     group.Spec.MTDCost.CPUCost,
		MtdMemoryCost:  group.Spec.MTDCost.MemoryCost,
		MtdStorageCost: group.Spec.MTDCost.StorageCost,
		MtdCost:        group.Spec.MTDCost.TotalCost,
	}
	if uid != "" {
		grp.ID = dgraph.ID{Xid: xid, UID: uid}
	}
	return dgraph.MutateNode(grp, dgraph.CREATE)
}
