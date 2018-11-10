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