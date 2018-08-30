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
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/purser_controller/buffering"
	"github.com/vmware/purser/pkg/purser_controller/client"
	"github.com/vmware/purser/pkg/purser_controller/config"
	"github.com/vmware/purser/pkg/purser_controller/controller"
	"github.com/vmware/purser/pkg/purser_controller/eventprocessor"
)

var conf *config.Config

var groupcrdclient *client.GroupCrdClient
var subscriberclient *client.SubscriberCrdClient

func init() {
	setlogFile()
	conf = &config.Config{}
	conf.Resource = config.Resource{Pod: true, Node: true, PersistentVolume: true, PersistentVolumeClaim: true, ReplicaSet: true,
		Deployment: true, StatefulSet: true, DaemonSet: true, Job: true, Service: true}
	conf.RingBuffer = &buffering.RingBuffer{Size: buffering.BUFFER_SIZE, Mutex: &sync.Mutex{}}
	// initialize client for api extension server
	conf.Groupcrdclient, conf.Subscriberclient = controller.GetApiExtensionClient()
}

func setlogFile() {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	go eventprocessor.ProcessEvents(conf)
	controller.Start(conf)
}
