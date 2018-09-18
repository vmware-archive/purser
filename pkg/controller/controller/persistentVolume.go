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
	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/crd"
	"github.com/vmware/purser/pkg/controller/utils"
	api_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// activePodVolumeClaims: current active(bounded) podVolumeClaims for the pod.
// old podVolumeClaims: Claims(map) for the pod before this update.
// compares old podVolumeClaims and activePodVolumeClaims. This function hadles 3 cases.
// Case Unbound pvc:
//		Present as 'bounded' in old podVolumeClaims. Not present in activePodVolumeClaims.
// Case New PVC:
//		Not present in old podVolumeClaims. Present in activePodVolumeClaims.
// Case Bound an unbounded pvc:
//		Present as 'unbounded' in old podVolumClaimss. Present in activePodVolumeClaims.
func updatePodVolumeClaims(pod api_v1.Pod, podDetails crd.PodDetails, eventTime meta_v1.Time) crd.PodDetails {
	activePodVolumeClaims := getactivePodVolumeClaims(pod)

	podVolumeClaims := podDetails.PodVolumeClaims
	if podVolumeClaims == nil {
		podVolumeClaims = map[string]*crd.PersistentVolumeClaim{}
	}

	for claimName := range podVolumeClaims {
		// isPresent: true if pvc is present in activePodVolumeClaims.
		_, isPresent := activePodVolumeClaims[claimName]

		// isBounded: true if pvc is present in old podVolumeClaims as 'bounded'
		isBounded := checkBounded(podVolumeClaims[claimName])

		if (!isPresent) && isBounded {
			// Case Unbound pvc
			log.Info("Unbounded pvc: " + claimName + " from pod: " + podDetails.Name)
			podVolumeClaims[claimName].UnboundTimes = append(podVolumeClaims[claimName].UnboundTimes, eventTime)
		} else if isPresent && (!isBounded) {
			// Case Bound an unbounded pvc
			log.Info("Bounded pvc: " + claimName + " to pod: " + podDetails.Name)
			podVolumeClaims[claimName].BoundTimes = append(podVolumeClaims[claimName].BoundTimes, eventTime)
		}
	}

	// check for new pvc
	for claimKey := range activePodVolumeClaims {
		_, isPresent := podVolumeClaims[claimKey]
		if !isPresent {
			// Case New PVC
			log.Info("Bounded new pvc: " + claimKey + "to pod: " + podDetails.Name)
			podVolumeClaims[claimKey] = activePodVolumeClaims[claimKey]
		}
	}

	// TODO: handle Case Resizing of PVC

	podDetails.PodVolumeClaims = podVolumeClaims
	return podDetails
}

func getactivePodVolumeClaims(pod api_v1.Pod) map[string]*crd.PersistentVolumeClaim {
	namespace := pod.GetNamespace()
	podVolumeClaims := map[string]*crd.PersistentVolumeClaim{}
	for j := 0; j < len(pod.Spec.Volumes); j++ {
		vol := pod.Spec.Volumes[j]
		if vol.PersistentVolumeClaim != nil {
			claimName := vol.PersistentVolumeClaim.ClaimName
			podVolumeClaims[claimName] = collectPersistentVolumeClaim(claimName, namespace)
			podVolumeClaims[claimName].BoundTimes = append(podVolumeClaims[claimName].BoundTimes, pod.GetCreationTimestamp())
		}
	}
	return podVolumeClaims
}

func collectPersistentVolumeClaim(claimName, namespace string) *crd.PersistentVolumeClaim {
	pvc, err := Kubeclient.CoreV1().PersistentVolumeClaims(namespace).Get(claimName, meta_v1.GetOptions{})
	if isPVCError(err, claimName) {
		return nil
	}

	request := pvc.Spec.Resources.Requests["storage"].DeepCopy()
	capacity := pvc.Status.Capacity["storage"].DeepCopy()

	return &crd.PersistentVolumeClaim{
		Name:                pvc.GetObjectMeta().GetName(),
		VolumeName:          pvc.Spec.VolumeName,
		RequestSizeInGB:     []float64{utils.BytesToGB(request.Value())},
		CapacityAllotedInGB: []float64{utils.BytesToGB(capacity.Value())},
		BoundTimes:          []meta_v1.Time{},
		UnboundTimes:        []meta_v1.Time{},
	}
}

// action to be taken when pod is deleted.
// Unbound all bounded pvcs.
func pvcHandlePodDeletion(podDetails *crd.PodDetails) {
	pvMap := podDetails.PodVolumeClaims
	for claimName := range pvMap {
		if checkBounded(pvMap[claimName]) {
			log.Info("Unbounded pvc: " + claimName + " Reason deletion of pod: " + podDetails.Name)
			pvMap[claimName].UnboundTimes = append(pvMap[claimName].UnboundTimes, podDetails.EndTime)
		}
	}
	podDetails.PodVolumeClaims = pvMap
}

// If length of bound times is 1 more than unbound times it means that the pvc is still bound to pod.
func checkBounded(pvc *crd.PersistentVolumeClaim) bool {
	return len(pvc.BoundTimes)-len(pvc.UnboundTimes) == 1
}

// false if no error
func isPVCError(err error, claimName string) bool {
	if err == nil {
		return false
	}

	if errors.IsNotFound(err) {
		log.Errorf("Persistent Volume Claim %s not found\n", claimName)
		return true
	}
	if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Errorf("Error getting persistence volume Claim %s : %v\n", claimName, statusError.ErrStatus.Message)
		return true
	}
	log.Errorf("Error: Unable to get PVC: %s\n", claimName)
	return true
}
