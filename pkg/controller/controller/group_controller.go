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
	"encoding/json"
	"flag"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/client"
	"github.com/vmware/purser/pkg/controller/crd"

	"github.com/vmware/purser/pkg/controller/metrics"
	api_v1 "k8s.io/api/core/v1"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const environment = "prod"

// GetClientConfig returns rest config, if path not specified assume in cluster config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	log.Println("Using In cluster config.")
	return rest.InClusterConfig()
}

// GetAPIExtensionClient returns an instance of CRD client.
func GetAPIExtensionClient() (*client.GroupCrdClient, *client.SubscriberCrdClient) {
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

// CreateGroupCRDInstance creates group CRD instance.
func CreateGroupCRDInstance(crdclient *client.GroupCrdClient, groupName string, groupType string) *crd.Group {
	// Create a new Example object and write to k8s
	example := &crd.Group{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: groupName,
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

// CreateSubscriberCRDInstance creates subscriber CRD instance.
func CreateSubscriberCRDInstance(crdclient *client.SubscriberCrdClient, subscriberName string) *crd.Subscriber {
	// Create a new Example object and write to k8s
	example := &crd.Subscriber{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: subscriberName,
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

// ListGroupCrdInstances fetches list of Group CRD instances.
func ListGroupCrdInstances(crdclient *client.GroupCrdClient) {
	items, err := crdclient.ListGroups(meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("List:\n%v\n", items)
}

// ListSubscriberCrdInstances fetches list of subscriber CRD instances.
func ListSubscriberCrdInstances(crdclient *client.SubscriberCrdClient) {
	items, err := crdclient.ListSubscriber(meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("List:\n%v\n", items)
}

// GetGroupCrdByName return group CRD by name.
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

// GetAllGroups returns the collection of all groups.
func GetAllGroups(crdclient *client.GroupCrdClient) []*crd.Group {
	items, err := crdclient.ListGroups(meta_v1.ListOptions{})
	if err != nil {
		log.Error("Error while fetching groups ", err)
		return nil
	}
	return items.Items
}

// UpdateCustomGroups modifies custom group deifinitions.
func UpdateCustomGroups(payloads []*interface{}, groups []*crd.Group, crdclient *client.GroupCrdClient) {

	processPayload(groups, payloads)

	// update all the groups
	for _, group := range groups {
		_, err := crdclient.UpdateGroup(group)

		if err != nil {
			log.Errorf("There is an error while updating the crd for group = "+group.Name, err)
		} else {
			log.Debug("Updating the crd for group = " + group.Name + " is successful")
		}
	}
}

func processPayload(groups []*crd.Group, payloads []*interface{}) {
	for _, event := range payloads {
		payload := (*event).(*Payload)
		if payload.ResourceType != "Pod" {
			continue
		}
		pod := api_v1.Pod{}
		err := json.Unmarshal([]byte(payload.Data), &pod)
		if err != nil {
			log.Errorf("Error unmarshalling payload " + payload.Data)
		}

		log.Info("Started updating User Created Groups for pod "+pod.Name+" update.", pod.Name)

		for _, group := range groups {
			if isPodBelongsToGroup(group, &pod) {
				log.Info("Updating the user group " + group.Spec.Name + " with pod " + pod.Name + " details")
				updatePodDetails(group, pod, *payload)
			}
		}
		log.Debug("Completed updating User Created Groups for pod " + pod.Name + " update.")
	}
}

// nolint
func updatePodDetails(group *crd.Group, pod api_v1.Pod, payload Payload) {
	podKey := pod.GetObjectMeta().GetNamespace() + ":" + pod.GetObjectMeta().GetName()
	podDetails := group.Spec.PodsDetails

	if podDetails == nil {
		podDetails = map[string]*crd.PodDetails{}
	}

	existingPodDetails := podDetails[podKey]
	if existingPodDetails != nil {
		if payload.EventType == Create {
			// TODO:
			// This case means we have lost a Delete event for this pod. So we need to update
			// the pod details with the new one
		} else if payload.EventType == Delete {
			// Here we are not using pod.GetObjectMeta().GetDeletionTimestamp() because
			// by the time controller gets to this part of the code the object(pod) might have been
			// removed from etcd.
			existingPodDetails.EndTime = *pod.GetDeletionTimestamp()
		}
	} else {
		if payload.EventType == Create {
			newPodDetails := crd.PodDetails{Name: pod.Name, StartTime: pod.GetCreationTimestamp()}
			containers := []*crd.Container{}
			for _, cont := range pod.Spec.Containers {
				container := getContainerWithMetrics(cont)
				containers = append(containers, container)
			}
			newPodDetails.Containers = containers
			podDetails[podKey] = &newPodDetails
		} else if payload.EventType == Delete {
			// TODO:
			// This case means we have lost a Create event for this pod.
			// If we can retrieve pod details(metrics and creation time) then we can
			// include that in podDetails
		}
	}
	group.Spec.PodsDetails = podDetails
}

func getContainerWithMetrics(cont api_v1.Container) *crd.Container {
	container := crd.Container{Name: cont.Name}
	metrics := metrics.Metrics{}
	if cont.Resources.Requests != nil {
		metrics.CPURequest = cont.Resources.Requests.Cpu()
		metrics.MemoryRequest = cont.Resources.Requests.Memory()
	}
	if cont.Resources.Limits != nil {
		metrics.CPULimit = cont.Resources.Limits.Cpu()
		metrics.MemoryLimit = cont.Resources.Limits.Memory()
	}
	container.Metrics = &metrics
	return &container
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
