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
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

// RetrieveAllServicesWithDstPods returns all pods in the dgraph
func RetrieveAllServicesWithDstPods() ([]models.Service, error) {
	const q = `query {
		services(func: has(isService)) {
			xid
			name
			pod {
				xid
				name
				pod {
					xid
					name
				}
			}
		}
	}`

	type root struct {
		Services []models.Service `json:"services"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}

	return newRoot.Services, nil
}

// RetrieveServicesInteractionsForAllLiveServices returns all services in the dgraph
func RetrieveServicesInteractionsForAllLiveServices() ([]models.Service, error) {
	q := `query {
		services(func: has(isService)) @filter((NOT has(endTime))) {
			name
			service {
				name
			}
		}
	}`

	type root struct {
		Services []models.Service `json:"services"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Services, nil
}