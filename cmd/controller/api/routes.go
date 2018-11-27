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
		"/",
		GetHomePage,
	},
	Route{
		"GetPodInteractions",
		"GET",
		"/interactions/pod",
		GetPodInteractions,
	},
	Route{
		"GetClusterHierarchy",
		"GET",
		"/hierarchy",
		GetClusterHierarchy,
	},
	Route{
		"GetNamespaceHierarchy",
		"GET",
		"/hierarchy/namespace",
		GetNamespaceHierarchy,
	},
	Route{
		"GetDeploymentHierarchy",
		"GET",
		"/hierarchy/deployment",
		GetDeploymentHierarchy,
	},
	Route{
		"GetReplicasetHierarchy",
		"GET",
		"/hierarchy/replicaset",
		GetReplicasetHierarchy,
	},
	Route{
		"GetStatefulsetHierarchy",
		"GET",
		"/hierarchy/statefulset",
		GetStatefulsetHierarchy,
	},
	Route{
		"GetPodHierarchy",
		"GET",
		"/hierarchy/pod",
		GetPodHierarchy,
	},
	Route{
		"GetContainerHierarchy",
		"GET",
		"/hierarchy/container",
		GetContainerHierarchy,
	},
	Route{
		"GetProcessHierarchy",
		"GET",
		"/hierarchy/process",
		GetEmptyHierarchy,
	},
	Route{
		"GetNodeHierarchy",
		"GET",
		"/hierarchy/node",
		GetNodeHierarchy,
	},
	Route{
		"GetPVHierarchy",
		"GET",
		"/hierarchy/pv",
		GetPVHierarchy,
	},
	Route{
		"GetPVCHierarchy",
		"GET",
		"/hierarchy/pvc",
		GetEmptyHierarchy,
	},
	Route{
		"GetDaemonsetHierarchy",
		"GET",
		"/hierarchy/daemonset",
		GetDaemonsetHierarchy,
	},
	Route{
		"GetJobHierarchy",
		"GET",
		"/hierarchy/job",
		GetJobHierarchy,
	},
	Route{
		"GetClusterMetrics",
		"GET",
		"/metrics",
		GetClusterMetrics,
	},
	Route{
		"GetNamespaceMetrics",
		"GET",
		"/metrics/namespace",
		GetNamespaceMetrics,
	},
	Route{
		"GetDeploymentMetrics",
		"GET",
		"/metrics/deployment",
		GetDeploymentMetrics,
	},
	Route{
		"GetDaemonsetMetrics",
		"GET",
		"/metrics/daemonset",
		GetDaemonsetMetrics,
	},
	Route{
		"GetJobMetrics",
		"GET",
		"/metrics/job",
		GetJobMetrics,
	},
	Route{
		"GetStatefulsetMetrics",
		"GET",
		"/metrics/statefulset",
		GetStatefulsetMetrics,
	},
	Route{
		"GetReplicasetMetrics",
		"GET",
		"/metrics/replicaset",
		GetReplicasetMetrics,
	},
	Route{
		"GetNodeMetrics",
		"GET",
		"/metrics/node",
		GetNodeMetrics,
	},
	Route{
		"GetPodMetrics",
		"GET",
		"/metrics/pod",
		GetPodMetrics,
	},
	Route{
		"GetContainerMetrics",
		"GET",
		"/metrics/container",
		GetContainerMetrics,
	},
	Route{
		"GetPVMetrics",
		"GET",
		"/metrics/pv",
		GetPVMetrics,
	},
	Route{
		"GetPVCMetrics",
		"GET",
		"/metrics/pvc",
		GetPVCMetrics,
	},
	Route{
		"GetPodDiscoveryNodes",
		"GET",
		"/nodes",
		GetPodDiscoveryNodes,
	},
	Route{
		"GetPodDiscoveryEdges",
		"GET",
		"/edges",
		GetPodDiscoveryEdges,
	},
}
