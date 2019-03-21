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
	"net/http"
	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
	"github.com/Sirupsen/logrus"
)

// GetGroupsData listens on /groups endpoint
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