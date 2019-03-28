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
	"github.com/dgraph-io/dgo/protos/api"
	groups_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/utils"
)

// Group constants
const (
	IsGroup        = "isGroup"
	groupXIDPrefix = "purser-group-"
)

// Group schema in dgraph
type Group struct {
	dgraph.ID
	IsGroup                  bool    `json:"isGroup,omitempty"`
	Name                     string  `json:"name,omitempty"`
	PodsCount                int     `json:"podsCount,omitempty"`
	MtdCPU                   float64 `json:"mtdCPU,omitempty"`
	MtdMemory                float64 `json:"mtdMemory,omitempty"`
	MtdStorage               float64 `json:"mtdStorage,omitempty"`
	CPU                      float64 `json:"cpu,omitempty"`
	Memory                   float64 `json:"memory,omitempty"`
	Storage                  float64 `json:"storage,omitempty"`
	MtdCPUCost               float64 `json:"mtdCPUCost,omitempty"`
	MtdMemoryCost            float64 `json:"mtdMemoryCost,omitempty"`
	MtdStorageCost           float64 `json:"mtdStorageCost,omitempty"`
	MtdCost                  float64 `json:"mtdCost,omitempty"`
	ProjectedCPUCost         float64 `json:"projectedCPUCost,omitempty"`
	ProjectedMemoryCost      float64 `json:"projectedMemoryCost,omitempty"`
	ProjectedStorageCost     float64 `json:"projectedStorageCost,omitempty"`
	ProjectedCost            float64 `json:"projectedCost,omitempty"`
	LastMonthCPUCost         float64 `json:"lastMonthCPUCost,omitempty"`
	LastMonthMemoryCost      float64 `json:"lastMonthMemoryCost,omitempty"`
	LastMonthStorageCost     float64 `json:"lastMonthStorageCost,omitempty"`
	LastMonthCost            float64 `json:"lastMonthCost,omitempty"`
	LastLastMonthCPUCost     float64 `json:"lastLastMonthCPUCost,omitempty"`
	LastLastMonthMemoryCost  float64 `json:"lastLastMonthMemoryCost,omitempty"`
	LastLastMonthStorageCost float64 `json:"lastLastMonthStorageCost,omitempty"`
	LastLastMonthCost        float64 `json:"lastLastMonthCost,omitempty"`
}

// CreateOrUpdateGroup updates group if it is already present in dgraph else it creates one
func CreateOrUpdateGroup(group *groups_v1.Group, podsCount int) (*api.Assigned, error) {
	xid := groupXIDPrefix + group.Name
	uid := dgraph.GetUID(xid, IsGroup)

	hoursRemainingInCurrentMonth := utils.GetHoursRemainingInCurrentMonth()
	grp := Group{
		ID:                   dgraph.ID{Xid: xid},
		IsGroup:              true,
		Name:                 group.Name,
		PodsCount:            podsCount,
		MtdCPU:               group.Spec.MTDMetrics.CPURequest,
		MtdMemory:            group.Spec.MTDMetrics.MemoryRequest,
		MtdStorage:           group.Spec.MTDMetrics.StorageClaim,
		CPU:                  group.Spec.PITMetrics.CPURequest,
		Memory:               group.Spec.PITMetrics.MemoryRequest,
		Storage:              group.Spec.PITMetrics.StorageClaim,
		MtdCPUCost:           group.Spec.MTDCost.CPUCost,
		MtdMemoryCost:        group.Spec.MTDCost.MemoryCost,
		MtdStorageCost:       group.Spec.MTDCost.StorageCost,
		MtdCost:              group.Spec.MTDCost.TotalCost,
		ProjectedCPUCost:     group.Spec.MTDCost.CPUCost + group.Spec.PerHourCost.CPUCost*hoursRemainingInCurrentMonth,
		ProjectedMemoryCost:  group.Spec.MTDCost.MemoryCost + group.Spec.PerHourCost.MemoryCost*hoursRemainingInCurrentMonth,
		ProjectedStorageCost: group.Spec.MTDCost.StorageCost + group.Spec.PerHourCost.StorageCost*hoursRemainingInCurrentMonth,
		ProjectedCost:        group.Spec.MTDCost.TotalCost + group.Spec.PerHourCost.TotalCost*hoursRemainingInCurrentMonth,
	}
	if uid != "" {
		grp.ID = dgraph.ID{Xid: xid, UID: uid}
	}
	return dgraph.MutateNode(grp, dgraph.CREATE)
}

// DeleteGroup deletes group from dgraph
func DeleteGroup(name string) {
	xid := groupXIDPrefix + name
	uid := dgraph.GetUID(xid, IsGroup)

	if uid != "" {
		grp := Group{ID: dgraph.ID{UID: uid}}
		_, err := dgraph.MutateNode(grp, dgraph.DELETE)
		if err != nil {
			logrus.Errorf("error while deleting group: %v, err: %v", name, err)
		}
		return
	}
	logrus.Infof("Group: %s not yet persisted", name)
}
