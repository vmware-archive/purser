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

package query

import (
	"strconv"

	"github.com/Sirupsen/logrus"
)

// Cluster resource constants
const (
	ContainerCheck = "isContainer"
	ContainerType  = "container"
	IsProcFilter   = "@filter(has(isProc))"

	DaemonsetCheck = "isDaemonset"
	DaemonsetType  = "daemonset"
	IsPodFilter    = "@filter(has(isPod))"

	DeploymentCheck    = "isDeployment"
	DeploymentType     = "deployment"
	IsReplicasetFilter = "@filter(has(isReplicaset))"

	JobCheck = "isJob"
	JobType  = "job"

	NamespaceCheck       = "isNamespace"
	NamespaceType        = "namespace"
	NamespaceChildFilter = "@filter(has(isDeployment) OR has(isStatefulset) OR has(isJob) OR has(isDaemonset) OR (has(isReplicaset) AND (NOT has(deployment))))"

	NodeCheck = "isNode"
	NodeType  = "node"

	PodCheck          = "isPod"
	PodType           = "pod"
	IsContainerFilter = "@filter(has(isContainer))"

	PVCheck     = "isPersistentVolume"
	PVType      = "pv"
	IsPVCFilter = "@filter(has(isPersistentVolumeClaim))"

	PVCCheck = "isPersistentVolumeClaim"
	PVCType  = "pvc"

	ReplicasetCheck = "isReplicaset"
	ReplicasetType  = "replicaset"

	StatefulsetCheck = "isStatefulset"
	StatefulsetType  = "statefulset"
)

// Resource structure
type Resource struct {
	Check       string
	Type        string
	Name        string
	ChildFilter string
}

// RetrieveResourceHierarchy returns hierarchy for a given resource
func (r *Resource) RetrieveResourceHierarchy() JSONDataWrapper {
	if r.Name == All {
		logrus.Errorf("wrong type of query, empty name is given")
		return JSONDataWrapper{}
	}
	query := r.getQueryForHierarchy()
	return getJSONDataFromQuery(query)
}

// RetrieveResourceMetrics returns metrics for a given resource
func (r *Resource) RetrieveResourceMetrics() JSONDataWrapper {
	if r.Name == All {
		logrus.Errorf("wrong type of query, empty name is given")
		return JSONDataWrapper{}
	}
	query := r.getQueryForResourceMetrics()
	return getJSONDataFromQuery(query)
}

func (r *Resource) getQueryForResourceMetrics() string {
	switch r.Type {
	case DeploymentType:
		return getQueryForDeploymentMetrics(r.Name)
	case NamespaceType:
		return getQueryForNamespaceMetrics(r.Name)
	case NodeType:
		return getQueryForNodeMetrics(r.Name)
	case PVType:
		return getQueryForPVMetrics(r.Name)
	case PVCType:
		return getQueryForPVCMetrics(r.Name)
	case ContainerType:
		return getQueryForContainerMetrics(r.Name)
	case PodType:
		cpuPriceInFloat64, memoryPriceInFloat64 := getPricePerResourceForPod(r.Name)
		cpuPrice := strconv.FormatFloat(cpuPriceInFloat64, 'f', 11, 64)
		memoryPrice := strconv.FormatFloat(memoryPriceInFloat64, 'f', 11, 64)
		return getQueryForPodMetrics(r.Name, cpuPrice, memoryPrice)
	}
	return r.getQueryForPodParentMetrics()
}

// getJSONDataFromQuery executes query and wraps the data in a desired structure(JSONDataWrapper)
func getJSONDataFromQuery(query string) JSONDataWrapper {
	parentRoot := ParentWrapper{}
	err := executeQuery(query, &parentRoot)
	if err != nil || len(parentRoot.Parent) == 0 {
		logrus.Errorf("Unable to execute query, err: (%v)", err)
		return JSONDataWrapper{}
	}
	root := JSONDataWrapper{
		Data: parentRoot.Parent[0],
	}
	return root
}
