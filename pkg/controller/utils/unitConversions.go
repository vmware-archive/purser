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

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/resource"
)

// BytesToGB converts from bytes(int64) to GB(float64)
func BytesToGB(val int64) float64 {
	return float64BytesToFloat64GB(float64(val))
}

// ConvertToFloat64GB quantity to float64 GB
func ConvertToFloat64GB(quantity *resource.Quantity) float64 {
	return float64BytesToFloat64GB(resourceToFloat64(quantity))
}

// ConvertToFloat64CPU quantity to float64 vCPU
func ConvertToFloat64CPU(quantity *resource.Quantity) float64 {
	return resourceToFloat64(quantity)
}

// AddResourceAToResourceB ...
func AddResourceAToResourceB(resA, resB *resource.Quantity) {
	if resA != nil {
		resB.Add(*resA)
	}
}

// float64BytesToFloat64GB from bytes (float64) to GB(float64)
func float64BytesToFloat64GB(val float64) float64 {
	return val / (1024.0 * 1024.0 * 1024.0)
}

// resourceToFloat64 ...
func resourceToFloat64(quantity *resource.Quantity) float64 {
	decVal := quantity.AsDec()
	decValueFloat, err := strconv.ParseFloat(decVal.String(), 64)
	if err != nil {
		log.Errorf("error while converting into string: (%s) to float\n", decVal.String())
	}
	return decValueFloat // 0 if not isSuccess
}
