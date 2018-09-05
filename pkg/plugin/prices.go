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

package plugin

import (
	"time"
)

func getMonthToDateCostForInstanceType(instanceType string) (float64, float64, float64) {
	hours := totalHoursTillNow()
	basePrice := getPriceForInstanceType(instanceType)
	totalPrice := hours * basePrice
	cpuMemoryRatio := cpuMemoryRatio(instanceType)
	return totalPrice, totalPrice * cpuMemoryRatio, totalPrice * (1 - cpuMemoryRatio)
}

func getMonthToDateCostForStorageClass(storageClass string) float64 {
	percentageOfHoursElapsed := percentageOfHoursElapsedInCurrentMonth()
	basePrice := getPriceForVolumeType(storageClass)
	return basePrice * percentageOfHoursElapsed
}

func percentageOfHoursElapsedInCurrentMonth() float64 {
	now := time.Now()
	hoursTillNow := 24.0*float64(now.Day()-1) + float64(now.Hour())
	totalDays := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).Add(-time.Nanosecond).Day()
	return hoursTillNow / (float64)(totalDays)
}

func totalHoursTillNow() float64 {
	now := time.Now()
	return 24.0*float64(now.Day()-1) + float64(now.Hour())
}

func projectToMonth(val float64) float64 {
	// TODO: enhance this.
	return (val * 31 * 24) / totalHoursTillNow()
}

func cpuMemoryRatio(instanceType string) float64 {
	// TODO: enhance this.
	return 0.3
}
