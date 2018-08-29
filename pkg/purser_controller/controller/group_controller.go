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
	"flag"
	"time"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/purser_controller/client"
	"github.com/vmware/purser/pkg/purser_controller/crd"
	"github.com/vmware/purser/pkg/purser_controller/eventprocessor"
	"github.com/vmware/purser/pkg/purser_controller/metrics"
	api_v1 "k8s.io/api/core/v1"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const environment = "dev"

// return rest config, if path not specified assume in cluster config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	log.Println("Using In cluster config.")
	return rest.InClusterConfig()
}

func GetApiExtensionClient() (*client.GroupCrdClient, *client.SubscriberCrdClient) {
	var config *rest.Config
	var err error
	if environment == "dev" {
		kubeconf := flag.String("kubeconf", "/Users/gurusreekanthc/.kube/config", "path to Kubernetes config file")
		flag.Parse()
		config, err = GetClientConfig(*kubeconf)
	} else {
		config, err = GetClientConfig("")
	}

	if err != nil {
		log.Println(err)
		panic(err.Error())
	}

	// create clientset and create our CRD, this only need to run once
	clientset, err := apiextcs.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// note: if the CRD exist our CreateGroupCRD function is set to exit without an error
	err = crd.CreateGroupCRD(clientset)
	if err != nil {
		panic(err)
	}

	err = crd.CreateSubscriberCRD(clientset)
	if err != nil {
		panic(err)
	}

	// Wait for the CRD to be created before we use it (only needed if its a new one)
	time.Sleep(3 * time.Second)

	// Create a new clientset which include our CRD schema
	gcrdcs, gscheme, err := crd.NewGroupClient(config)
	if err != nil {
		panic(err)
	}

	// Create a CRD client interface
	groupcrdclient := client.CreateGroupCrdClient(gcrdcs, gscheme, "default")

	// Create a new clientset which include our CRD schema
	crdcs, scheme, err := crd.NewSubscriberClient(config)
	if err != nil {
		panic(err)
	}

	// Create a CRD client interface
	subcrdclient := client.CreateSubscriberCrdClient(crdcs, scheme, "default")

	return groupcrdclient, subcrdclient
}

func CreateGroupCRDInstance(crdclient *client.GroupCrdClient, groupName string, groupType string) *crd.Group {
	// Create a new Example object and write to k8s
	example := &crd.Group{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: groupName,
			//Labels: map[string]string{"mylabel": "test"},
		},
		Spec: crd.GroupSpec{
			Name: groupName,
			Type: groupType,
		},
		Status: crd.GroupStatus{
			State:   "created",
			Message: "Done",
		},
	}

	result, err := crdclient.CreateGroup(example)
	if err == nil {
		log.Printf("Created Group : %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		log.Printf("Group already exists : %#v\n", result)
	} else {
		panic(err)
	}
	return result
}

func CreateSubscriberCRDInstance(crdclient *client.SubscriberCrdClient, subscriberName string) *crd.Subscriber {
	// Create a new Example object and write to k8s
	example := &crd.Subscriber{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: subscriberName,
			//Labels: map[string]string{"mylabel": "test"},
		},
		Spec: crd.SubscriberSpec{
			Name: subscriberName,
		},
		Status: crd.SubscriberStatus{
			State:   "created",
			Message: "Done",
		},
	}

	result, err := crdclient.CreateSubscriber(example)
	if err == nil {
		log.Printf("Created Subscriber : %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		log.Printf("Subscriber already exists : %#v\n", result)
	} else {
		panic(err)
	}
	return result
}

func ListGroupCrdInstances(crdclient *client.GroupCrdClient) {
	items, err := crdclient.ListGroups(meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("List:\n%s\n", items)
}

func ListSubscriberCrdInstances(crdclient *client.SubscriberCrdClient) {
	items, err := crdclient.ListSubscriber(meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("List:\n%s\n", items)
}

func GetGroupCrdByName(crdclient *client.GroupCrdClient, groupName string, groupType string) *crd.Group {
	group, err := crdclient.GetGroup(groupName)

	if err == nil {
		return group
	} else if apierrors.IsNotFound(err) {
		// create group if not exist
		return CreateGroupCRDInstance(crdclient, groupName, groupType)
	} else {
		panic(err)
	}
}

func GetAllGroups(crdclient *client.GroupCrdClient) []crd.Group {
	items, err := crdclient.ListGroups(meta_v1.ListOptions{})
	if err != nil {
		log.Error("Error while fetching groups ", err)
		return nil
	}
	return items.Items
}

func UpdateCustomGroups(crdclient *client.GroupCrdClient, payloads []*interface{}) {
	groups := GetAllGroups(crdclient)
	if groups == nil {
		return
	}

	for _, event := range payloads {

		payload := (*event).(eventprocessor.Payload)
		if payload.ResourceType != "Pod" {
			continue
		}
		pod := api_v1.Pod{}
		json.Unmarshal([]byte(payload.Data), &pod)

		log.Info("Started updating User Created Groups for pod {} update.\n", pod.Name)

		for _, group := range groups {
			if isPodBelongsToGroup(&group, &pod) {
				log.Info("Updating the user group {} with pod {} details\n", group.Spec.Name, pod.Name)

				podKey := pod.GetObjectMeta().GetNamespace() + ":" + pod.GetObjectMeta().GetName()
				podDetails := group.Spec.PodsDetails

				if podDetails == nil {
					podDetails = map[string]*crd.PodDetails{}
				}

				existingPodDetails := podDetails[podKey]
				if existingPodDetails != nil {
					if payload.EventType == "create" {
						// do nothing.
					} else if payload.EventType == "delete" {
						existingPodDetails.EndTime = pod.GetObjectMeta().GetCreationTimestamp()
					}
				} else {
					if payload.EventType == "create" {
						newPodDetails := crd.PodDetails{Name: pod.Name, StartTime: pod.GetCreationTimestamp()}
						containers := []*crd.Container{}
						for _, cont := range pod.Spec.Containers {
							container := crd.Container{Name: cont.Name}
							metrics := metrics.Metrics{}
							if cont.Resources.Requests != nil {
								metrics.CpuRequest = cont.Resources.Requests.Cpu()
								metrics.MemoryRequest = cont.Resources.Requests.Memory()
							}
							if cont.Resources.Limits != nil {
								metrics.CpuLimit = cont.Resources.Limits.Cpu()
								metrics.MemoryLimit = cont.Resources.Limits.Memory()
							}
							container.Metrics = &metrics
							containers = append(containers, &container)
						}
						newPodDetails.Containers = containers
						podDetails[podKey] = &newPodDetails
					} else if payload.EventType == "delete" {
						// do nothing.
					}
				}
			}
		}
		log.Info("Completed updating User Created Groups for pod {} update.\n", pod.Name)
	}

	// update all the groups
	for _, group := range groups {
		_, err := crdclient.UpdateGroup(&group)

		if err != nil {
			log.Error("There is an error while updating the crd for group = %s\n", group.Name)
		} else {
			log.Info("Updating the crd for group = %s is successful\n", group.Name)
		}
	}
}

func isPodBelongsToGroup(group *crd.Group, pod *api_v1.Pod) bool {
	for groupLabelKey, groupLabelVal := range group.Spec.Labels {
		for podLabelKey, podLabelVal := range pod.Labels {
			if groupLabelKey == podLabelKey && groupLabelVal == podLabelVal {
				return true
			}
		}
	}
	return false
}

func calculatedAggregatedPodMetric(met map[string]*metrics.Metrics) *metrics.Metrics {
	cpuLimit := &resource.Quantity{}
	memoryLimit := &resource.Quantity{}
	cpuRequest := &resource.Quantity{}
	memoryRequest := &resource.Quantity{}
	for _, c := range met {
		cpuLimit.Add(*c.CpuLimit)
		memoryLimit.Add(*c.MemoryLimit)
		cpuRequest.Add(*c.CpuRequest)
		memoryRequest.Add(*c.MemoryRequest)
	}
	return &metrics.Metrics{
		CpuLimit:      cpuLimit,
		MemoryLimit:   memoryLimit,
		CpuRequest:    cpuRequest,
		MemoryRequest: memoryRequest,
	}
}
