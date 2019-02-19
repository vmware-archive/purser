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
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

const (
	na              = "NA"
	aws             = "aws"
	gbMonth         = "GB-Mo"
	deliminator     = "-"
	storageInstance = "Storage"
	computeInstance = "Compute Instance"
	priceError      = -1.0
	hoursInMonth    = 720

	// TODO: Determine priceSplitRatio according to instance type i.e, compute optimized or memory optimized etc
	priceSplitRatio             = 0.5
	defaultCPUCostPerCPUPerHour = "0.024"
	defaultMemCostPerGBPerHour  = "0.01"
	defaultStorageCostInFloat64 = 0.00013888888
)

// GetRateCardForAWS takes region as input and returns RateCard and error if any
func GetRateCardForAWS(region string) *models.RateCard {
	awsPricing, err := GetAWSPricing(region)
	if err == nil {
		return convertAWSPricingToPurserRateCard(region, awsPricing)
	}
	return nil
}

func convertAWSPricingToPurserRateCard(region string, awsPricing *Pricing) *models.RateCard {
	nodePrices, storagePrices := getResourcePricesFromAWSPricing(awsPricing)
	return &models.RateCard{
		ID:            dgraph.ID{Xid: models.RateCardXID},
		IsRateCard:    true,
		CloudProvider: aws,
		Region:        region,
		NodePrices:    nodePrices,
		StoragePrices: storagePrices,
	}
}

func getResourcePricesFromAWSPricing(awsPricing *Pricing) ([]*models.NodePrice, []*models.StoragePrice) {
	var nodePrices []*models.NodePrice
	var storagePrices []*models.StoragePrice
	products := awsPricing.Products
	planList := awsPricing.Terms

	duplicateComputeInstanceChecker := make(map[string]bool)
	for _, product := range products {
		priceInFloat64, unit := getResourcePrice(product, planList)
		switch product.ProductFamily {
		case computeInstance:
			pricePerCPU, pricePerGB := getPriceForUnitResource(product, priceInFloat64)
			nodePrices = updateComputeInstancePrices(
				product, priceInFloat64, pricePerCPU, pricePerGB, duplicateComputeInstanceChecker, nodePrices)
		case storageInstance:
			storagePrices = updateStorageInstancePrices(product, priceInFloat64, unit, storagePrices)
		}

	}
	return nodePrices, storagePrices
}

func getResourcePrice(product Product, planList PlanList) (float64, string) {
	for _, pricingAttributes := range planList.OnDemand[product.Sku] {
		for _, pricingData := range pricingAttributes.PriceDimensions {
			for _, pricePerUnit := range pricingData.PricePerUnit {
				priceInFloat64, err := strconv.ParseFloat(pricePerUnit, 64)
				if err != nil {
					logrus.Errorf("unable to parse string: %s to float. err: %v", pricePerUnit, err)
					return priceError, "" // negative price means error
				}
				return priceInFloat64, pricingData.Unit
			}
		}
	}
	return priceError, ""
}

func getPriceForUnitResource(product Product, priceInFloat64 float64) (string, string) {
	pricePerCPU := defaultCPUCostPerCPUPerHour
	pricePerGB := defaultMemCostPerGBPerHour
	if priceInFloat64 != priceError {
		cpu, err := strconv.ParseFloat(product.Attributes.Vcpu, 64)
		if err == nil {
			pricePerCPU = strconv.FormatFloat(priceSplitRatio*priceInFloat64/cpu, 'f', 11, 64)
		}

		memWithUnits := product.Attributes.Memory
		mem, err := strconv.ParseFloat(strings.Join(strings.Split(strings.Split(memWithUnits, " GiB")[0], ","), ""), 64)
		if err == nil {
			pricePerCPU = strconv.FormatFloat(priceSplitRatio*priceInFloat64/mem, 'f', 11, 64)
		}
	}
	return pricePerCPU, pricePerGB
}

func updateComputeInstancePrices(product Product, priceInFloat64 float64, pricePerCPU, pricePerGB string, duplicateComputeInstanceChecker map[string]bool, nodePrices []*models.NodePrice) []*models.NodePrice {
	key := product.Sku + product.Attributes.InstanceType + product.Attributes.OperatingSystem
	if _, isPresent := duplicateComputeInstanceChecker[key]; !isPresent && product.Attributes.PreInstalledSW == na {
		// Unit of Compute price USD-perHour
		productXID := product.Attributes.InstanceType + deliminator + product.Attributes.InstanceFamily + deliminator + product.Attributes.OperatingSystem
		nodePrice := &models.NodePrice{
			ID:              dgraph.ID{Xid: productXID},
			IsNodePrice:     true,
			InstanceType:    product.Attributes.InstanceType,
			InstanceFamily:  product.Attributes.InstanceFamily,
			OperatingSystem: product.Attributes.OperatingSystem,
			Price:           priceInFloat64,
			PricePerCPU:     pricePerCPU,
			PricePerMemory:  pricePerGB,
		}
		duplicateComputeInstanceChecker[key] = true
		uid := models.StoreNodePrice(nodePrice, productXID)
		if uid != "" {
			nodePrice.ID = dgraph.ID{UID: uid, Xid: productXID}
			nodePrices = append(nodePrices, nodePrice)
		}
	}
	return nodePrices
}

func updateStorageInstancePrices(product Product, priceInFloat64 float64, unit string, storagePrices []*models.StoragePrice) []*models.StoragePrice {
	if priceInFloat64 == priceError {
		priceInFloat64 = defaultStorageCostInFloat64
	} else if unit == gbMonth {
		// convert to GBHour
		priceInFloat64 = priceInFloat64 / hoursInMonth
	}

	productXID := product.Attributes.VolumeType + deliminator + product.Attributes.UsageType
	storagePrice := &models.StoragePrice{
		ID:             dgraph.ID{Xid: productXID},
		IsStoragePrice: true,
		VolumeType:     product.Attributes.VolumeType,
		UsageType:      product.Attributes.UsageType,
		Price:          priceInFloat64,
		PricePerGB:     strconv.FormatFloat(priceInFloat64, 'f', 11, 64),
	}
	uid := models.StoreStoragePrice(storagePrice, productXID)
	if uid != "" {
		storagePrice.ID = dgraph.ID{UID: uid, Xid: productXID}
		storagePrices = append(storagePrices, storagePrice)
	}
	return storagePrices
}
