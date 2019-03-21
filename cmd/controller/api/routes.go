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
		GetHomePage,
	},
	Route{
		"GetPodInteractions",
		"GET",
		"/api/interactions/pod",
		GetPodInteractions,
	},
	Route{
		"GetClusterHierarchy",
		"GET",
		"/api/hierarchy",
		GetClusterHierarchy,
	},
	Route{
		"GetNamespaceHierarchy",
		"GET",
		"/api/hierarchy/namespace",
		GetNamespaceHierarchy,
	},
	Route{
		"GetDeploymentHierarchy",
		"GET",
		"/api/hierarchy/deployment",
		GetDeploymentHierarchy,
	},
	Route{
		"GetReplicasetHierarchy",
		"GET",
		"/api/hierarchy/replicaset",
		GetReplicasetHierarchy,
	},
	Route{
		"GetStatefulsetHierarchy",
		"GET",
		"/api/hierarchy/statefulset",
		GetStatefulsetHierarchy,
	},
	Route{
		"GetPodHierarchy",
		"GET",
		"/api/hierarchy/pod",
		GetPodHierarchy,
	},
	Route{
		"GetContainerHierarchy",
		"GET",
		"/api/hierarchy/container",
		GetContainerHierarchy,
	},
	Route{
		"GetProcessHierarchy",
		"GET",
		"/api/hierarchy/process",
		GetEmptyHierarchy,
	},
	Route{
		"GetNodeHierarchy",
		"GET",
		"/api/hierarchy/node",
		GetNodeHierarchy,
	},
	Route{
		"GetPVHierarchy",
		"GET",
		"/api/hierarchy/pv",
		GetPVHierarchy,
	},
	Route{
		"GetPVCHierarchy",
		"GET",
		"/api/hierarchy/pvc",
		GetEmptyHierarchy,
	},
	Route{
		"GetDaemonsetHierarchy",
		"GET",
		"/api/hierarchy/daemonset",
		GetDaemonsetHierarchy,
	},
	Route{
		"GetJobHierarchy",
		"GET",
		"/api/hierarchy/job",
		GetJobHierarchy,
	},
	Route{
		"GetClusterMetrics",
		"GET",
		"/api/metrics",
		GetClusterMetrics,
	},
	Route{
		"GetNamespaceMetrics",
		"GET",
		"/api/metrics/namespace",
		GetNamespaceMetrics,
	},
	Route{
		"GetDeploymentMetrics",
		"GET",
		"/api/metrics/deployment",
		GetDeploymentMetrics,
	},
	Route{
		"GetDaemonsetMetrics",
		"GET",
		"/api/metrics/daemonset",
		GetDaemonsetMetrics,
	},
	Route{
		"GetJobMetrics",
		"GET",
		"/api/metrics/job",
		GetJobMetrics,
	},
	Route{
		"GetStatefulsetMetrics",
		"GET",
		"/api/metrics/statefulset",
		GetStatefulsetMetrics,
	},
	Route{
		"GetReplicasetMetrics",
		"GET",
		"/api/metrics/replicaset",
		GetReplicasetMetrics,
	},
	Route{
		"GetNodeMetrics",
		"GET",
		"/api/metrics/node",
		GetNodeMetrics,
	},
	Route{
		"GetPodMetrics",
		"GET",
		"/api/metrics/pod",
		GetPodMetrics,
	},
	Route{
		"GetContainerMetrics",
		"GET",
		"/api/metrics/container",
		GetContainerMetrics,
	},
	Route{
		"GetPVMetrics",
		"GET",
		"/api/metrics/pv",
		GetPVMetrics,
	},
	Route{
		"GetPVCMetrics",
		"GET",
		"/api/metrics/pvc",
		GetPVCMetrics,
	},
	Route{
		"GetPodDiscoveryNodes",
		"GET",
		"/api/nodes",
		GetPodDiscoveryNodes,
	},
	Route{
		"GetPodDiscoveryEdges",
		"GET",
		"/api/edges",
		GetPodDiscoveryEdges,
	},
	Route{
		"GetGroupsData",
		"GET",
		"/api/groups",
		GetGroupsData,
	},
	Route{
		"Login",
		"POST",
		"/auth/login",
		LoginUser,
	},
	Route{
		"Logout",
		"POST",
		"/auth/logout",
		LogoutUser,
	},
	Route{
		"ChangePassword",
		"POST",
		"/auth/changePassword",
		ChangePassword,
	},
}
