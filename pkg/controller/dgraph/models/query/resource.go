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

// RetrieveResourceHierarchy returns hierarchy for a given resource
func RetrieveResourceHierarchy(resourceCheck, resourceType, resourceName, childFilter string) JSONDataWrapper {
	if resourceName == All {
		logrus.Errorf("wrong type of query, empty name is given")
		return JSONDataWrapper{}
	}
	query := getQueryForHierarchy(resourceCheck, resourceType, resourceName, childFilter)
	return getJSONDataFromQuery(query)
}

// RetrieveResourceMetrics returns metrics for a given resource
func RetrieveResourceMetrics(resourceCheck, resourceType, resourceName string) JSONDataWrapper {
	if resourceName == All {
		logrus.Errorf("wrong type of query, empty name is given")
		return JSONDataWrapper{}
	}
	query := getQueryForResourceMetrics(resourceCheck, resourceType, resourceName)
	return getJSONDataFromQuery(query)
}

func getQueryForResourceMetrics(resourceCheck, resourceType, resourceName string) string {
	switch resourceType {
	case DeploymentType:
		return getQueryForDeploymentMetrics(resourceName)
	case NamespaceType:
		return getQueryForNamespaceMetrics(resourceName)
	case NodeType:
		return getQueryForNodeMetrics(resourceName)
	case PVType:
		return getQueryForPVMetrics(resourceName)
	case PVCType:
		return getQueryForPVCMetrics(resourceName)
	case ContainerType:
		return getQueryForContainerMetrics(resourceName)
	case PodType:
		cpuPriceInFloat64, memoryPriceInFloat64 := getPricePerResourceForPod(resourceName)
		cpuPrice := strconv.FormatFloat(cpuPriceInFloat64, 'f', 11, 64)
		memoryPrice := strconv.FormatFloat(memoryPriceInFloat64, 'f', 11, 64)
		return getQueryForPodMetrics(resourceName, cpuPrice, memoryPrice)
	}
	return getQueryForPodParentMetrics(resourceCheck, resourceType, resourceName)
}

// getJSONDataFromQuery executes query and wraps the data in a desired structure(JSONDataWrapper)
func getJSONDataFromQuery(query string) JSONDataWrapper {
	parentRoot := ParentWrapper{}
	err := executeQuery(query, &parentRoot)
	if err != nil || len(parentRoot.Parent) == 0 {
		logrus.Errorf("Unable to execute query, err: (%v), length of output: (%d)", err, len(parentRoot.Parent))
		return JSONDataWrapper{}
	}
	root := JSONDataWrapper{
		Data: ParentWrapper{
			Name:        parentRoot.Parent[0].Name,
			Type:        parentRoot.Parent[0].Type,
			Children:    parentRoot.Parent[0].Children,
			CPU:         parentRoot.Parent[0].CPU,
			Memory:      parentRoot.Parent[0].Memory,
			Storage:     parentRoot.Parent[0].Storage,
			CPUCost:     parentRoot.Parent[0].CPUCost,
			MemoryCost:  parentRoot.Parent[0].MemoryCost,
			StorageCost: parentRoot.Parent[0].StorageCost,
		},
	}
	return root
}
