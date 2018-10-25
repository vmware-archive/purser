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

package controller

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/client"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	subscriber_v1 "github.com/vmware/purser/pkg/client/clientset/typed/subscriber/v1"
)

// TestCrdFlow executes the CRD flow.
func TestCrdFlow(t *testing.T) {
	clientset, clusterConfig := client.GetAPIExtensionClient("")
	subcrdclient := subscriber_v1.NewSubscriberClient(clientset, clusterConfig)
	ListSubscriberCrdInstances(subcrdclient)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

// ListSubscriberCrdInstances fetches list of subscriber CRD instances.
func ListSubscriberCrdInstances(crdclient *subscriber_v1.SubscriberClient) {
	items, err := crdclient.ListSubscriber(meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("List:\n%v\n", items)
}
