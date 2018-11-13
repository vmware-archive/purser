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
		"GetInventoryPods",
		"GET",
		"/inventory/pod",
		GetInventoryPods,
	},
	Route{
		"GetPodInteractions",
		"GET",
		"/interactions/pod",
		GetPodInteractions,
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
		GetProcessHierarchy,
	},
	Route{
		"GetNodeHierarchy",
		"GET",
		"/hierarchy/node",
		GetNodeHierarchy,
	},
	Route{
		"GetDaemonsetHierarchy",
		"GET",
		"/hierarchy/daemonset",
		GetDeamonsetHierarchy,
	},
	Route{
		"GetJobHierarchy",
		"GET",
		"/hierarchy/job",
		GetJobHierarchy,
	},
	Route{
		"GetNamespaceWithMetrics",
		"GET",
		"/metrics/namespace",
		GetNamespaceWithMetrics,
	},
	Route{
		"GetDeploymentWithMetrics",
		"GET",
		"/metrics/deployment",
		GetDeploymentWithMetrics,
	},
	Route{
		"GetDaemonsetWithMetrics",
		"GET",
		"/metrics/daemon",
		GetDaemonsetWithMetrics,
	},
	Route{
		"GetStatefulsetWithMetrics",
		"GET",
		"/metrics/statefulset",
		GetStatefulsetWithMetrics,
	},
	Route{
		"GetReplicasetWithMetrics",
		"GET",
		"/metrics/replicaset",
		GetReplicasetWithMetrics,
	},
	Route{
		"GetDaemonsetWithMetrics",
		"GET",
		"/metrics/daemonset",
		GetDaemonsetWithMetrics,
	},
	Route{
		"GetJobWithMetrics",
		"GET",
		"/metrics/job",
		GetJobWithMetrics,
	},
	Route{
		"GetNodeWithMetrics",
		"GET",
		"/metrics/node",
		GetNodeWithMetrics,
	},
	Route{
		"GetPodWithMetrics",
		"GET",
		"/metrics/pod",
		GetPodWithMetrics,
	},
	Route{
		"GetContainerWithMetrics",
		"GET",
		"/metrics/container",
		GetContainerWithMetrics,
	},
}