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

package aws

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/utils"
)

const (
	httpTimeout = 100 * time.Second
)

// Pricing structure
type Pricing struct {
	Products map[string]Product
	Terms    PlanList
}

// PlanList structure
type PlanList struct {
	OnDemand map[string]map[string]TermAttributes
}

// TermAttributes structure
type TermAttributes struct {
	PriceDimensions map[string]PricingData
}

// PricingData structure
type PricingData struct {
	Unit         string
	PricePerUnit map[string]string
}

// Product structure
type Product struct {
	Sku           string
	ProductFamily string
	Attributes    ProductAttributes
}

// ProductAttributes structure
type ProductAttributes struct {
	InstanceType    string
	InstanceFamily  string
	OperatingSystem string
	PreInstalledSW  string
	VolumeType      string
	UsageType       string
}

// GetAWSPricing function details
// input: region
// retrieves data from http get to the corresponding url for that region
func GetAWSPricing(region string) (*Pricing, error) {
	var myClient = &http.Client{Timeout: httpTimeout}
	rateCard := Pricing{}
	err := utils.GetJSONResponse(myClient, getURLForRegion(region), &rateCard)
	if err != nil {
		logrus.Errorf("Unable to get aws pricing. Reason: %v", err)
		return nil, err
	}
	return &rateCard, nil
}

func getURLForRegion(region string) string {
	return "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/" + region + "/index.json"
}
