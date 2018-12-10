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

package v1

import (
	"github.com/vmware/purser/pkg/controller/metrics"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CRD Group attributes
const (
	CRDPlural   string = "groups"
	CRDGroup    string = "vmware.purser.com"
	CRDVersion  string = "v1"
	FullCRDName string = CRDPlural + "." + CRDGroup
)

// Group describes our custom Group resource
type Group struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               GroupSpec   `json:"spec"`
	Status             GroupStatus `json:"status,omitempty"`
}

// GroupSpec is the spec for the Group resource
type GroupSpec struct {
	Name               string                         `json:"name"`
	Type               string                         `json:"type,omitempty"`
	Expressions        map[string]map[string][]string `json:"labels,omitempty"`
	AllocatedResources *GroupMetrics                  `json:"metrics,omitempty"`
	PITMetrics         *GroupMetrics                  `json:"pitMetrics,omitempty"`
	MTDMetrics         *GroupMetrics                  `json:"mtdMetrics,omitempty"`
	MTDCost            *Cost                          `json:"mtdCost,omitempty"`
	PodsDetails        map[string]*PodDetails         `json:"podDetails,omitempty"`
}

// GroupMetrics ...
type GroupMetrics struct {
	CPULimit        float64
	MemoryLimit     float64
	StorageCapacity float64
	CPURequest      float64
	MemoryRequest   float64
	StorageClaim    float64
}

// Cost details
type Cost struct {
	TotalCost   float64
	CPUCost     float64
	MemoryCost  float64
	StorageCost float64
}

// GroupList is the list of Group resources
type GroupList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []*Group `json:"items"`
}

// GroupStatus holds the status information for each Group resource
type GroupStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// PodDetails information for the pods associated with the Group resource
type PodDetails struct {
	Name            string
	StartTime       meta_v1.Time
	EndTime         meta_v1.Time
	Containers      []*Container
	PodVolumeClaims map[string]*PersistentVolumeClaim
}

// PersistentVolumeClaim information for the pods associated with the Group resource
// A PVC can bound and unbound to a pod many times, so maintaining
// BoundTimes and UnboundTimes as lists.
// A PVC can be upgraded or downgraded, so maintaining capacityAllocated as a list
// Whenever a PVC capacity changes will update UnboundTime for old capacity, and
// append new capacity to capacityAllocated with bound time appended to BoundTimes
// The i-th capacity allocated corresponds to the i-th bound time and to i-th unbound time.
// Similarly for RequestSizeInGB
type PersistentVolumeClaim struct {
	Name                  string
	VolumeName            string
	RequestSizeInGB       []float64
	CapacityAllocatedInGB []float64
	BoundTimes            []meta_v1.Time
	UnboundTimes          []meta_v1.Time
}

// Container information for the pods associated with the Group resource
type Container struct {
	Name    string
	Metrics *metrics.Metrics
}
