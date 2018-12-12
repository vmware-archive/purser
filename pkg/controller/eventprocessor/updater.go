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

	log "github.com/Sirupsen/logrus"

	groups_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	"github.com/vmware/purser/pkg/controller"
	"github.com/vmware/purser/pkg/controller/metrics"

	api_v1 "k8s.io/api/core/v1"
)

func updateCustomGroups(payloads []*interface{}, groups *groups_v1.GroupList) {
	for _, event := range payloads {
		payload := (*event).(*controller.Payload)
		if payload.ResourceType != "Pod" {
			continue
		}
		pod := api_v1.Pod{}
		err := json.Unmarshal([]byte(payload.Data), &pod)
		if err != nil {
			log.Errorf("error unmarshalling payload %s, %v", payload.Data, err)
		}

		log.Debugf("Started updating user created groups for pod %s update.", pod.Name)

		for _, group := range groups.Items {
			if isPodBelongsToGroup(group, &pod) {
				log.Infof("Updating the user group %s with pod %s details.", group.Spec.Name, pod.Name)
				updatePodDetails(group, pod, *payload)
			}
		}
		log.Debugf("Completed updating user created groups for pod %s update.", pod.Name)
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
		}
	} else if payload.EventType == controller.Update {
		// TODO: handle all pod updates

		// handle pod pvc updates
		*existingPodDetails = controller.UpdatePodVolumeClaims(pod, *existingPodDetails, payload.CaptureTime)
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
	containerMetrics := metrics.Metrics{}
	if cont.Resources.Requests != nil {
		containerMetrics.CPURequest = cont.Resources.Requests.Cpu()
		containerMetrics.MemoryRequest = cont.Resources.Requests.Memory()
	}
	if cont.Resources.Limits != nil {
		containerMetrics.CPULimit = cont.Resources.Limits.Cpu()
		containerMetrics.MemoryLimit = cont.Resources.Limits.Memory()
	}
	container.Metrics = &containerMetrics
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
