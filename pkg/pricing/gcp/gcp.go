/*
 * Copyright (c) 2019 VMware Inc. All Rights Reserved.
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

package gcp

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/utils"
)

const (
	httpTimeout = 100 * time.Second
	apiKey      = "your-api-key-here"
	computeURL  = "https://cloudbilling.googleapis.com/v1/services/6F81-5844-456A/skus"
	storageURL  = "https://cloudbilling.googleapis.com/v1/services/95FF-2EF5-5EA1/skus"
)

/*
 * The following structs are modelled according to the GCP Pricing API response
 * https://cloud.google.com/billing/v1/how-tos/catalog-api
 */

// Pricing structure
type Pricing struct {
	Skus          []Skus `json:"skus"`
	NextPageToken string `json:"nextPageToken"`
}

// Category structure
type Category struct {
	ServiceDisplayName string `json:"serviceDisplayName"`
	ResourceFamily     string `json:"resourceFamily"`
	ResourceGroup      string `json:"resourceGroup"`
	UsageType          string `json:"usageType"`
}

// UnitPrice structure
type UnitPrice struct {
	CurrencyCode string  `json:"currencyCode"`
	Units        string  `json:"units"`
	Nanos        float64 `json:"nanos"`
}

// TieredRates structure
type TieredRates struct {
	StartUsageAmount int       `json:"startUsageAmount"`
	UnitPrice        UnitPrice `json:"unitPrice"`
}

// PricingExpression structure
type PricingExpression struct {
	UsageUnit            string        `json:"usageUnit"`
	UsageUnitDescription string        `json:"usageUnitDescription"`
	DisplayQuantity      int           `json:"displayQuantity"`
	TieredRates          []TieredRates `json:"tieredRates"`
}

// AggregationInfo structure
type AggregationInfo struct {
	AggregationLevel    string `json:"aggregationLevel"`
	AggregationInterval string `json:"aggregationInterval"`
	AggregationCount    int    `json:"aggregationCount"`
}

// PricingInfo structure
type PricingInfo struct {
	EffectiveTime          string            `json:"effectiveTime"`
	Summary                string            `json:"summary"`
	PricingExpression      PricingExpression `json:"pricingExpression"`
	AggregationInfo        AggregationInfo   `json:"aggregationInfo"`
	CurrencyConversionRate int               `json:"currencyConversionRate"`
}

// Skus structure
type Skus struct {
	Name                string        `json:"name"`
	SkuID               string        `json:"skuId"`
	Description         string        `json:"description"`
	Category            Category      `json:"category"`
	ServiceRegions      []string      `json:"serviceRegions"`
	PricingInfo         []PricingInfo `json:"pricingInfo"`
	ServiceProviderName string        `json:"serviceProviderName"`
}

func getGCPPRicingHelper(url string, region string) (*Pricing, error) {
	var myClient = &http.Client{Timeout: httpTimeout}
	pricing := Pricing{}
	for {
		temp := Pricing{}
		err := utils.GetJSONResponse(myClient, url, &temp)
		if err != nil {
			logrus.Errorf("Unable to get gcp pricing. Reason: %v", err)
			return nil, err
		}

		for _, sku := range temp.Skus {
			if contains(sku.ServiceRegions, region) {
				// handle "global" region?
				pricing.Skus = append(pricing.Skus, sku)
			}
		}

		if temp.NextPageToken != "" {
			url = fmt.Sprintf("%s&pageToken=%s", url, temp.NextPageToken)
		} else {
			break
		}
	}

	if len(pricing.Skus) == 0 {
		return nil, errors.New("Unable to fetch data from api")
	}

	return &pricing, nil
}

// GetGCPPricingCompute calls the gcp pricing api and filters on region
func GetGCPPricingCompute(region string) (*Pricing, error) {
	pricing, err := getGCPPRicingHelper(fmt.Sprintf("%s?key=%s", computeURL, apiKey), region)
	return pricing, err
}

// GetGCPPricingStorage calls the gcp pricing api and filters on region
func GetGCPPricingStorage(region string) (*Pricing, error) {
	pricing, err := getGCPPRicingHelper(fmt.Sprintf("%s?key=%s", storageURL, apiKey), region)
	return pricing, err
}
