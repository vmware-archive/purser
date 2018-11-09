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
	addHeaders(w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)
	if name, isName := queryParams["name"]; isName {
		jsonResp, err := models.RetrievePodsInteractionsForGivenPod(name[0])
		if err != nil {
			logrus.Errorf("Unable to get response: (%v)", err)
		}
		err = json.NewEncoder(w).Encode(jsonResp)
		if err != nil {
			logrus.Errorf("Unable to encode to json: (%v)", err)
		}
	} else {
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