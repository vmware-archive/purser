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

// RetrieveNamespaceHierarchy returns hierarchy for a given namespace
func RetrieveNamespaceHierarchy(name string) JSONDataWrapper {
	if name == All {
		return RetrieveClusterHierarchy(Logical)
	}

	query := `query {
		parent(func: has(isNamespace)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~namespace @filter(has(isDeployment) OR has(isStatefulset) OR has(isJob) OR has(isDaemonset) OR (has(isReplicaset) AND (NOT has(deployment)))) {
				name
				type
			}
        }
    }`
	return getJSONDataFromQuery(query)
}

// getJSONDataFromQuery executes query and wraps the data in a desired structure(JSONDataWrapper)
func getJSONDataFromQuery(query string) JSONDataWrapper {
	parentRoot := ParentWrapper{}
	err := dgraph.ExecuteQuery(query, &parentRoot)
	if err != nil || len(parentRoot.Parent) == 0 {
		logrus.Errorf("Unable to execute query, err: (%v), length of output: (%d)", err, len(parentRoot.Parent))
		return JSONDataWrapper{}
	}
	root := JSONDataWrapper{
		Data: ParentWrapper{
			Name:     parentRoot.Parent[0].Name,
			Type:     parentRoot.Parent[0].Type,
			Children: parentRoot.Parent[0].Children,
		},
	}
	return root
}
