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
)

// RetrieveReplicasetHierarchy returns hierarchy for a given replicaset
func RetrieveReplicasetHierarchy(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for replicaset, empty name is given")
		return JSONDataWrapper{}
	}
	query := getQueryForHierarchy("isReplicaset", "replicaset", name, "@filter(has(isPod))")
	return getJSONDataFromQuery(query)
}

// RetrieveReplicasetMetrics returns replicaset for a given replicaset
func RetrieveReplicasetMetrics(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for replicaset, empty name is given")
		return JSONDataWrapper{}
	}
	query := getQueryForReplicasetMetrics(name)
	return getJSONDataFromQuery(query)
}
