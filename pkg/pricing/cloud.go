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
	switch c.CloudProvider {
	case models.AWS:
		rateCard := aws.GetRateCardForAWS(c.Region)
		models.StoreRateCard(rateCard)
	}
}
