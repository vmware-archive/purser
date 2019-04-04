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

package dgraph

import (
	"github.com/vmware/purser/pkg/controller/utils"

	log "github.com/Sirupsen/logrus"
	"time"
)

type resource struct {
	ID
}

// RemoveResourcesInactive deletes all resources which have their deletion time stamp before
// the start of current month.
func RemoveResourcesInactive() {
	err := removeOldDeletedResources()
	if err != nil {
		log.Println(err)
	}

	err = removeOldDeletedPods()
	if err != nil {
		log.Error(err)
	}
}

func removeOldDeletedResources() error {
	uids, err := retrieveResourcesWithEndTimeBeforeCurrentMonthStart()
	if err != nil {
		return err
	}
	if len(uids) == 0 {
		log.Println("No old deleted resources are present in dgraph")
		return nil
	}

	_, err = MutateNode(uids, DELETE)
	return err
}

func removeOldDeletedPods() error {
	uids, err := retrievePodsWithEndTimeBeforeThreeMonths()
	if err != nil {
		return err
	}
	if len(uids) == 0 {
		log.Println("No old deleted pods are present in dgraph")
		return nil
	}

	_, err = MutateNode(uids, DELETE)
	return err
}

func retrieveResourcesWithEndTimeBeforeCurrentMonthStart() ([]resource, error) {
	q := `query {
		resources(func: le(endTime, "` + utils.ConverTimeToRFC3339(utils.GetCurrentMonthStartTime()) + `")) @filter(NOT(has(isPod))) {
			uid
		}
	}`

	type root struct {
		Resources []resource `json:"resources"`
	}
	newRoot := root{}
	err := ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Resources, nil
}

func retrievePodsWithEndTimeBeforeThreeMonths() ([]resource, error) {
	q := `query {
		resources(func: le(endTime, "` + utils.ConverTimeToRFC3339(utils.GetCurrentMonthStartTime().Add(-time.Hour*24*30*2)) + `")) @filter(has(isPod)) {
			uid
		}
	}`

	type root struct {
		Resources []resource `json:"resources"`
	}
	newRoot := root{}
	err := ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}
	return newRoot.Resources, nil
}
