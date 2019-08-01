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

package main

import (
	"flag"
	"time"

	"github.com/vmware/purser/pkg/controller/dgraph/models/query"

	"github.com/vmware/purser/pkg/pricing"

	log "github.com/Sirupsen/logrus"

	"github.com/robfig/cron"
	"github.com/vmware/purser/cmd/controller/api"
	"github.com/vmware/purser/cmd/controller/config"
	"github.com/vmware/purser/pkg/controller"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/discovery/processor"
	"github.com/vmware/purser/pkg/controller/eventprocessor"
	"github.com/vmware/purser/pkg/utils"
)

var conf controller.Config

// InClusterConfigPath should be empty to get client and config for InCluster environment.
const InClusterConfigPath = ""

var interactions *string

func init() {
	logLevel := flag.String("log", "info", "set log level as info or debug")
	dgraphURL := flag.String("dgraphURL", "purser-db", "dgraph zero url")
	dgraphPort := flag.String("dgraphPort", "9080", "dgraph zero port")
	interactions = flag.String("interactions", "disable", "enable discovery of interactions")
	kubeconfig := flag.String("kubeconfig", InClusterConfigPath, "path to the kubeconfig file")
	flag.Parse()

	utils.InitializeLogger(*logLevel)
	config.Setup(&conf, *kubeconfig)

	// start dgraph and create login if not exists
	dgraph.Start(*dgraphURL, *dgraphPort)
	dgraph.StoreLogin()
}

func main() {
	go api.StartServer(conf)
	// go startCronJobForPopulatingRateCard()
	// time.Sleep(time.Minute * 5)
	pricing.TestRateCards()
	go eventprocessor.ProcessEvents(&conf)

	if *interactions == "enable" {
		go startInteractionsDiscovery()
	}
	go startCronJobForUpdatingCustomGroups()
	controller.Start(&conf)
}

// starts first discovery after 5 min of controller starting. Next runs will occur in every 59 min
func startInteractionsDiscovery() {
	time.Sleep(time.Minute * 5)
	runDiscovery()

	c := cron.New()
	err := c.AddFunc("@every 0h59m", runDiscovery)
	if err != nil {
		log.Error(err)
	}
	err = c.AddFunc("@daily", dgraph.RemoveResourcesInactive)
	if err != nil {
		log.Error(err)
	}
	c.Start()
}

func runDiscovery() {
	processor.ProcessPodInteractions(conf)
	processor.ProcessServiceInteractions(conf)
}

func startCronJobForUpdatingCustomGroups() {
	query.ComputeClusterAllocationAndCapacity()
	runGroupUpdate()

	c := cron.New()
	err := c.AddFunc("@every 0h5m", runGroupUpdate)
	if err != nil {
		log.Error(err)
	}
	err = c.AddFunc("@every 0h5m", query.ComputeClusterAllocationAndCapacity)
	if err != nil {
		log.Error(err)
	}
	c.Start()
}

func runGroupUpdate() {
	eventprocessor.UpdateGroups(conf.Groupcrdclient)
}

func startCronJobForPopulatingRateCard() {
	cloud := &pricing.Cloud{Kubeclient: conf.Kubeclient}
	// find cloud provider and region
	cloud.CloudProvider, cloud.Region = pricing.GetClusterProviderAndRegion()
	cloud.PopulateRateCard()

	c := cron.New()

	err := c.AddFunc("@every 168h", cloud.PopulateRateCard)
	if err != nil {
		log.Error(err)
	}
	c.Start()
}
