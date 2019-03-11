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
	"testing"

	"github.com/vmware/purser/test/utils"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestBytesToGB(t *testing.T) {
	act := BytesToGB(124235312345978)
	exp := 115703.15095221438
	utils.Equals(t, exp, act)
}

func TestConvertToFloat64GB(t *testing.T) {
	quantities := getTestQuantities()
	exp := [3]float64{0.011175870895385742, 0.01171875, 0.011175870895385742}
	for index, quantity := range quantities {
		act := ConvertToFloat64GB(&quantity)
		utils.Equals(t, exp[index], act)
	}
}

func TestConvertToFloat64CPU(t *testing.T) {
	quantities := getTestQuantities()
	exp := [3]float64{1.2e+07, 1.2582912e+07, 1.2e+07}
	for index, quantity := range quantities {
		act := ConvertToFloat64CPU(&quantity)
		utils.Equals(t, exp[index], act)
	}
}

func getTestQuantities() [3]resource.Quantity {
	var quantities [3]resource.Quantity
	quantities[0], _ = resource.ParseQuantity("12e6")
	quantities[1], _ = resource.ParseQuantity("12Mi")
	quantities[2], _ = resource.ParseQuantity("12M")

	return quantities
}
