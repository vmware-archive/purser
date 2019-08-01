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
	"fmt"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

// AWS specific constants
const (
	na              = "NA"
	gbMonth         = "GB-Mo"
	deliminator     = "-"
	storageInstance = "Storage"
	computeInstance = "Compute Instance"

	// TODO: Determine priceSplitRatio according to instance type i.e, compute optimized or memory optimized etc
	priceSplitRatio = 0.5
)

//
type bestNodePrice struct {
	CPU         float64
	Memory      float64
	CPUPrice    float64
	MemoryPrice float64
	Total       float64
	NodePrice   *models.NodePrice
}

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
		ID:            dgraph.ID{Xid: models.RateCardXID}, //look into dgraph
		IsRateCard:    true,
		CloudProvider: models.AWS,
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
			nodePrices = updateComputeInstancePrices(product, priceInFloat64, duplicateComputeInstanceChecker, nodePrices)
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
					return models.PriceError, "" // negative price means error
				}
				return priceInFloat64, pricingData.Unit
			}
		}
	}
	return models.PriceError, ""
}

func updateComputeInstancePrices(product Product, priceInFloat64 float64, duplicateComputeInstanceChecker map[string]bool, nodePrices []*models.NodePrice) []*models.NodePrice {
	key := product.Sku + product.Attributes.InstanceType + product.Attributes.OperatingSystem
	if _, isPresent := duplicateComputeInstanceChecker[key]; !isPresent && product.Attributes.PreInstalledSW == na {
		// Unit of Compute price USD-perHour
		productXID := product.Attributes.InstanceType + deliminator + product.Attributes.OperatingSystem
		pricePerCPU, pricePerGB := getPriceForUnitResource(product, priceInFloat64)
		nCPU, _ := strconv.ParseFloat(product.Attributes.Vcpu, 64)
		nMemory, _ := strconv.ParseFloat(strings.Split(product.Attributes.Memory, " ")[0], 64)
		nodePrice := &models.NodePrice{
			ID:              dgraph.ID{Xid: productXID},
			IsNodePrice:     true,
			InstanceType:    product.Attributes.InstanceType,
			InstanceFamily:  product.Attributes.InstanceFamily,
			OperatingSystem: product.Attributes.OperatingSystem,
			Price:           priceInFloat64,
			PricePerCPU:     pricePerCPU,
			PricePerMemory:  pricePerGB,
			CPU:             nCPU,
			Memory:          nMemory,
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
	if priceInFloat64 == models.PriceError {
		priceInFloat64 = models.DefaultStorageCostInFloat64
	} else if unit == gbMonth {
		// convert to GBHour
		priceInFloat64 = priceInFloat64 / models.HoursInMonth
	}

	productXID := product.Attributes.VolumeType + deliminator + product.Attributes.UsageType
	storagePrice := &models.StoragePrice{
		ID:             dgraph.ID{Xid: productXID},
		IsStoragePrice: true,
		VolumeType:     product.Attributes.VolumeType,
		UsageType:      product.Attributes.UsageType,
		Price:          priceInFloat64,
	}
	uid := models.StoreStoragePrice(storagePrice, productXID)
	if uid != "" {
		storagePrice.ID = dgraph.ID{UID: uid, Xid: productXID}
		storagePrices = append(storagePrices, storagePrice)
	}
	return storagePrices
}

func getPriceForUnitResource(product Product, priceInFloat64 float64) (float64, float64) {
	pricePerCPU := models.DefaultCPUCostInFloat64
	pricePerGB := models.DefaultMemCostInFloat64

	// priceInFloat64 should be greater than 0 otherwise this function returns default pricing
	if priceInFloat64 != models.PriceError && priceInFloat64 != 0 {
		cpu, err := strconv.ParseFloat(product.Attributes.Vcpu, 64)
		if err == nil {
			pricePerCPU = priceSplitRatio * priceInFloat64 / cpu
		}

		memWithUnits := product.Attributes.Memory
		// memWithUnits format: "3,126 GiB"
		mem, err := strconv.ParseFloat(strings.Join(strings.Split(strings.Split(memWithUnits, " GiB")[0], ","), ""), 64)
		if err == nil {
			pricePerGB = (1 - priceSplitRatio) * priceInFloat64 / mem
		}
	}
	return pricePerCPU, pricePerGB
}

// GetAwsNodesCost ..
func GetAwsNodesCost(nodes []models.Node, region string) []models.ClusterNodePrice {
	nodePrices, _ := models.GetRateCardForRegion(models.AWS, region)
	var clusterNodePrices []models.ClusterNodePrice
	for _, node := range nodes {
		nodePrice, _ := getBestNodePriceForNode(node, nodePrices)
		clusterNodePrices = append(clusterNodePrices, models.ClusterNodePrice{
			InstanceType:    nodePrice.InstanceType,
			OperatingSystem: nodePrice.OperatingSystem,
			Price:           nodePrice.Price,
			CPUCost:         nodePrice.CPU * nodePrice.PricePerCPU,
			MemoryCost:      nodePrice.Memory * nodePrice.PricePerMemory,
			CPU:             nodePrice.CPU,
			Memory:          nodePrice.Memory,
		})
	}
	logrus.Printf("%#v", clusterNodePrices)
	return clusterNodePrices
}

//getBestNodePriceForNode ..
func getBestNodePriceForNode(node models.Node, nodePrices []*models.NodePrice) (models.NodePrice, error) {
	var bestNP bestNodePrice
	logrus.Printf("%#v", node)

	for _, nodePrice := range nodePrices {
		if nodePrice.CPU == node.CPUCapacity && nodePrice.Memory == node.MemoryCapacity {
			fmt.Println("Node with matching details found")
			return *nodePrice, nil
		}

		// logrus.Printf("nodeprice: %v %v", nodePrice.CPU, nodePrice.Memory)
		if nodePrice.CPU >= node.CPUCapacity && nodePrice.Memory >= node.MemoryCapacity {
			if bestNP.NodePrice == nil || (nodePrice.CPU <= bestNP.CPU && nodePrice.Memory <= bestNP.Memory) {
				bestNP.CPU = nodePrice.CPU
				bestNP.CPUPrice = nodePrice.PricePerCPU
				bestNP.Memory = nodePrice.Memory
				bestNP.MemoryPrice = nodePrice.PricePerMemory
				bestNP.NodePrice = nodePrice
			}
		}
	}
	if bestNP.NodePrice == nil {
		logrus.Printf("no satisfying node price found")
		return *nodePrices[0], nil
	}
	return *bestNP.NodePrice, nil
}
