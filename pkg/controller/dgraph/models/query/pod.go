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

// RetrievePodsInteractions returns inbound and outbound interactions of a pod
func RetrievePodsInteractions(name string, isOrphan bool) []byte {
	var query string
	if name == All {
		if isOrphan {
			query = `query {
				pods(func: has(isPod)) {
					name
					outbound: pod {
						name
					}
					inbound: ~pod @filter(has(isPod)) {
						name
					}
				}
			}`
		} else {
			query = `query {
				pods(func: has(isPod)) @filter(has(pod)) {
					name
					outbound: pod {
						name
					}
					inbound: ~pod @filter(has(isPod)) {
						name
					}
				}
			}`
		}
	} else {
		query = `query {
			pods(func: has(isPod)) @filter(eq(name, "` + name + `")) {
				name
				outbound: pod {
					name
				}
				inbound: ~pod @filter(has(isPod)) {
					name
				}
			}
		}`
	}

	result, err := dgraph.ExecuteQueryRaw(query)
	if err != nil {
		logrus.Errorf("Error while retrieving query for pods interactions. Name: (%v), isOrphan: (%v), error: (%v)", name, isOrphan, err)
		return nil
	}
	return result
}
