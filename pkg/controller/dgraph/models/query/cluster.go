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

package query

import (
	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
)

// RetrieveClusterHierarchy returns all namespaces if view is logical and returns all nodes with disks if view is physical
func RetrieveClusterHierarchy(view string) JSONDataWrapper {
	var query string
	if view == Physical {
		query = `query {
			children(func: has(name)) @filter(has(isNode) OR has(isPersistentVolume)) {
				name
				type
			}
		}`
	} else {
		query = `query {
			children(func: has(isNamespace)) {
				name
				type
			}
		}`
	}

	parentRoot := Parent{}
	err := dgraph.ExecuteQuery(query, &parentRoot)
	if err != nil {
		logrus.Errorf("Unable to execute query for retrieving cluster hierarchy: (%v)", err)
		return JSONDataWrapper{}
	}
	root := JSONDataWrapper{
		Data: Parent{
			Name:     "cluster",
			Type:     "cluster",
			Children: parentRoot.Children,
		},
	}
	logrus.Debugf("data: (%v)", root.Data)
	return root
}
