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
	"fmt"
	"log"
	"time"

	"github.com/vmware/purser/pkg/controller/utils"
)

type resource struct {
	ID
}

// RemoveResourcesInactiveInCurrentMonth deletes all resources which have their deletion time stamp before
// the start of current month.
func RemoveResourcesInactiveInCurrentMonth() {
	err := removeOldDeletedResources()
	if err != nil {
		log.Println(err)
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

	toDelete := utils.JSONMarshal(uids)
	if toDelete == nil {
		return fmt.Errorf("Unable to marshal toDelete uids")
	}
	return DeleteNodes(toDelete)
}

func retrieveResourcesWithEndTimeBeforeCurrentMonthStart() ([]resource, error) {
	q := `query {
		resources(func: le(endTime, "` + getCurrentMonthStartTime() + `")) {
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

// getCurrentMonthStartTime returns month start time as k8s apimachinery Time object
func getCurrentMonthStartTime() string {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	return monthStart.Format(time.RFC3339)
}
