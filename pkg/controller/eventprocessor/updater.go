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
	"github.com/vmware/purser/pkg/controller/discovery/processor"
	"k8s.io/apimachinery/pkg/labels"

	log "github.com/Sirupsen/logrus"

	groups_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	"github.com/vmware/purser/pkg/controller"
	"github.com/vmware/purser/pkg/controller/metrics"

	api_v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func updateCustomGroups(payloads []*interface{}, groups *groups_v1.GroupList, conf *controller.Config) {
	for _, event := range payloads {
		payload := (*event).(*controller.Payload)
		if payload.ResourceType == "Pod" {
			pod := api_v1.Pod{}
			err := json.Unmarshal([]byte(payload.Data), &pod)
			if err != nil {
				log.Errorf("error un marshalling payload %s, %v", payload.Data, err)
				return
			}

			log.Debugf("Started updating user created groups for pod %s update.", pod.Name)

			for _, group := range groups.Items {
				if isPodBelongsToGroup(group, &pod) {
					log.Infof("Updating the user group %s with pod %s details.", group.Spec.Name, pod.Name)
					updatePodDetails(group, pod, *payload)
				}
			}
			log.Infof("Completed updating user created groups for pod %s update.", pod.Name)
		} else if payload.ResourceType == "Group" {
			group := getGroupFromPayload(payload, conf)
			syncNewGroup(group, conf)
		}
	}
}

// nolint
func updatePodDetails(group *groups_v1.Group, pod api_v1.Pod, payload controller.Payload) {
	podKey := pod.GetObjectMeta().GetNamespace() + ":" + pod.GetObjectMeta().GetName()
	podDetails := group.Spec.PodsDetails

	if podDetails == nil {
		podDetails = map[string]*groups_v1.PodDetails{}
	}

	existingPodDetails := podDetails[podKey]
	if existingPodDetails != nil {
		if payload.EventType == controller.Create {
			// TODO:
			// This case means we have lost a Delete event for this pod. So we need to update
			// the pod details with the new one
		} else if payload.EventType == controller.Delete {
			existingPodDetails.EndTime = *pod.GetDeletionTimestamp()
			controller.PvcHandlePodDeletion(existingPodDetails)
		} else if payload.EventType == controller.Update {
			// TODO: handle all pod updates

			// handle pod pvc updates
			*existingPodDetails = controller.UpdatePodVolumeClaims(pod, *existingPodDetails, payload.CaptureTime)
		}
	} else {
		if payload.EventType == controller.Create {
			newPodDetails := groups_v1.PodDetails{Name: pod.Name, StartTime: pod.GetCreationTimestamp()}
			containers := []*groups_v1.Container{}
			for _, cont := range pod.Spec.Containers {
				container := getContainerWithMetrics(cont)
				containers = append(containers, container)
			}
			newPodDetails.Containers = containers
			newPodDetails = controller.UpdatePodVolumeClaims(pod, newPodDetails, pod.GetCreationTimestamp())
			podDetails[podKey] = &newPodDetails
		} else if payload.EventType == controller.Delete {
			// TODO:
			// This case means we have lost a Create event for this pod.
			// If we can retrieve pod details(metrics and creation time) then we can
			// include that in podDetails
		} else if payload.EventType == controller.Update {
			// TODO:
			// This case means we have lost a Create event for this pod.
			// We can retrieve pod details(metrics and creation time) then we can
			// include that in podDetails
		}
	}
	group.Spec.PodsDetails = podDetails
}

func getContainerWithMetrics(cont api_v1.Container) *groups_v1.Container {
	container := groups_v1.Container{Name: cont.Name}
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

func isPodBelongsToGroup(group *groups_v1.Group, pod *api_v1.Pod) bool {
	for groupLabelKey, groupLabelVal := range group.Spec.Labels {
		for podLabelKey, podLabelVal := range pod.Labels {
			if groupLabelKey == podLabelKey && groupLabelVal == podLabelVal {
				return true
			}
		}
	}
	return false
}

func getGroupFromPayload(payload *controller.Payload, conf *controller.Config) *groups_v1.Group {
	// convert the payload into Group object
	groupCRD := groups_v1.Group{}
	err := json.Unmarshal([]byte(payload.Data), &groupCRD)
	if err != nil {
		log.Errorf("error un marshalling payload " + payload.Data)
		return nil
	}

	group, err := conf.Groupcrdclient.Get(groupCRD.Name)
	if err != nil {
		log.Errorf("cannot update group details. Reason: Unable to retrieve group from client: (%v)", err)
		return nil
	}
	return group
}

func syncNewGroup(group *groups_v1.Group, conf *controller.Config) {
	podDetails := group.Spec.PodsDetails
	if podDetails == nil {
		podDetails = map[string]*groups_v1.PodDetails{}
	}

	// get all live pods with matching labels of group
	podList := processor.RetrievePodList(conf.Kubeclient, metav1.ListOptions{LabelSelector: labels.SelectorFromSet(group.Spec.Labels).String()})
	// add pod details in group spec
	for _, pod := range podList.Items {
		podKey := pod.GetObjectMeta().GetNamespace() + ":" + pod.GetObjectMeta().GetName()
		newPodDetails := groups_v1.PodDetails{Name: pod.Name, StartTime: pod.GetCreationTimestamp()}
		containers := []*groups_v1.Container{}
		for _, cont := range pod.Spec.Containers {
			container := getContainerWithMetrics(cont)
			containers = append(containers, container)
		}
		newPodDetails.Containers = containers
		newPodDetails = controller.UpdatePodVolumeClaims(pod, newPodDetails, pod.GetCreationTimestamp())
		podDetails[podKey] = &newPodDetails
	}
	group.Spec.PodsDetails = podDetails
}
