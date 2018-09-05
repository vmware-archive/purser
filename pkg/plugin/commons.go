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
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
func getCurrentTime() meta_v1.Time {
	return meta_v1.Now()
}

// getCurrentMonthStartTime returns month start time as k8s apimachinery Time object
func getCurrentMonthStartTime() meta_v1.Time {
	now := time.Now()
	monthStart := meta_v1.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	return monthStart
}

/*
currentMonthActiveTimeInHours returns active time (endTime - startTime) in the current month.
1. If startTime is before month start then it is set as month start
2. If endTime is not set(isZero) then it is set as current time
These two conditions ensures that the active time we compute is within the current month.
*/
func currentMonthActiveTimeInHours(startTime, endTime, currentTime, monthStart meta_v1.Time) float64 {
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

func resourceQuantityToFloat64(quantity *resource.Quantity) float64 {
	val, isSuccess := quantity.AsInt64()
	if !isSuccess {
		fmt.Println("Unable to convert resource quantity into int64")
	}
	return float64(val)
}
