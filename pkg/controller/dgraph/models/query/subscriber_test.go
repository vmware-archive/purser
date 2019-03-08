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
	"fmt"
	"testing"

	"github.com/vmware/purser/pkg/controller/dgraph/models"

	"github.com/stretchr/testify/assert"
)

func mockDgraphForSubscriberQueries(queryType string) {
	executeQuery = func(query string, root interface{}) error {
		dummySubscriberList, ok := root.(*subscriberRoot)
		if !ok {
			return fmt.Errorf("wrong root received")
		}

		if queryType == testRetrieveSubscribers {
			dummySubscriber := models.SubscriberCRD{
				Name: "subscriber-purser",
				Spec: models.SubscriberSpec{
					URL: "http://purer.com",
				},
			}
			dummySubscriberList.Subscribers = []models.SubscriberCRD{dummySubscriber}
			return nil
		}

		return fmt.Errorf("no data found")
	}
}

// TestRetrieveSubscribersWithDgraphError ...
func TestRetrieveSubscribersWithDgraphError(t *testing.T) {
	mockDgraphForSubscriberQueries(testWrongQuery)
	_, err := RetrieveSubscribers()
	assert.Error(t, err)
}

// TestRetrieveSubscribers ...
func TestRetrieveSubscribers(t *testing.T) {
	mockDgraphForSubscriberQueries(testRetrieveSubscribers)
	got, err := RetrieveSubscribers()
	expected := []models.SubscriberCRD{{
		Name: "subscriber-purser",
		Spec: models.SubscriberSpec{
			URL: "http://purer.com",
		},
	}}
	assert.Equal(t, expected, got)
	assert.NoError(t, err)
}
