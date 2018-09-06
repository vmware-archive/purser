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
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// nolint
func executeCommand(command string) []byte {
	slice := strings.Fields(command)
	cmd := exec.Command(slice[0], slice[1:]...)
	cmd.Env = os.Environ()
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return out.Bytes()
}

// getCurrentTime returns the current time as k8s apimachinery Time object
func getCurrentTime() metav1.Time {
	return metav1.Now()
}

// getCurrentMonthStartTime returns month start time as k8s apimachinery Time object
func getCurrentMonthStartTime() metav1.Time {
	now := time.Now()
	monthStart := metav1.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	return monthStart
}

/*
currentMonthActiveTimeInHours returns active time (endTime - startTime) in the current month.
1. If startTime is before month start then it is set as month start
2. If endTime is not set(isZero) then it is set as current time
These two conditions ensures that the active time we compute is within the current month.
*/
func currentMonthActiveTimeInHours(startTime, endTime metav1.Time) float64 {
	currentTime := getCurrentTime()
	monthStart := getCurrentMonthStartTime()
	return currentMonthActiveTimeInHoursMulti(startTime, endTime, currentTime, monthStart)
}

/*
currentMonthActiveTimeInHoursMulti is same as currentMonthActiveTimeInHours but it needs extra inputs:
currentTime and monthStart.
Use this method(currentMonthActiveTimeInHoursMulti) if you want to caclculate active time multiple times (ex: inside a loop).
*/
func currentMonthActiveTimeInHoursMulti(startTime, endTime, currentTime, monthStart metav1.Time) float64 {
	if startTime.Time.Before(monthStart.Time) {
		startTime = monthStart
	}

	if endTime.IsZero() {
		endTime = currentTime
	}

	duration := endTime.Time.Sub(startTime.Time)
	durationInHours := duration.Hours()
	return durationInHours
}

// totalHoursTillNow return number of hours from month start to current time.
func totalHoursTillNow() float64 {
	monthStart := getCurrentMonthStartTime()
	currentTime := getCurrentTime()
	return currentMonthActiveTimeInHours(monthStart, currentTime)
}

func projectToMonth(val float64) float64 {
	// TODO: enhance this.
	return (val * 31 * 24) / totalHoursTillNow()
}

func bytesToGB(val int64) float64 {
	return float64(val) / (1024.0 * 1024.0 * 1024.0)
}
