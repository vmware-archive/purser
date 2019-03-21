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

package eventprocessor

import (
	"encoding/json"
	"time"

	"github.com/vmware/purser/pkg/controller/dgraph/models/query"

	log "github.com/Sirupsen/logrus"

	groups_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	subcriber_v1 "github.com/vmware/purser/pkg/apis/subscriber/v1"
	"github.com/vmware/purser/pkg/controller"
	"github.com/vmware/purser/pkg/controller/dgraph/models"

	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	batch_v1 "k8s.io/api/batch/v1"
	api_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
)

// ProcessEvents processes the event and notifies the subscribers.
func ProcessEvents(conf *controller.Config) {

	for {
		conf.RingBuffer.PrintDetails()

		for {
			data, size := conf.RingBuffer.ReadN(ReadSize)

			if size == 0 {
				log.Debug("No new events to process.")
				break
			}

			ProcessPayloads(data, conf)

			subscribers, err := query.RetrieveSubscribers()
			if err == nil {
				notifySubscribers(data, subscribers)
			} else {
				log.Errorf("unable to retrieve subscribers from dgraph: %v", err)
			}

			conf.RingBuffer.RemoveN(size)
			conf.RingBuffer.PrintDetails()
		}
		time.Sleep(10 * time.Second)
	}
}

// ProcessPayloads store payload info in dgraph. If payload is of type group then it updates its group spec
func ProcessPayloads(payloads []*interface{}, conf *controller.Config) {
	for _, event := range payloads {
		payload := (*event).(*controller.Payload)
		handlePayloadBasedOnResource(payload, conf)
	}
}

// nolint: gocyclo
func handlePayloadBasedOnResource(payload *controller.Payload, conf *controller.Config) {
	var err error
	switch payload.ResourceType {
	case "Pod":
		pod := api_v1.Pod{}
		unmarshalPayload(payload, &pod)
		err = models.StorePod(pod)
	case "Service":
		service := api_v1.Service{}
		unmarshalPayload(payload, &service)
		err = models.StoreService(service)
	case "Node":
		node := api_v1.Node{}
		unmarshalPayload(payload, &node)
		_, err = models.StoreNode(node)
	case "Namespace":
		ns := api_v1.Namespace{}
		unmarshalPayload(payload, &ns)
		_, err = models.StoreNamespace(ns)
	case "Deployment":
		deployment := apps_v1beta1.Deployment{}
		unmarshalPayload(payload, &deployment)
		_, err = models.StoreDeployment(deployment)
	case "ReplicaSet":
		replicaset := ext_v1beta1.ReplicaSet{}
		unmarshalPayload(payload, &replicaset)
		_, err = models.StoreReplicaset(replicaset)
	case "StatefulSet":
		statefulset := apps_v1beta1.StatefulSet{}
		unmarshalPayload(payload, &statefulset)
		_, err = models.StoreStatefulset(statefulset)
	case "PersistentVolume":
		pv := api_v1.PersistentVolume{}
		unmarshalPayload(payload, &pv)
		_, err = models.StorePersistentVolume(pv, conf.Kubeclient)
	case "PersistentVolumeClaim":
		pvc := api_v1.PersistentVolumeClaim{}
		unmarshalPayload(payload, &pvc)
		_, err = models.StorePersistentVolumeClaim(pvc)
	case "DaemonSet":
		daemonset := ext_v1beta1.DaemonSet{}
		unmarshalPayload(payload, &daemonset)
		_, err = models.StoreDaemonset(daemonset)
	case "Job":
		job := batch_v1.Job{}
		unmarshalPayload(payload, &job)
		_, err = models.StoreJob(job)
	case "Group":
		groupCRD := &groups_v1.Group{}
		unmarshalPayload(payload, &groupCRD)
		handlePayloadForGroup(payload, conf, groupCRD.Name)
	case "Subscriber":
		subscriberCRD := subcriber_v1.Subscriber{}
		unmarshalPayload(payload, &subscriberCRD)
		_, err = models.StoreSubscriberCRD(subscriberCRD)
	}
	checkDgraphError(payload.ResourceType, err)
}

func unmarshalPayload(payload *controller.Payload, resource interface{}) {
	err := json.Unmarshal([]byte(payload.Data), resource)
	if err != nil {
		log.Errorf("Error un marshalling payload " + payload.Data)
	}
}

func checkDgraphError(resource string, err error) {
	if err != nil {
		log.Errorf("Error while persisting %s %v", resource, err)
	}
}

func handlePayloadForGroup(payload *controller.Payload, conf *controller.Config, groupName string) {
	if payload.EventType == controller.Delete {
		models.DeleteGroup(groupName)
	} else {
		group, err := conf.Groupcrdclient.Get(groupName)
		if err != nil {
			log.Errorf("Unable to get group from client: (%v)", err)
		} else {
			UpdateGroup(group, conf.Groupcrdclient)
		}
	}
}
