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
