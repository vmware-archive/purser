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
	"math"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

const (
	defaultOS = "linux"
)

// put somewhere else
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// GetRateCardForGCP gets the rate card for a given region
func GetRateCardForGCP(region string) *models.RateCard {
	computePricing, err1 := GetGCPPricingCompute(region)
	storagePricing, err2 := GetGCPPricingStorage(region)
	if err1 == nil && err2 == nil {
		return getPurserRateCard(region, computePricing, storagePricing)
	}
	return nil
}

func getPurserRateCard(region string, computePricing *Pricing, storagePricing *Pricing) *models.RateCard {
	nodePrices := getNodePricesFromGCPPricing(computePricing)
	storagePrices := getStoragePricesFromGCPPricing(storagePricing)
	return &models.RateCard{
		ID:            dgraph.ID{Xid: models.RateCardXID},
		IsRateCard:    true,
		CloudProvider: models.GCP,
		Region:        region,
		NodePrices:    nodePrices,
		StoragePrices: storagePrices,
	}
}

func getStoragePricesFromGCPPricing(pricing *Pricing) []*models.StoragePrice {
	var storagePrices []*models.StoragePrice
	filterStoragePricing(pricing)
	for _, sku := range pricing.Skus {
		price, err := getPriceFromSku(&sku)
		if err != nil {
			logrus.Printf("Unable to get price for sku %s", sku.SkuID)
			continue
		}
		switch sku.PricingInfo[0].PricingExpression.UsageUnit {
		case "GiBy.mo":
			price = price / models.HoursInMonth
		case "GiBy.d":
			price = price / 24
		default:
			logrus.Printf("Unexpected storage price unit for sku %s", sku.SkuID)
			continue
		}
		// see getPricePerUnitResourceFromNodePrice
		productXID := sku.Category.ResourceGroup + "-" + defaultOS
		storagePrice := &models.StoragePrice{
			ID:             dgraph.ID{Xid: productXID},
			IsStoragePrice: true,
			VolumeType:     "ANY",
			UsageType:      sku.Category.ResourceGroup,
			Price:          price,
		}
		uid := models.StoreStoragePrice(storagePrice, productXID)
		if uid != "" {
			storagePrice.ID = dgraph.ID{UID: uid, Xid: productXID}
			storagePrices = append(storagePrices, storagePrice)
		}
	}
	return storagePrices
}

func getNodePricesFromGCPPricing(pricing *Pricing) []*models.NodePrice {
	var nodePrices []*models.NodePrice
	computeByRG := make(map[string][]Skus)
	var cpuPrice, memoryPrice float64
	var err1, err2 error
	filterComputePricing(pricing)
	groupNodesByResourceGroupCompute(pricing.Skus, computeByRG)
	for rg, skus := range computeByRG {
		// there should ideally be only 1 or 2 skus for a single resource group
		if len(skus) == 1 {
			// only compute
			cpuPrice, err1 = getPriceFromSku(&skus[0])
			memoryPrice = 0
		} else if len(skus) == 2 {
			if isCPUSku(&skus[0]) {
				cpuPrice, err1 = getPriceFromSku(&skus[0])
				memoryPrice, err2 = getPriceFromSku(&skus[1])
			} else {
				cpuPrice, err1 = getPriceFromSku(&skus[1])
				memoryPrice, err2 = getPriceFromSku(&skus[0])
			}
		} else {
			logrus.Errorf("Unexpected no. of skus in resource group. Skipping...")
			continue
		}
		if err1 != nil || err2 != nil {
			logrus.Printf("Unable to get cpu/memory price for resource group %s", rg)
			continue
		}
		if cpuPrice != 0 {
			// see getPricePerUnitResourceFromNodePrice
			productXID := resourceGroupToInstanceType(skus[0].Category.ResourceGroup) + "-" + defaultOS
			nodePrice := &models.NodePrice{
				ID:              dgraph.ID{Xid: productXID},
				IsNodePrice:     true,
				InstanceType:    resourceGroupToInstanceType(skus[0].Category.ResourceGroup),
				InstanceFamily:  "ANY",
				OperatingSystem: defaultOS,
				Price:           0,
				PricePerCPU:     cpuPrice,
				PricePerMemory:  memoryPrice,
			}
			uid := models.StoreNodePrice(nodePrice, productXID)
			if uid != "" {
				nodePrice.ID = dgraph.ID{UID: uid, Xid: productXID}
				nodePrices = append(nodePrices, nodePrice)
			}
		}
	}
	return nodePrices
}

func groupNodesByResourceGroupCompute(skus []Skus, group map[string][]Skus) {
	for _, sku := range skus {
		group[sku.Category.ResourceGroup] = append(group[sku.Category.ResourceGroup], sku)
	}
}

// used to differentiate between a cpu sku and its corresponding memory sku
// assumes all cpu rates are per hour
func isCPUSku(sku *Skus) bool {
	return sku.PricingInfo[0].PricingExpression.UsageUnit == "h"
}

func filterRegion(pricing *Pricing, region string) {
	var newSkus []Skus

	for _, el := range pricing.Skus {
		if contains(el.ServiceRegions, region) {
			newSkus = append(newSkus, el)
		}
	}
	pricing.Skus = newSkus
}

func filterComputePricing(pricing *Pricing) {
	allowedResourceGroups := []string{"N1Standard"} // put somewhere else
	var newSkus []Skus

	for _, el := range pricing.Skus {
		if el.Category.ResourceFamily == "Compute" &&
			contains(allowedResourceGroups, el.Category.ResourceGroup) &&
			el.Category.UsageType != "Preemptible" {
			newSkus = append(newSkus, el)
		}
	}
	pricing.Skus = newSkus
}

func filterStoragePricing(pricing *Pricing) {
	allowedResourceGroups := []string{"ColdlineStorage", "RegionalStorage", "NearlineStorage"} // put somewhere else
	var newSkus []Skus

	for _, el := range pricing.Skus {
		if el.Category.ResourceFamily == "Storage" &&
			contains(allowedResourceGroups, el.Category.ResourceGroup) &&
			el.Category.UsageType != "Preemptible" {
			newSkus = append(newSkus, el)
		}
	}
	pricing.Skus = newSkus
}

// PricingInfo and TieredRates are usually length 1
func getPriceFromSku(sku *Skus) (float64, error) {
	var price float64
	el := *sku
	if len(el.PricingInfo) != 1 {
		return 0, errors.New("unexpected PricingInfo")
	}
	for _, tr := range el.PricingInfo[0].PricingExpression.TieredRates {
		units, err := strconv.ParseFloat(tr.UnitPrice.Units, 64)
		if err != nil {
			return 0, errors.New("invalid units")
		}
		price = units + tr.UnitPrice.Nanos*math.Pow10(-9)
		if price != 0 {
			break
		}
	}
	if price == 0 {
		return 0, errors.New("unable to get price")
	}
	return price, nil
}

// according to how GCP sets beta.kubernetes.io/instance-type
func resourceGroupToInstanceType(rg string) string {
	return strings.ToLower(rg[:2]) + "-" + strings.ToLower(rg[2:])
}
