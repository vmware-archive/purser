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
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/utils"
	api_v1 "k8s.io/api/core/v1"
)

// Dgraph Model Constants
const (
	IsService = "isService"
)

// Service model structure in Dgraph
type Service struct {
	dgraph.ID
	IsService bool       `json:"isService,omitempty"`
	Name      string     `json:"name,omitempty"`
	StartTime time.Time  `json:"startTime,omitempty"`
	EndTime   time.Time  `json:"endTime,omitempty"`
	Pod       []*Pod     `json:"servicePods,omitempty"`
	Interacts []*Service `json:"interacts,omitempty"`
}

// StoreService create a new node in the Dgraph  if it is not present.
func StoreService(service api_v1.Service) error {
	xid := service.Namespace + ":" + service.Name
	uid := dgraph.GetUID(xid, IsService)

	if uid == "" {
		newService := Service{
			Name:      service.Name,
			IsService: true,
			ID:        dgraph.ID{Xid: xid},
		}
		bytes := utils.JSONMarshal(newService)
		_, err := dgraph.MutateNode(bytes)
		return err
	}
	return nil
}

// StoreServicesInteraction stores the service interaction data in the Dgraph
func StoreServicesInteraction(sourceServiceXID string, destinationServicesXIDs []string) error {
	uid := dgraph.GetUID(sourceServiceXID, IsService)
	if uid == "" {
		log.Println("Source Service " + sourceServiceXID + " is not persisted yet.")
		return fmt.Errorf("source service: %s is not persisted yet", sourceServiceXID)
	}

	services := []*Service{}
	for _, destinationServiceXID := range destinationServicesXIDs {
		uid = dgraph.GetUID(destinationServiceXID, IsService)
		if uid == "" {
			continue
		}

		service := &Service{
			ID: dgraph.ID{UID: uid},
		}
		services = append(services, service)
	}
	source := Service{
		ID:        dgraph.ID{UID: uid},
		Interacts: services,
	}

	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	_, err = dgraph.MutateNode(bytes)
	return err
}

// StorePodServiceEdges saves pods in Services object in the dgraph
func StorePodServiceEdges(svcToPod map[string][]string) {
	for svcXID, podXIDs := range svcToPod {
		svcUID := dgraph.GetUID(svcXID, IsService)
		if svcUID == "" {
			continue
		}

		svcPods := []*Pod{}
		for _, podXID := range podXIDs {
			podUID := dgraph.GetUID(podXID, IsPod)
			if podUID == "" {
				log.Debugf("Pod uid is empty for pod xid: %s", podXID)
				continue
			}
			pod := &Pod{
				ID: dgraph.ID{UID: podUID},
			}
			svcPods = append(svcPods, pod)
		}

		updatedService := Service{
			ID:  dgraph.ID{UID: svcUID},
			Pod: svcPods,
		}
		bytes := utils.JSONMarshal(updatedService)
		_, err := dgraph.MutateNode(bytes)
		if err != nil {
			log.Error(err)
		}
	}
}

// RetreiveAllServices returns all pods in the dgraph
func RetreiveAllServices() ([]Service, error) {
	const q = `query {
		services(func: has(isService)) {
			name
			interacts @facets {
				name
			}
			pod {
				name
			}
		}
	}`

	type root struct {
		Services []Service `json:"services"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}

	return newRoot.Services, nil
}

// RetreiveServiceList ...
func RetreiveServiceList() ([]Service, error) {
	const q = `query {
		serviceList(func: has(isService)) {
			name
		}
	}`

	type root struct {
		ServiceList []Service `json:"serviceList"`
	}
	newRoot := root{}
	err := dgraph.ExecuteQuery(q, &newRoot)
	if err != nil {
		return nil, err
	}

	return newRoot.ServiceList, nil
}
