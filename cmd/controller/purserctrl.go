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
	"time"

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

func init() {
	utils.InitializeLogger()
	config.Setup(&conf)
}

func main() {
	go api.StartServer()
	go eventprocessor.ProcessEvents(&conf)
	go startCronJobs()
	controller.Start(&conf)
}

// starts first discovery after 10 min of controller starting. Next runs will occur in every 30 min
func startCronJobs() {
	time.Sleep(time.Minute * 5)
	runDiscovery()

	c := cron.New()
	err := c.AddFunc("@every 0h59m", runDiscovery)
	if err != nil {
		log.Fatal(err)
	}
	err = c.AddFunc("@daily", dgraph.RemoveResourcesInactiveInCurrentMonth)
	if err != nil {
		log.Error(err)
	}
	c.Start()
}

func runDiscovery() {
	processor.ProcessPodInteractions(conf)
	processor.ProcessServiceInteractions(conf)
}
