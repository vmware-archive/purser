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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

// GetHomePage is the default appD home page
func GetHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Application Discovery!")
}

// GetInventoryPods gives all the pods stored in dgraph
func GetInventoryPods(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if orphanVal, isOrphan := queryParams["orphan"]; isOrphan && orphanVal[0] == "true" {
		jsonResp, err := models.RetrievePodsInteractionsForAllPodsOrphanedTrue()
		if err != nil {
			logrus.Errorf("Unable to get response: (%v)", err)
		}
		err = json.NewEncoder(w).Encode(jsonResp)
		if err != nil {
			logrus.Errorf("Unable to encode to json: (%v)", err)
		}
	} else {
		jsonResp, err := models.RetrievePodsInteractionsForAllPodsOrphanedFalse()
		if err != nil {
			logrus.Errorf("Unable to get response: (%v)", err)
		}
		err = json.NewEncoder(w).Encode(jsonResp)
		if err != nil {
			logrus.Errorf("Unable to encode to json: (%v)", err)
		}
	}
}

// GetPodInteractions listens on /interactions/pod endpoint and returns pod interactions
func GetPodInteractions(w http.ResponseWriter, r *http.Request) {
	var pod []models.Pod
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		pod, err = models.RetrievePodsInteractionsForGivenPod(name[0])
	} else {
		if orphanVal, isOrphan := queryParams["orphan"]; isOrphan && orphanVal[0] == "true" {
			pod, err = models.RetrievePodsInteractionsForAllPodsOrphanedTrue()
		} else {
			pod, err = models.RetrievePodsInteractionsForAllPodsOrphanedFalse()
		}
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}
	err = json.NewEncoder(w).Encode(pod)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetNamespaceHierarchy listens on /hierarchy/namespace endpoint and returns all namespace and their children up to 2 levels
func GetNamespaceHierarchy(w http.ResponseWriter, r *http.Request) {
	var namespace []byte
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		namespace, err = models.RetrieveNamespace(name[0])
	} else {
		namespace, err = models.RetrieveAllNamespaces()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	_, err = w.Write(namespace)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetDeploymentHierarchy listens on /hierarchy/deployment endpoint and returns all deployments and their children up to 2 levels
func GetDeploymentHierarchy(w http.ResponseWriter, r *http.Request) {
	var deployment []byte
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		deployment, err = models.RetrieveDeployment(name[0])
	} else {
		deployment, err = models.RetrieveAllDeployments()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	_, err = w.Write(deployment)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetReplicasetHierarchy listens on /hierarchy/replicaset endpoint and returns all replicasets and their children up to 2 levels
func GetReplicasetHierarchy(w http.ResponseWriter, r *http.Request) {
	var replicaset []byte
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		replicaset, err = models.RetrieveReplicaset(name[0])
	} else {
		replicaset, err = models.RetrieveAllReplicasets()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	_, err = w.Write(replicaset)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetStatefulsetHierarchy listens on /hierarchy/statefulset endpoint and returns all statefulsets and their children up to 2 levels
func GetStatefulsetHierarchy(w http.ResponseWriter, r *http.Request) {
	var statefulset []byte
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		statefulset, err = models.RetrieveStatefulset(name[0])
	} else {
		statefulset, err = models.RetrieveAllStatefulsets()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	_, err = w.Write(statefulset)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetPodHierarchy listens on /hierarchy/pod endpoint and returns all pods and their children up to 2 levels
func GetPodHierarchy(w http.ResponseWriter, r *http.Request) {
	var pod []byte
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		pod, err = models.RetrievePod(name[0])
	} else {
		pod, err = models.RetrieveAllPods()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	_, err = w.Write(pod)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetContainerHierarchy listens on /hierarchy/container endpoint and returns all containers along with process in them.
func GetContainerHierarchy(w http.ResponseWriter, r *http.Request) {
	var container []byte
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		container, err = models.RetrieveContainer(name[0])
	} else {
		container, err = models.RetrieveAllContainers()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	_, err = w.Write(container)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetProcessHierarchy listens on /hierarchy/process endpoint and returns all processes
func GetProcessHierarchy(w http.ResponseWriter, r *http.Request) {
	var process []byte
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		process, err = models.RetrieveProcess(name[0])
	} else {
		process, err = models.RetrieveAllProcess()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	_, err = w.Write(process)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetNodeHierarchy listens on /hierarchy/node endpoint and returns all nodes with their children up to 2 levels
func GetNodeHierarchy(w http.ResponseWriter, r *http.Request) {
	var node []byte
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		node, err = models.RetrieveNode(name[0])
	} else {
		node, err = models.RetrieveAllNodes()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	_, err = w.Write(node)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetDeamonsetHierarchy listens on /hierarchy/daemonset endpoint and returns all daemonsets and their children upto 2 levels
func GetDeamonsetHierarchy(w http.ResponseWriter, r *http.Request) {
	var daemonset []byte
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		daemonset, err = models.RetrieveDaemonset(name[0])
	} else {
		daemonset, err = models.RetrieveAllDaemonsets()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	_, err = w.Write(daemonset)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetJobHierarchy listens on /hierarchy/job endpoint and returns all jobs and their children upto 2 levels
func GetJobHierarchy(w http.ResponseWriter, r *http.Request) {
	var job []byte
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		job, err = models.RetrieveJob(name[0])
	} else {
		job, err = models.RetrieveAllJobs()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	_, err = w.Write(job)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

func GetClusterWithMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if view, isView := queryParams["view"]; isView && view[0] == "physical" {
		cluster, err := models.RetrieveAllNodesWithMetrics()
		if err != nil {
			logrus.Errorf("Unable to get response: (%v)", err)
		}

		err = json.NewEncoder(w).Encode(cluster)
		if err != nil {
			logrus.Errorf("Unable to encode to json: (%v)", err)
		}
	} else {
		cluster, err := models.RetrieveAllNamespacesWithMetrics()
		if err != nil {
			logrus.Errorf("Unable to get response: (%v)", err)
		}

		err = json.NewEncoder(w).Encode(cluster)
		if err != nil {
			logrus.Errorf("Unable to encode to json: (%v)", err)
		}
	}
}

// GetNamespaceWithMetrics ...
func GetNamespaceWithMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		namespace, err := models.RetrieveNamespaceWithMetrics(name[0])
		if err != nil {
			logrus.Errorf("Unable to get response: (%v)", err)
		}

		err = json.NewEncoder(w).Encode(namespace)
		if err != nil {
			logrus.Errorf("Unable to encode to json: (%v)", err)
		}
	} else {
		GetClusterWithMetrics(w, r)
	}

}

// GetDeploymentWithMetrics ...
func GetDeploymentWithMetrics(w http.ResponseWriter, r *http.Request) {
	var deployment models.DeploymentsWithMetrics
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		deployment, err = models.RetrieveDeploymentWithMetrics(name[0])
	} else {
		deployment, err = models.RetrieveAllDeploymentsWithMetrics()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	err = json.NewEncoder(w).Encode(deployment)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetDaemonsetWithMetrics ...
func GetDaemonsetWithMetrics(w http.ResponseWriter, r *http.Request) {
	var daemonset models.DaemonsetsWithMetrics
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		daemonset, err = models.RetrieveDaemonsetWithMetrics(name[0])
	} else {
		daemonset, err = models.RetrieveAllDaemonsetsWithMetrics()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	err = json.NewEncoder(w).Encode(daemonset)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetStatefulsetWithMetrics ...
func GetStatefulsetWithMetrics(w http.ResponseWriter, r *http.Request) {
	var statefulset models.StatefulsetsWithMetrics
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		statefulset, err = models.RetrieveStatefulsetWithMetrics(name[0])
	} else {
		statefulset, err = models.RetrieveAllStatefulsetsWithMetrics()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	err = json.NewEncoder(w).Encode(statefulset)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetReplicasetWithMetrics ...
func GetReplicasetWithMetrics(w http.ResponseWriter, r *http.Request) {
	var replicaset models.ReplicasetsWithMetrics
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		replicaset, err = models.RetrieveReplicasetWithMetrics(name[0])
	} else {
		replicaset, err = models.RetrieveAllReplicasetsWithMetrics()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	err = json.NewEncoder(w).Encode(replicaset)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetJobWithMetrics ...
func GetJobWithMetrics(w http.ResponseWriter, r *http.Request) {
	var job models.JobsWithMetrics
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		job, err = models.RetrieveJobWithMetrics(name[0])
	} else {
		job, err = models.RetrieveAllJobsWithMetrics()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	err = json.NewEncoder(w).Encode(job)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetNodeWithMetrics ...
func GetNodeWithMetrics(w http.ResponseWriter, r *http.Request) {
	var node models.NodesWithMetrics
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		node, err = models.RetrieveNodeWithMetrics(name[0])
		if err != nil {
			logrus.Errorf("Unable to get response: (%v)", err)
		}

		err = json.NewEncoder(w).Encode(node)
		if err != nil {
			logrus.Errorf("Unable to encode to json: (%v)", err)
		}
	}
}

// GetPodWithMetrics ...
func GetPodWithMetrics(w http.ResponseWriter, r *http.Request) {
	var pod models.PodsWithMetrics
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		pod, err = models.RetrievePodWithMetrics(name[0])
	} else {
		pod, err = models.RetrieveAllPodsWithMetrics()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	err = json.NewEncoder(w).Encode(pod)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetContainerWithMetrics ...
func GetContainerWithMetrics(w http.ResponseWriter, r *http.Request) {
	var container models.ContainersWithMetrics
	var err error

	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		container, err = models.RetrieveContainerWithMetrics(name[0])
	} else {
		container, err = models.RetrieveAllContainersWithMetrics()
	}
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}

	err = json.NewEncoder(w).Encode(container)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

func addHeaders(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)
}