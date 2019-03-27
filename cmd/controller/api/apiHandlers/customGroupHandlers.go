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
	"github.com/Sirupsen/logrus"
	group_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	"github.com/vmware/purser/pkg/client/clientset/typed/groups/v1"
	"github.com/vmware/purser/pkg/controller"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
	"github.com/vmware/purser/pkg/controller/eventprocessor"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

var groupClient *v1.GroupClient

// GetGroupsData listens on /api/groups endpoint
func GetGroupsData(w http.ResponseWriter, r *http.Request) {
	if isUserAuthenticated(w, r) {
		addHeaders(&w, r)

		groupsData, err := query.RetrieveGroupsData()
		if err != nil {
			logrus.Errorf("unable to retrieve groups data from dgraph, %v", err)
		} else {
			encodeAndWrite(w, groupsData)
		}
	}
}

// DeleteGroup listens on /api/group/delete
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	if isUserAuthenticated(w, r) {
		addAccessControlHeaders(&w, r)
		queryParams := r.URL.Query()
		logrus.Debugf("Query params: (%v)", queryParams)
		var err error
		if name, isName := queryParams[query.Name]; isName {
			err = getGroupClient().Delete(name[0], &meta_v1.DeleteOptions{})
			if err == nil {
				w.WriteHeader(http.StatusOK)
				models.DeleteGroup(name[0])
				return
			}
		}
		logrus.Errorf("unable to delete: query params: %v, err: %v", queryParams, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// CreateGroup listens on /api/group/create
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	if isUserAuthenticated(w, r) {
		addAccessControlHeaders(&w, r)
		groupData, err := convertRequestBodyToJSON(r)
		if err != nil {
			logrus.Errorf("unable to parse request as either JSON or YAML, err: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newGroup := group_v1.Group{}
		if jsonErr := json.Unmarshal(groupData, &newGroup); jsonErr != nil {
			logrus.Errorf("unable to parse object as group, err: %v", jsonErr)
			http.Error(w, jsonErr.Error(), http.StatusBadRequest)
			return
		}
		if _, groupErr := getGroupClient().Create(&newGroup); groupErr != nil {
			logrus.Errorf("unable to create group: %v", groupErr)
			http.Error(w, groupErr.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		eventprocessor.UpdateGroup(&newGroup, getGroupClient())
	}
}

// SetGroupClient sets groupcrd client
func SetGroupClient(conf controller.Config) {
	groupClient = conf.Groupcrdclient
}

func getGroupClient() *v1.GroupClient {
	return groupClient
}
