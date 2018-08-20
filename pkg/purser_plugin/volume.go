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

package purser_plugin

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

// PersistentVolume details
type PersistentVolume struct {
	name         string
	capacityInGB float64
	storageClass string
}

// PersistentVolumeClaim details
type PersistentVolumeClaim struct {
	name                string
	volumeName          string
	requestSizeInGB     float64
	capacityAllotedInGB float64
	storageClass        *string
	pv                  *PersistentVolume
}

func collectPersistentVolumeClaims(pvcs map[string]*PersistentVolumeClaim) map[string]*PersistentVolumeClaim {
	for key := range pvcs {
		pvc := collectPersistentVolumeClaim(key)
		pvcs[key] = pvc
	}
	return pvcs
}

func collectPersistentVolumeClaim(claimName string) *PersistentVolumeClaim {
	pvc, err := ClientSetInstance.CoreV1().PersistentVolumeClaims("default").Get(claimName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Persistent Volume Claim %s not found\n", claimName)
		return nil
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting persistence volume Claim %s : %v\n", claimName, statusError.ErrStatus.Message)
		return nil
	} else if err != nil {
		panic(err.Error())
	} else {
		var claim PersistentVolumeClaim
		claim.name = pvc.GetObjectMeta().GetName()
		claim.volumeName = pvc.Spec.VolumeName
		claim.storageClass = pvc.Spec.StorageClassName
		request := pvc.Spec.Resources.Requests["storage"].DeepCopy()
		claim.requestSizeInGB = (float64)(request.Value()) / (float64)(1024.0*1024.0*1024.0)
		capacity := pvc.Status.Capacity["storage"].DeepCopy()
		claim.capacityAllotedInGB = (float64)(capacity.Value()) / (float64)(1024.0*1024.0*1024.0)
		return &claim
	}
}

func collectPersistentVolume(volName string) *PersistentVolume {
	pv, err := ClientSetInstance.CoreV1().PersistentVolumes().Get(volName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Persistent Volume %s not found\n", volName)
		return nil
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting persistence volume %s : %v\n", volName, statusError.ErrStatus.Message)
		return nil
	} else if err != nil {
		panic(err.Error())
	} else {
		var persistentVolume PersistentVolume
		persistentVolume.storageClass = pv.Spec.StorageClassName
		persistentVolume.name = pv.Name
		q := pv.Spec.Capacity["storage"].DeepCopy()
		persistentVolume.capacityInGB = (float64)(q.Value()) / (float64)(1024.0*1024.0*1024.0)
		return &persistentVolume
	}
}

func GetClusterVolumes() []v1.PersistentVolume {
	pvs, err := ClientSetInstance.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pvs.Items
}

func GetClusterPersistentVolumeClaims() []v1.PersistentVolumeClaim {
	pvcs, err := ClientSetInstance.CoreV1().PersistentVolumeClaims("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pvcs.Items
}
