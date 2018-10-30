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

package models

import (
	"time"

	"github.com/vmware/purser/pkg/controller/dgraph"
	apps_v1beta1 "k8s.io/api/apps/v1beta1"
)

// Dgraph Model Constants
const (
	IsDeployment = "isDeployment"
)

// Deployment schema in dgraph
type Deployment struct {
	dgraph.ID
	IsDeployment bool      `json:"isDeployment,omitempty"`
	Name         string    `json:"name,omitempty"`
	StartTime    time.Time `json:"startTime,omitempty"`
	EndTime      time.Time `json:"endTime,omitempty"`
	Pods         []*Pod    `json:"pods,omitempty"`
}

func createDeploymentObject(deployment apps_v1beta1.Deployment) Deployment {
	newDeployment := Deployment{
		Name:         deployment.Name,
		IsDeployment: true,
		ID:           dgraph.ID{Xid: deployment.Namespace + ":" + deployment.Name},
		StartTime:    deployment.GetCreationTimestamp().Time,
	}
	deploymentDeletionTimestamp := deployment.GetDeletionTimestamp()
	if !deploymentDeletionTimestamp.IsZero() {
		newDeployment.EndTime = deploymentDeletionTimestamp.Time
	}
	return newDeployment
}

// StoreDeployment create a new deployment in the Dgraph and updates if already present.
func StoreDeployment(deployment apps_v1beta1.Deployment) (string, error) {
	xid := deployment.Namespace + ":" + deployment.Name
	uid := dgraph.GetUID(xid, IsDeployment)

	newDeployment := createDeploymentObject(deployment)
	if uid != "" {
		newDeployment.UID = uid
	}
	assigned, err := dgraph.MutateNode(newDeployment, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	return assigned.Uids["blank-0"], nil
}
