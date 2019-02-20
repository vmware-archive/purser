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

	"github.com/vmware/purser/pkg/controller/dgraph/models"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/utils"
)

// RetrievePVHierarchy returns hierarchy for a given pv
func RetrievePVHierarchy(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for PV, empty name is given")
		return JSONDataWrapper{}
	}

	query := `query {
		parent(func: has(isPersistentVolume)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~pv @filter(has(isPersistentVolumeClaim)) {
				name
				type
			}
        }
    }`
	return getJSONDataFromQuery(query)
}

// RetrievePVMetrics returns metrics for a given pv
func RetrievePVMetrics(name string) JSONDataWrapper {
	if name == All {
		logrus.Errorf("wrong type of query for PV, empty name is given")
		return JSONDataWrapper{}
	}

	secondsSinceMonthStart := fmt.Sprintf("%f", utils.GetSecondsSince(utils.GetCurrentMonthStartTime()))
	query := `query {
		parent(func: has(isPersistentVolume)) @filter(eq(name, "` + name + `")) {
			name
			type
			children: ~pv @filter(has(isPersistentVolumeClaim)) {
				name
				type
				storage: pvcStorage as storageCapacity
				stChild as startTime
				stSecondsChild as math(since(stChild))
				secondsSinceStartChild as math(cond(stSecondsChild > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, stSecondsChild))
				etChild as endTime
				isTerminatedChild as count(endTime)
				secondsSinceEndChild as math(cond(isTerminatedChild == 0, 0.0, since(etChild)))
				durationInHoursChild as math((secondsSinceStartChild - secondsSinceEndChild) / 3600)
				storageCost: math(pvcStorage * durationInHoursChild * ` + models.DefaultStorageCostPerGBPerHour + `)
			}
			storage: storage as storageCapacity
			st as startTime
			stSeconds as math(since(st))
			secondsSinceStart as math(cond(stSeconds > ` + secondsSinceMonthStart + `, ` + secondsSinceMonthStart + `, stSeconds))
			et as endTime
			isTerminated as count(endTime)
			secondsSinceEnd as math(cond(isTerminated == 0, 0.0, since(et)))
			durationInHours as math((secondsSinceStart - secondsSinceEnd) / 3600)
			storageCost: math(storage * durationInHours * ` + models.DefaultStorageCostPerGBPerHour + `)
        }
    }`
	return getJSONDataFromQuery(query)
}
