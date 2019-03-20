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
	"io"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
	"github.com/vmware/purser/pkg/controller/discovery/generator"
)

// Credentials structure
type Credentials struct {
	Password    string `json:"password"`
	Username    string `json:"username"`
	NewPassword string `json:"newPassword"`
}

const cookieName = "session-token-purser"

var store = sessions.NewCookieStore([]byte(cookieKey))

// LoginUser listens on /auth/login endpoint
func LoginUser(w http.ResponseWriter, r *http.Request) {
	addAccessControlHeaders(&w, r)
	var cred Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !query.CheckLogin(cred.Username, cred.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, err := store.Get(r, cookieName)
	if err != nil {
		logrus.Errorf("unable to get session from cookie store, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.Values["authenticated"] = true
	saveSession(session, w, r)
}

// LogoutUser listens on /auth/logout endpoint
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	addAccessControlHeaders(&w, r)
	session, err := store.Get(r, cookieName)
	if err != nil {
		logrus.Errorf("unable to get session from cookie store, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.Values["authenticated"] = false
	saveSession(session, w, r)
}

func saveSession(session *sessions.Session, w http.ResponseWriter, r *http.Request) {
	err := session.Save(r, w)
	if err != nil {
		logrus.Errorf("unable to get session from cookie store, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func isUserAuthenticated(r *http.Request) bool {
	session, err := store.Get(r, cookieName)
	if err != nil {
		logrus.Errorf("unable to get session from cookie store, err: %v", err)
		return false
	}
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false
	}
	return true
}

// ChangePassword listens on /auth/changePassword endpoint
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	addAccessControlHeaders(&w, r)
	var cred Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !query.UpdateLogin(cred.Username, cred.Password, cred.NewPassword) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// GetHomePage is the default api home page
func GetHomePage(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	_, err := fmt.Fprintf(w, "Welcome to the Purser!")
	if err != nil {
		logrus.Errorf("Unable to write welcome message to Homepage: (%v)", err)
	}
}

// GetPodInteractions listens on /interactions/pod endpoint and returns pod interactions
func GetPodInteractions(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
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
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
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
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check:       query.NamespaceCheck,
			Type:        query.NamespaceType,
			Name:        name[0],
			ChildFilter: query.NamespaceChildFilter,
		}
		jsonData = resourceQuery.RetrieveResourceHierarchy()
	} else {
		jsonData = query.RetrieveClusterHierarchy(query.Logical)
	}
	encodeAndWrite(w, jsonData)
}

// GetDeploymentHierarchy listens on /hierarchy/deployment endpoint and returns all children of deployment
func GetDeploymentHierarchy(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check:       query.DeploymentCheck,
			Type:        query.DeploymentType,
			Name:        name[0],
			ChildFilter: query.IsReplicasetFilter,
		}
		jsonData = resourceQuery.RetrieveResourceHierarchy()
	} else {
		logrus.Errorf("wrong type of query for deployment, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetReplicasetHierarchy listens on /hierarchy/replicaset endpoint and returns all children of replicaset
func GetReplicasetHierarchy(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check:       query.ReplicasetCheck,
			Type:        query.ReplicasetType,
			Name:        name[0],
			ChildFilter: query.IsPodFilter,
		}
		jsonData = resourceQuery.RetrieveResourceHierarchy()
	} else {
		logrus.Errorf("wrong type of query for replicaset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetStatefulsetHierarchy listens on /hierarchy/statefulset endpoint and returns all children of statefulset
func GetStatefulsetHierarchy(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check:       query.StatefulsetCheck,
			Type:        query.StatefulsetType,
			Name:        name[0],
			ChildFilter: query.IsPodFilter,
		}
		jsonData = resourceQuery.RetrieveResourceHierarchy()
	} else {
		logrus.Errorf("wrong type of query for statefulset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetPodHierarchy listens on /hierarchy/pod endpoint and returns all children of pod
func GetPodHierarchy(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check:       query.PodCheck,
			Type:        query.PodType,
			Name:        name[0],
			ChildFilter: query.IsContainerFilter,
		}
		jsonData = resourceQuery.RetrieveResourceHierarchy()
	} else {
		logrus.Errorf("wrong type of query for pod, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetContainerHierarchy listens on /hierarchy/container endpoint and returns all children of container
func GetContainerHierarchy(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check:       query.ContainerCheck,
			Type:        query.ContainerType,
			Name:        name[0],
			ChildFilter: query.IsProcFilter,
		}
		jsonData = resourceQuery.RetrieveResourceHierarchy()
	} else {
		logrus.Errorf("wrong type of query for container, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetEmptyHierarchy listens on /hierarchy/process and /hierarchy/pvc endpoint and returns empty data
func GetEmptyHierarchy(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	encodeAndWrite(w, jsonData)
}

// GetNodeHierarchy listens on /hierarchy/node endpoint and returns all children of node
func GetNodeHierarchy(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check:       query.NodeCheck,
			Type:        query.NodeType,
			Name:        name[0],
			ChildFilter: query.IsPodFilter,
		}
		jsonData = resourceQuery.RetrieveResourceHierarchy()
	} else {
		logrus.Errorf("wrong type of query for node, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetPVHierarchy listens on /hierarchy/pv endpoint and returns all children of PV
func GetPVHierarchy(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check:       query.PVCheck,
			Type:        query.PVType,
			Name:        name[0],
			ChildFilter: query.IsPVCFilter,
		}
		jsonData = resourceQuery.RetrieveResourceHierarchy()
	} else {
		logrus.Errorf("wrong type of query for PV, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetDaemonsetHierarchy listens on /hierarchy/daemonset endpoint and returns all children of Daemonset
func GetDaemonsetHierarchy(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check:       query.DaemonsetCheck,
			Type:        query.DaemonsetType,
			Name:        name[0],
			ChildFilter: query.IsPodFilter,
		}
		jsonData = resourceQuery.RetrieveResourceHierarchy()
	} else {
		logrus.Errorf("wrong type of query for Daemonset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetJobHierarchy listens on /hierarchy/job endpoint and returns all children of Job
func GetJobHierarchy(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check:       query.JobCheck,
			Type:        query.JobType,
			Name:        name[0],
			ChildFilter: query.IsPodFilter,
		}
		jsonData = resourceQuery.RetrieveResourceHierarchy()
	} else {
		logrus.Errorf("wrong type of query for Job, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetClusterMetrics listens on /metrics endpoint with option for view(physical or logical)
func GetClusterMetrics(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
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
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check: query.NamespaceCheck,
			Type:  query.NamespaceType,
			Name:  name[0],
		}
		jsonData = resourceQuery.RetrieveResourceMetrics()
	} else {
		jsonData = query.RetrieveClusterMetrics(query.Logical)
	}
	encodeAndWrite(w, jsonData)
}

// GetDeploymentMetrics listens on /metrics/deployment
func GetDeploymentMetrics(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check: query.DeploymentCheck,
			Type:  query.DeploymentType,
			Name:  name[0],
		}
		jsonData = resourceQuery.RetrieveResourceMetrics()
	} else {
		logrus.Errorf("wrong type of query for deployment, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetDaemonsetMetrics listens on /metrics/daemonset
func GetDaemonsetMetrics(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check: query.DaemonsetCheck,
			Type:  query.DaemonsetType,
			Name:  name[0],
		}
		jsonData = resourceQuery.RetrieveResourceMetrics()
	} else {
		logrus.Errorf("wrong type of query for daemonset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetJobMetrics listens on /metrics/job
func GetJobMetrics(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check: query.JobCheck,
			Type:  query.JobType,
			Name:  name[0],
		}
		jsonData = resourceQuery.RetrieveResourceMetrics()
	} else {
		logrus.Errorf("wrong type of query for job, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetStatefulsetMetrics listens on /metrics/statefulset
func GetStatefulsetMetrics(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check: query.StatefulsetCheck,
			Type:  query.StatefulsetType,
			Name:  name[0],
		}
		jsonData = resourceQuery.RetrieveResourceMetrics()
	} else {
		logrus.Errorf("wrong type of query for statefulset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetReplicasetMetrics listens on /metrics/replicaset
func GetReplicasetMetrics(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check: query.ReplicasetCheck,
			Type:  query.ReplicasetType,
			Name:  name[0],
		}
		jsonData = resourceQuery.RetrieveResourceMetrics()
	} else {
		logrus.Errorf("wrong type of query for statefulset, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetNodeMetrics listens on /metrics/node
func GetNodeMetrics(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check: query.NodeCheck,
			Type:  query.NodeType,
			Name:  name[0],
		}
		jsonData = resourceQuery.RetrieveResourceMetrics()
	} else {
		logrus.Errorf("wrong type of query for node, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetPodMetrics listens on /metrics/pod
func GetPodMetrics(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check: query.PodCheck,
			Type:  query.PodType,
			Name:  name[0],
		}
		jsonData = resourceQuery.RetrieveResourceMetrics()
	} else {
		logrus.Errorf("wrong type of query for pod, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetContainerMetrics listens on /metrics/container
func GetContainerMetrics(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check: query.ContainerCheck,
			Type:  query.ContainerType,
			Name:  name[0],
		}
		jsonData = resourceQuery.RetrieveResourceMetrics()
	} else {
		logrus.Errorf("wrong type of query for container, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetPVMetrics listens on /metrics/pv
func GetPVMetrics(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check: query.PVCheck,
			Type:  query.PVType,
			Name:  name[0],
		}
		jsonData = resourceQuery.RetrieveResourceMetrics()
	} else {
		logrus.Errorf("wrong type of query for PV, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetPVCMetrics listens on /metrics/pv
func GetPVCMetrics(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	var jsonData query.JSONDataWrapper
	if name, isName := queryParams[query.Name]; isName {
		resourceQuery := query.Resource{
			Check: query.PVCCheck,
			Type:  query.PVCType,
			Name:  name[0],
		}
		jsonData = resourceQuery.RetrieveResourceMetrics()
	} else {
		logrus.Errorf("wrong type of query for PVC, no name is given")
	}
	encodeAndWrite(w, jsonData)
}

// GetPodDiscoveryNodes listens on /discovery/pod/nodes endpoint
func GetPodDiscoveryNodes(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var pods []models.Pod
	var err error

	addHeaders(&w, r)
	pods, err = query.RetrievePodsInteractionsForAllLivePodsWithCount()
	generator.GeneratePodNodesAndEdges(pods)
	if err != nil {
		logrus.Errorf("Unable to get response: (%v)", err)
	}
	nodes := generator.GetGraphNodes()
	if nodes != nil {
		logrus.Infof("No nodes found")
		return
	}
	err = json.NewEncoder(w).Encode(nodes)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetPodDiscoveryEdges listens on /discovery/pod/edges endpoint
func GetPodDiscoveryEdges(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var err error
	addHeaders(&w, r)

	edges := generator.GetGraphEdges()
	if edges == nil {
		logrus.Infof("No edges found")
		return
	}
	err = json.NewEncoder(w).Encode(edges)
	if err != nil {
		logrus.Errorf("Unable to encode to json: (%v)", err)
	}
}

// GetGroupsData listens on /groups endpoint
func GetGroupsData(w http.ResponseWriter, r *http.Request) {
	if !isUserAuthenticated(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	addHeaders(&w, r)

	groupsData, err := query.RetrieveGroupsData()
	if err != nil {
		logrus.Errorf("unable to retrieve groups data from dgraph, %v", err)
	} else {
		encodeAndWrite(w, groupsData)
	}
}

func addHeaders(w *http.ResponseWriter, r *http.Request) {
	addAccessControlHeaders(w, r)
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).WriteHeader(http.StatusOK)
}

func addAccessControlHeaders(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
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
