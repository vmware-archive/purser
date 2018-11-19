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
	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
	"io"
	"net/http"
)

// GetHomePage is the default api home page
func GetHomePage(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Welcome to the Purser!")
	if err != nil {
		logrus.Errorf("Unable to write welcome message to Homepage: (%v)", err)
	}
}

// GetPodInteractions listens on /interactions/pod endpoint and returns pod interactions
func GetPodInteractions(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonResp []byte
	if name, isName := queryParams[query.Name]; isName {
		jsonResp = query.RetrievePodsInteractions(name[0], false)
	} else {
		if orphanVal, isOrphan := queryParams[query.Orphan]; isOrphan && orphanVal[0] == query.False {
			jsonResp = query.RetrievePodsInteractions(query.All, false)
		} else {
			jsonResp = query.RetrievePodsInteractions(query.All, true)
		}
	}
	writeBytes(w, jsonResp)
}

// GetClusterHierarchy listens on /hierarchy endpoint and returns all namespaces(or nodes and PV) in the cluster
func GetClusterHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if view, isView := queryParams[query.View]; isView && view[0] == query.Physical {
		jsonData = query.RetrieveClusterHierarchy(query.Physical)
	} else {
		jsonData = query.RetrieveClusterHierarchy(query.Logical)
	}
	encodeAndWrite(w, jsonData)
}

// GetNamespaceHierarchy listens on /hierarchy/namespace endpoint and returns all children of namespace
func GetNamespaceHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveNamespaceHierarchy(name[0])
	} else {
		jsonData = query.RetrieveNamespaceHierarchy(query.All)
	}
	encodeAndWrite(w, jsonData)
}

// GetDeploymentHierarchy listens on /hierarchy/deployment endpoint and returns all children of deployment
func GetDeploymentHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveDeploymentHierarchy(name[0])
	} else {
		logrus.Errorf("wrong type of query for deployment, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetReplicasetHierarchy listens on /hierarchy/replicaset endpoint and returns all children of replicaset
func GetReplicasetHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveReplicasetHierarchy(name[0])
	} else {
		logrus.Errorf("wrong type of query for replicaset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetStatefulsetHierarchy listens on /hierarchy/statefulset endpoint and returns all children of statefulset
func GetStatefulsetHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveStatefulsetHierarchy(name[0])
	} else {
		logrus.Errorf("wrong type of query for statefulset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetPodHierarchy listens on /hierarchy/pod endpoint and returns all children of pod
func GetPodHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrievePodHierarchy(name[0])
	} else {
		logrus.Errorf("wrong type of query for pod, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetContainerHierarchy listens on /hierarchy/container endpoint and returns all children of container
func GetContainerHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveContainerHierarchy(name[0])
	} else {
		logrus.Errorf("wrong type of query for container, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetEmptyHierarchy listens on /hierarchy/process and /hierarchy/pvc endpoint and returns empty data
func GetEmptyHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	encodeAndWrite(w, jsonData)
}

// GetNodeHierarchy listens on /hierarchy/node endpoint and returns all children of node
func GetNodeHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveNodeHierarchy(name[0])
	} else {
		logrus.Errorf("wrong type of query for node, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetPVHierarchy listens on /hierarchy/pv endpoint and returns all children of PV
func GetPVHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrievePVHierarchy(name[0])
	} else {
		logrus.Errorf("wrong type of query for PV, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetDaemonsetHierarchy listens on /hierarchy/daemonset endpoint and returns all children of Daemonset
func GetDaemonsetHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveDaemonsetHierarchy(name[0])
	} else {
		logrus.Errorf("wrong type of query for Daemonset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetJobHierarchy listens on /hierarchy/job endpoint and returns all children of Job
func GetJobHierarchy(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveJobHierarchy(name[0])
	} else {
		logrus.Errorf("wrong type of query for Job, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetClusterMetrics listens on /metrics endpoint with option for view(physical or logical)
func GetClusterMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if view, isView := queryParams[query.View]; isView && view[0] == query.Physical {
		jsonData = query.RetrieveClusterMetrics(query.Physical)
	} else {
		jsonData = query.RetrieveClusterMetrics(query.Logical)
	}
	encodeAndWrite(w, jsonData)
}

// GetNamespaceMetrics listens on /metrics/namespace
func GetNamespaceMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveNamespaceMetrics(name[0])
	} else {
		logrus.Errorf("wrong type of query for namespace, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetDeploymentMetrics listens on /metrics/deployment
func GetDeploymentMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveDeploymentMetrics(name[0])
	} else {
		logrus.Errorf("wrong type of query for deployment, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetDaemonsetMetrics listens on /metrics/daemonset
func GetDaemonsetMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveDaemonsetMetrics(name[0])
	} else {
		logrus.Errorf("wrong type of query for daemonset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetJobMetrics listens on /metrics/job
func GetJobMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveJobMetrics(name[0])
	} else {
		logrus.Errorf("wrong type of query for job, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetStatefulsetMetrics listens on /metrics/statefulset
func GetStatefulsetMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveStatefulsetMetrics(name[0])
	} else {
		logrus.Errorf("wrong type of query for statefulset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetReplicasetMetrics listens on /metrics/replicaset
func GetReplicasetMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveReplicasetMetrics(name[0])
	} else {
		logrus.Errorf("wrong type of query for statefulset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetNodeMetrics listens on /metrics/node
func GetNodeMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveNodeMetrics(name[0])
	} else {
		logrus.Errorf("wrong type of query for node, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetPodMetrics listens on /metrics/pod
func GetPodMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrievePodMetrics(name[0])
	} else {
		logrus.Errorf("wrong type of query for pod, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetContainerMetrics listens on /metrics/container
func GetContainerMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrieveContainerMetrics(name[0])
	} else {
		logrus.Errorf("wrong type of query for container, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetPVMetrics listens on /metrics/pv
func GetPVMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrievePVMetrics(name[0])
	} else {
		logrus.Errorf("wrong type of query for PV, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetPVCMetrics listens on /metrics/pv
func GetPVCMetrics(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		jsonData = query.RetrievePVCMetrics(name[0])
	} else {
		logrus.Errorf("wrong type of query for PVC, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

func addHeaders(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).WriteHeader(http.StatusOK)
}

func writeBytes(w io.Writer, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

func encodeAndWrite(w io.Writer, obj interface{}) {
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}
