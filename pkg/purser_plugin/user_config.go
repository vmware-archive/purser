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

package purser_plugin

import (
	"fmt"
	"strconv"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace                      = "default"
	userCostsConfigMap             = "purser-user-costs"
	defaultCPUCostPerCPUPerHour    = float64(0.024)
	defaultMemCostPerGBPerHour     = float64(0.01)
	defaultStorageCostPerGBPerHour = float64(0.0033)
)

/*
SaveUserCosts stores the cpu, memory and storage cost per unit per hour in the cluster as config maps.
*/
func SaveUserCosts(cpuCostPerCPUPerHour, memCostPerGBPerHour, storageCostPerGBPerHour string) bool {
	cm, err := ClientSetInstance.CoreV1().ConfigMaps(namespace).Get(userCostsConfigMap, metav1.GetOptions{})
	if err != nil {
		// no configmap so create new one
		mp := map[string]string{}
		mp["cpuCostPerCPUPerHour"] = cpuCostPerCPUPerHour
		mp["memCostPerGBPerHour"] = memCostPerGBPerHour
		mp["storageCostPerGBPerHour"] = storageCostPerGBPerHour
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: userCostsConfigMap,
			},
			Data: mp,
		}
		_, err2 := ClientSetInstance.CoreV1().ConfigMaps(namespace).Create(cm)
		if err2 != nil {
			fmt.Printf("Error in createing config map: %s", err2)
			return false
		}
	} else {
		// update configmap
		cm.Data["cpuCostPerCPUPerHour"] = cpuCostPerCPUPerHour
		cm.Data["memCostPerGBPerHour"] = memCostPerGBPerHour
		cm.Data["storageCostPerGBPerHour"] = storageCostPerGBPerHour
		_, err2 := ClientSetInstance.CoreV1().ConfigMaps(namespace).Update(cm)
		if err2 != nil {
			fmt.Printf("Error in updating config map: %s", err2)
			return false
		}
		fmt.Printf("Updated config map\n")
	}

	return true
}

/*
GetUserCosts gives the cpu, memory and storage cost per unit per hour which are stored in the cluster as config maps.
*/
func GetUserCosts() (float64, float64, float64) {
	var cpuCostPerCPUPerHour, memCostPerGBPerHour, storageCostPerGBPerHour float64
	cm, err := ClientSetInstance.CoreV1().ConfigMaps(namespace).Get(userCostsConfigMap, metav1.GetOptions{})
	if err != nil {
		// no user configed costs. so return default values
		cpuCostPerCPUPerHour = defaultCPUCostPerCPUPerHour
		memCostPerGBPerHour = defaultMemCostPerGBPerHour
		storageCostPerGBPerHour = defaultStorageCostPerGBPerHour
	} else {
		cpuCostPerCPUPerHour = strToFloat64(cm.Data["cpuCostPerCPUPerHour"])
		if cpuCostPerCPUPerHour < 0.0 {
			cpuCostPerCPUPerHour = defaultCPUCostPerCPUPerHour
		}

		memCostPerGBPerHour = strToFloat64(cm.Data["memCostPerGBPerHour"])
		if memCostPerGBPerHour < 0.0 {
			memCostPerGBPerHour = defaultMemCostPerGBPerHour
		}

		storageCostPerGBPerHour = strToFloat64(cm.Data["storageCostPerGBPerHour"])
		if storageCostPerGBPerHour < 0.0 {
			storageCostPerGBPerHour = defaultStorageCostPerGBPerHour
		}
	}

	return cpuCostPerCPUPerHour, memCostPerGBPerHour, storageCostPerGBPerHour
}

func strToFloat64(strNum string) float64 {
	n, err := strconv.ParseFloat(strNum, 64)
	if err != nil {
		fmt.Printf("Error converting string %s to float\n", strNum)
		return -1.0
	}
	return n
}
