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

package pricing

import (
	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
	"github.com/vmware/purser/pkg/pricing/aws"
	"github.com/vmware/purser/pkg/pricing/azure"
	"github.com/vmware/purser/pkg/pricing/gcp"
	"github.com/vmware/purser/pkg/pricing/pks"
	"k8s.io/client-go/kubernetes"
)

// Cloud structure used for pricing
type Cloud struct {
	CloudProvider string
	Region        string
	Kubeclient    *kubernetes.Clientset
}

// GetClusterProviderAndRegion returns cluster provider(ex: aws) and region(ex: us-east-1)
func GetClusterProviderAndRegion() (string, string) {
	// TODO: https://github.com/vmware/purser/issues/143
	cloudProvider := models.AWS
	region := "us-east-1"

	logrus.Infof("CloudProvider: %s, Region: %s", cloudProvider, region)
	return cloudProvider, region
}

// PopulateRateCard given a cloud (cloudProvider and region) it populates corresponding rate card in dgraph
func (c *Cloud) PopulateRateCard() {
	PopulateAllRateCards()
}

// TestRateCards ...
func TestRateCards() {
	// costs := models.GetCost(models.CloudRegionInfo{
	// 	CloudRegions: []models.CloudRegion{
	// 		models.CloudRegion{Region: "us-east-1", CloudProvider: models.AWS},
	// 		models.CloudRegion{Region: "westus", CloudProvider: models.AZURE},
	// 		models.CloudRegion{Region: "us-east1", CloudProvider: models.GCP},
	// 		models.CloudRegion{Region: "US-East-1", CloudProvider: models.PKS},
	// 	}})
	// fmt.Printf("%#v", costs)
}

//PopulateAllRateCards take region as input and saves the rate card for all cloud providers
func PopulateAllRateCards() {
	awsRegions := []string{"us-east-1", "us-west-1"}
	azureRegions := []string{"eastus", "westus"}
	gcpRegions := []string{"us-east1", "us-west1"}
	pksRegions := []string{"US-East-1", "US-West-2"}
	for _, region := range awsRegions {
		go models.StoreRateCard(aws.GetRateCardForAWS(region))
	}
	for _, region := range azureRegions {
		go models.StoreRateCard(azure.GetRateCardForAzure(region))
	}
	for _, region := range gcpRegions {
		go models.StoreRateCard(gcp.GetRateCardForGCP(region))
	}
	for _, region := range pksRegions {
		go models.StoreRateCard(pks.GetRateCardForPKS(region))
	}
}

//PopulateRateCard ...
func PopulateRateCard(region string, cloudProvider string) {
	var rateCard *models.RateCard
	switch cloudProvider {
	case models.AWS:
		rateCard = aws.GetRateCardForAWS(region)
		models.StoreRateCard(rateCard)
	case models.GCP:
		rateCard = gcp.GetRateCardForGCP(region)
		if rateCard != nil {
			models.StoreRateCard(rateCard)
		}
	case models.AZURE:
		rateCard := azure.GetRateCardForAzure(region)
		models.StoreRateCard(rateCard)
	case models.PKS:
		rateCard := pks.GetRateCardForPKS(region)
		models.StoreRateCard(rateCard)
	}

}
