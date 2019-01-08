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

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

const (
	na              = "na"
	aws             = "aws"
	gbMonth         = "GB-Mo"
	deliminator     = "-"
	storageInstance = "Storage"
	computeInstance = "Compute Instance"
)

// GetRateCardForAWS takes region as input and returns RateCard and error if any
func GetRateCardForAWS(region string) (*models.RateCard, error) {
	awsPricing, err := GetAWSPricing(region)
	if err != nil {
		return nil, err
	}
	rateCard := convertAWSPricingToPurserRateCard(region, awsPricing)
	return rateCard, nil
}

func convertAWSPricingToPurserRateCard(region string, awsPricing *Pricing) *models.RateCard {
	nodePrice, storagePrice := getResourcePricesFromAWSPricing(awsPricing)
	return &models.RateCard{
		IsRateCard:    true,
		CloudProvider: aws,
		Region:        region,
		NodePrices:    nodePrice,
		StoragePrices: storagePrice,
	}
}

func getResourcePricesFromAWSPricing(awsPricing *Pricing) ([]*models.NodePrice, []*models.StoragePrice) {
	var nodePrices []*models.NodePrice
	var storagePrices []*models.StoragePrice
	products := awsPricing.Products
	planList := awsPricing.Terms

	duplicateComputeInstanceChecker := make(map[string]bool)
	for _, product := range products {
		if product.ProductFamily == computeInstance {
			key := product.Sku + product.Attributes.InstanceType + product.Attributes.OperatingSystem
			if _, isPresent := duplicateComputeInstanceChecker[key]; !isPresent && product.Attributes.PreInstalledSW == na {
				nodePrice := getComputeInstancePrice(product, planList)
				if nodePrice != nil {
					duplicateComputeInstanceChecker[key] = true
					logrus.Debugf("Node Price: %v", *nodePrice)
					// TODO: store/update nodePrice in dgraph
					nodePrices = append(nodePrices, nodePrice)
				}
			}
		} else if product.ProductFamily == storageInstance {
			storagePrice := getStorageInstancePrice(product, planList)
			if storagePrice != nil {
				logrus.Debugf("Storage Price: %v", *storagePrice)
				// TODO: store/update storagePrice in dgraph
				storagePrices = append(storagePrices, storagePrice)
			}
		}
	}
	return nodePrices, storagePrices
}

func getComputeInstancePrice(product Product, planList PlanList) *models.NodePrice {
	for _, pricingAttributes := range planList.OnDemand[product.Sku] {
		for _, pricingData := range pricingAttributes.PriceDimensions {
			for _, pricePerUnit := range pricingData.PricePerUnit {
				priceInFloat64, err := strconv.ParseFloat(pricePerUnit, 64)
				if err != nil {
					logrus.Errorf("unable to parse string: %s to float. err: %v", pricePerUnit, err)
					return nil
				}
				// Unit of Compute price USD-perHour
				productXID := product.Attributes.InstanceType + deliminator + product.Attributes.InstanceFamily + deliminator + product.Attributes.OperatingSystem
				return &models.NodePrice{
					ID:              dgraph.ID{Xid: productXID},
					IsNodePrice:     true,
					InstanceType:    product.Attributes.InstanceType,
					InstanceFamily:  product.Attributes.InstanceFamily,
					OperatingSystem: product.Attributes.OperatingSystem,
					Price:           priceInFloat64,
				}
			}
		}
	}
	return nil
}

func getStorageInstancePrice(product Product, planList PlanList) *models.StoragePrice {
	for _, pricingAttributes := range planList.OnDemand[product.Sku] {
		for _, pricingData := range pricingAttributes.PriceDimensions {
			for _, pricePerUnit := range pricingData.PricePerUnit {
				priceInFloat64, err := strconv.ParseFloat(pricePerUnit, 64)
				if err != nil {
					logrus.Errorf("unable to parse string: %s to float. err: %v", pricePerUnit, err)
					return nil
				}
				if pricingData.Unit == gbMonth {
					// convert to GBHour
					priceInFloat64 = priceInFloat64 / (30 * 24)
				}
				productXID := product.Attributes.VolumeType + deliminator + product.Attributes.UsageType
				return &models.StoragePrice{
					ID:             dgraph.ID{Xid: productXID},
					IsStoragePrice: true,
					VolumeType:     product.Attributes.VolumeType,
					UsageType:      product.Attributes.UsageType,
					Price:          priceInFloat64,
				}
			}
		}
	}
	return nil
}
