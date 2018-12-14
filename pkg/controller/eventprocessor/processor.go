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
	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
	"time"

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
// nolint: gocyclo
func ProcessPayloads(payloads []*interface{}, conf *controller.Config) {
	for _, event := range payloads {
		payload := (*event).(*controller.Payload)
		if payload.ResourceType == "Pod" {
			pod := api_v1.Pod{}
			err := json.Unmarshal([]byte(payload.Data), &pod)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			err = models.StorePod(pod)
			if err != nil {
				log.Errorf("Error while persisting pod %v", err)
			}
		} else if payload.ResourceType == "Service" {
			service := api_v1.Service{}
			err := json.Unmarshal([]byte(payload.Data), &service)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			err = models.StoreService(service)
			if err != nil {
				log.Errorf("Error while persisting service %v", err)
			}
		} else if payload.ResourceType == "Node" {
			node := api_v1.Node{}
			err := json.Unmarshal([]byte(payload.Data), &node)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			_, err = models.StoreNode(node)
			if err != nil {
				log.Errorf("Error while persisting node %v", err)
			}
		} else if payload.ResourceType == "Namespace" {
			ns := api_v1.Namespace{}
			err := json.Unmarshal([]byte(payload.Data), &ns)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			_, err = models.StoreNamespace(ns)
			if err != nil {
				log.Errorf("Error while persisting namespace %v", err)
			}
		} else if payload.ResourceType == "Deployment" {
			deployment := apps_v1beta1.Deployment{}
			err := json.Unmarshal([]byte(payload.Data), &deployment)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			_, err = models.StoreDeployment(deployment)
			if err != nil {
				log.Errorf("Error while persisting deployment %v", err)
			}
		} else if payload.ResourceType == "ReplicaSet" {
			replicaset := ext_v1beta1.ReplicaSet{}
			err := json.Unmarshal([]byte(payload.Data), &replicaset)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			_, err = models.StoreReplicaset(replicaset)
			if err != nil {
				log.Errorf("Error while persisting replicaset %v", err)
			}
		} else if payload.ResourceType == "StatefulSet" {
			statefulset := apps_v1beta1.StatefulSet{}
			err := json.Unmarshal([]byte(payload.Data), &statefulset)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			_, err = models.StoreStatefulset(statefulset)
			if err != nil {
				log.Errorf("Error while persisting statefulset %v", err)
			}
		} else if payload.ResourceType == "PersistentVolume" {
			pv := api_v1.PersistentVolume{}
			err := json.Unmarshal([]byte(payload.Data), &pv)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			_, err = models.StorePersistentVolume(pv)
			if err != nil {
				log.Errorf("Error while persisting persistent volume %v", err)
			}
		} else if payload.ResourceType == "PersistentVolumeClaim" {
			pvc := api_v1.PersistentVolumeClaim{}
			err := json.Unmarshal([]byte(payload.Data), &pvc)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			_, err = models.StorePersistentVolumeClaim(pvc)
			if err != nil {
				log.Errorf("Error while persisting persistent volume claim %v", err)
			}
		} else if payload.ResourceType == "DaemonSet" {
			daemonset := ext_v1beta1.DaemonSet{}
			err := json.Unmarshal([]byte(payload.Data), &daemonset)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			_, err = models.StoreDaemonset(daemonset)
			if err != nil {
				log.Errorf("Error while persisting daemonset %v", err)
			}
		} else if payload.ResourceType == "Job" {
			job := batch_v1.Job{}
			err := json.Unmarshal([]byte(payload.Data), &job)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			_, err = models.StoreJob(job)
			if err != nil {
				log.Errorf("Error while persisting job %v", err)
			}
		} else if payload.ResourceType == "Group" {
			groupCRD := groups_v1.Group{}
			err := json.Unmarshal([]byte(payload.Data), &groupCRD)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			} else {
				group, err := conf.Groupcrdclient.Get(groupCRD.Name)
				if err != nil {
					log.Errorf("Unable to get group from client: (%v)", err)
				}
				UpdateGroup(group, conf.Groupcrdclient)
			}
		} else if payload.ResourceType == "Subscriber" {
			subscriberCRD := subcriber_v1.Subscriber{}
			err := json.Unmarshal([]byte(payload.Data), &subscriberCRD)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			_, err = models.StoreSubscriberCRD(subscriberCRD)
			if err != nil {
				log.Errorf("Error while persisting subscriber CRD %v", err)
			}
		}
	}
}
