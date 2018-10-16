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
	groups_v1 "github.com/vmware/purser/pkg/client/clientset/typed/groups/v1"
	subscriber_v1 "github.com/vmware/purser/pkg/client/clientset/typed/subscriber/v1"
	"github.com/vmware/purser/pkg/controller/buffering"
	"k8s.io/client-go/kubernetes"
)

// These are the event types supported for controllers
const (
	Create = "create"
	Delete = "delete"
	Update = "update"
)

// Resource contains resource configuration
type Resource struct {
	Pod                   bool `json:"po"`
	Node                  bool `json:"node"`
	PersistentVolume      bool `json:"pv"`
	PersistentVolumeClaim bool `json:"pvc"`
	Service               bool `json:"service"`
	ReplicaSet            bool `json:"replicaset"`
	StatefulSet           bool `json:"statefulset"`
	Deployment            bool `json:"deployment"`
	Job                   bool `json:"job"`
	DaemonSet             bool `json:"daemonset"`
}

// Config contains config objects
type Config struct {
	Resource         Resource `json:"resource"`
	RingBuffer       *buffering.RingBuffer
	Groupcrdclient   *groups_v1.GroupClient
	Subscriberclient *subscriber_v1.SubscriberClient
	Kubeclient       *kubernetes.Clientset
}
