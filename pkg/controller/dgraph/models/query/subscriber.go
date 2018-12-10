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

// RetrieveSubscribers gets all live subscribers
func RetrieveSubscribers() ([]models.SubscriberCRD, error) {
	q := `query {
		subscribers(func: has(isSubscriber)) @filter(NOT(has(endTime))) {
			name
			Spec {
				headers
				url
			}
		}
	}`
	type root struct {
		Subscribers []models.SubscriberCRD `json:"subscribers"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Subscribers, nil
}