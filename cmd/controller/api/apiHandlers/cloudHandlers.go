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

package apiHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

// GetCloudRegionList listens on /api/clouds/regions endpoint
func GetCloudRegionList(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	regionsData := ""
	encodeAndWrite(w, regionsData)
}

// CompareCloud listens on /api/clouds/compare endpoint
func CompareCloud(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w, r)
	queryParams := r.URL.Query()
	logrus.Debugf("Query params: (%v)", queryParams)

	regionData, err := convertRequestBodyToJSON(r)
	if err != nil {
		logrus.Errorf("unable to parse request as either JSON or YAML, err: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var region []models.CloudRegion

	if jsonErr := json.Unmarshal(regionData, &region); jsonErr != nil {
		logrus.Errorf("unable to parse object as group, err: %v", jsonErr)
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
		return
	}
	logrus.Printf("region  %#v ", region)
	nodeCost := models.GetCost(region)
	encodeAndWrite(w, nodeCost)
}
