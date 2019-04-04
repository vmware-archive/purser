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

package api

import (
	"github.com/vmware/purser/cmd/controller/api/apiHandlers"
	"net/http"
)

// Route structure
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes list
type Routes []Route

var routes = Routes{
	Route{
		"GetHomePage",
		"GET",
		"/api",
		apiHandlers.GetHomePage,
	},
	Route{
		"GetPodInteractions",
		"GET",
		"/api/interactions/pod",
		apiHandlers.GetPodInteractions,
	},
	Route{
		"GetClusterHierarchy",
		"GET",
		"/api/hierarchy",
		apiHandlers.GetClusterHierarchy,
	},
	Route{
		"GetNamespaceHierarchy",
		"GET",
		"/api/hierarchy/namespace",
		apiHandlers.GetNamespaceHierarchy,
	},
	Route{
		"GetDeploymentHierarchy",
		"GET",
		"/api/hierarchy/deployment",
		apiHandlers.GetDeploymentHierarchy,
	},
	Route{
		"GetReplicasetHierarchy",
		"GET",
		"/api/hierarchy/replicaset",
		apiHandlers.GetReplicasetHierarchy,
	},
	Route{
		"GetStatefulsetHierarchy",
		"GET",
		"/api/hierarchy/statefulset",
		apiHandlers.GetStatefulsetHierarchy,
	},
	Route{
		"GetPodHierarchy",
		"GET",
		"/api/hierarchy/pod",
		apiHandlers.GetPodHierarchy,
	},
	Route{
		"GetContainerHierarchy",
		"GET",
		"/api/hierarchy/container",
		apiHandlers.GetContainerHierarchy,
	},
	Route{
		"GetProcessHierarchy",
		"GET",
		"/api/hierarchy/process",
		apiHandlers.GetEmptyHierarchy,
	},
	Route{
		"GetNodeHierarchy",
		"GET",
		"/api/hierarchy/node",
		apiHandlers.GetNodeHierarchy,
	},
	Route{
		"GetPVHierarchy",
		"GET",
		"/api/hierarchy/pv",
		apiHandlers.GetPVHierarchy,
	},
	Route{
		"GetPVCHierarchy",
		"GET",
		"/api/hierarchy/pvc",
		apiHandlers.GetEmptyHierarchy,
	},
	Route{
		"GetDaemonsetHierarchy",
		"GET",
		"/api/hierarchy/daemonset",
		apiHandlers.GetDaemonsetHierarchy,
	},
	Route{
		"GetJobHierarchy",
		"GET",
		"/api/hierarchy/job",
		apiHandlers.GetJobHierarchy,
	},
	Route{
		"GetClusterMetrics",
		"GET",
		"/api/metrics",
		apiHandlers.GetClusterMetrics,
	},
	Route{
		"GetNamespaceMetrics",
		"GET",
		"/api/metrics/namespace",
		apiHandlers.GetNamespaceMetrics,
	},
	Route{
		"GetDeploymentMetrics",
		"GET",
		"/api/metrics/deployment",
		apiHandlers.GetDeploymentMetrics,
	},
	Route{
		"GetDaemonsetMetrics",
		"GET",
		"/api/metrics/daemonset",
		apiHandlers.GetDaemonsetMetrics,
	},
	Route{
		"GetJobMetrics",
		"GET",
		"/api/metrics/job",
		apiHandlers.GetJobMetrics,
	},
	Route{
		"GetStatefulsetMetrics",
		"GET",
		"/api/metrics/statefulset",
		apiHandlers.GetStatefulsetMetrics,
	},
	Route{
		"GetReplicasetMetrics",
		"GET",
		"/api/metrics/replicaset",
		apiHandlers.GetReplicasetMetrics,
	},
	Route{
		"GetNodeMetrics",
		"GET",
		"/api/metrics/node",
		apiHandlers.GetNodeMetrics,
	},
	Route{
		"GetPodMetrics",
		"GET",
		"/api/metrics/pod",
		apiHandlers.GetPodMetrics,
	},
	Route{
		"GetContainerMetrics",
		"GET",
		"/api/metrics/container",
		apiHandlers.GetContainerMetrics,
	},
	Route{
		"GetPVMetrics",
		"GET",
		"/api/metrics/pv",
		apiHandlers.GetPVMetrics,
	},
	Route{
		"GetPVCMetrics",
		"GET",
		"/api/metrics/pvc",
		apiHandlers.GetPVCMetrics,
	},
	Route{
		"GetPodDiscoveryNodes",
		"GET",
		"/api/nodes",
		apiHandlers.GetPodDiscoveryNodes,
	},
	Route{
		"GetPodDiscoveryEdges",
		"GET",
		"/api/edges",
		apiHandlers.GetPodDiscoveryEdges,
	},
	Route{
		"GetGroupsData",
		"GET",
		"/api/groups",
		apiHandlers.GetGroupsData,
	},
	Route{
		"Login",
		"POST",
		"/auth/login",
		apiHandlers.LoginUser,
	},
	Route{
		"Logout",
		"POST",
		"/auth/logout",
		apiHandlers.LogoutUser,
	},
	Route{
		"ChangePassword",
		"POST",
		"/auth/changePassword",
		apiHandlers.ChangePassword,
	},
	Route{
		"DeleteGroup",
		"POST",
		"/api/group/delete",
		apiHandlers.DeleteGroup,
	},
	Route{
		"CreateGroup",
		"POST",
		"/api/group/create",
		apiHandlers.CreateGroup,
	},
	Route{
		"SyncCluster",
		"GET",
		"/api/sync",
		apiHandlers.SyncCluster,
	},
}
