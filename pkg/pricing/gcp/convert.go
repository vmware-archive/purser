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

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
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
		price, _ := getPriceFromSku(&sku)
		switch sku.PricingInfo[0].PricingExpression.UsageUnit {
		case "GiBy.mo":
			price = price / models.HoursInMonth
		case "GiBy.d":
			price = price / 24
		default:
			price = 0 // change?
		}
		productXID := sku.Category.ResourceGroup
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
	var computeByRG map[string][]Skus
	computeByRG = make(map[string][]Skus)
	var cpuPrice, memoryPrice float64
	filterComputePricing(pricing)
	groupNodesByResourceGroupCompute(pricing.Skus, computeByRG)
	for _, skus := range computeByRG {
		// there should ideally be only 1 or 2 skus for a single resource group
		if len(skus) == 1 {
			// only compute
			cpuPrice, _ = getPriceFromSku(&skus[0])
			memoryPrice = 0
		} else if len(skus) == 2 {
			if isCPUSku(&skus[0]) {
				cpuPrice, _ = getPriceFromSku(&skus[0])
				memoryPrice, _ = getPriceFromSku(&skus[1])
			} else {
				cpuPrice, _ = getPriceFromSku(&skus[1])
				memoryPrice, _ = getPriceFromSku(&skus[0])
			}
		} else {
			logrus.Errorf("Unexpected no. of skus in resource group. Skipping...")
			continue
		}
		if cpuPrice != 0 {
			productXID := skus[0].Category.ResourceGroup
			nodePrice := &models.NodePrice{
				ID:              dgraph.ID{Xid: productXID},
				IsNodePrice:     true,
				InstanceType:    skus[0].Category.ResourceGroup,
				InstanceFamily:  "ANY",
				OperatingSystem: "ANY",
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

// assumes all rates are per hour
func isCPUSku(sku *Skus) bool {
	return sku.PricingInfo[0].PricingExpression.UsageUnit == "h"
}

// assumes all rates are per GB per hour
func isMemorySku(sku *Skus) bool {
	return sku.PricingInfo[0].PricingExpression.UsageUnit == "GiBy.h"
}

func filterComputePricing(pricing *Pricing) {
	exceptResourceGroups := []string{"RAM", "GPU", "CPU", "PdSnapshotEgress", "SecurityPolicy"} // put somewhere else
	var newSkus []Skus

	for _, el := range pricing.Skus {
		if el.Category.ResourceFamily == "Compute" &&
			!contains(exceptResourceGroups, el.Category.ResourceGroup) &&
			el.Category.UsageType != "Preemptible" {
			newSkus = append(newSkus, el)
		}
	}
	pricing.Skus = newSkus
}

func filterStoragePricing(pricing *Pricing) {
	allowedResourceGroups := []string{"ColdlineStorage", "RegionalStorage", "DRAStorage", "NearlineStorage"} // put somewhere else
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
		return 0, errors.New("Unexpected PricingInfo")
	}
	for _, tr := range el.PricingInfo[0].PricingExpression.TieredRates {
		units, err := strconv.ParseFloat(tr.UnitPrice.Units, 64)
		if err != nil {
			return 0, errors.New("Invalid units")
		}
		price = units + tr.UnitPrice.Nanos*math.Pow10(-9)
		if price != 0 {
			break
		}
	}
	if price == 0 {
		return 0, errors.New("Unable to get price")
	}
	return price, nil
}
