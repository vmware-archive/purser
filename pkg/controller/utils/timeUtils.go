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

package utils

import "time"

// GetCurrentMonthStartTime returns month start time as k8s apimachinery Time object
func GetCurrentMonthStartTime() time.Time {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	return monthStart
}

// ConverTimeToRFC3339 returns query time in RFC3339 format
func ConverTimeToRFC3339(queryTime time.Time) string {
	return queryTime.Format(time.RFC3339)
}

// GetSecondsSince returns number of seconds since query time
func GetSecondsSince(queryTime time.Time) float64 {
	return time.Since(queryTime).Seconds()
}

// GetHoursRemainingInCurrentMonth returns number of hours remaining in the month
func GetHoursRemainingInCurrentMonth() float64 {
	now := time.Now()
	monthEnd := time.Date(now.Year(), now.Month(), 30, 23, 59, 0, 0, time.Local)
	return -time.Since(monthEnd).Hours()
}
