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
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/vmware/purser/pkg/controller/dgraph"
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
	Pod       []*Pod     `json:"pod,omitempty"`
	Interacts []*Service `json:"interacts,omitempty"`
	Namespace *Namespace `json:"namespace,omitempty"`
}

func newService(svc api_v1.Service) (*api.Assigned, error) {
	newService := Service{
		Name:      svc.Name,
		IsService: true,
		ID:        dgraph.ID{Xid: svc.Namespace + ":" + svc.Name},
		StartTime: svc.GetCreationTimestamp().Time,
	}
	namespaceUID, err := createOrGetNamespaceByID(svc.Namespace)
	if err == nil {
		newService.Namespace = &Namespace{ID: dgraph.ID{UID: namespaceUID, Xid: svc.Namespace}}
	}
	return dgraph.MutateNode(newService, dgraph.CREATE)
}

// StoreService create a new node in the Dgraph  if it is not present.
func StoreService(service api_v1.Service) error {
	xid := service.Namespace + ":" + service.Name
	uid := dgraph.GetUID(xid, IsService)

	if uid == "" {
		assigned, err := newService(service)
		if err != nil {
			return err
		}
		log.Infof("Service with xid: (%s) persisted in dgraph", xid)
		uid = assigned.Uids["blank-0"]
	}

	svcDeletionTimestamp := service.GetDeletionTimestamp()
	if !svcDeletionTimestamp.IsZero() {
		updatedService := Service{
			ID:      dgraph.ID{Xid: xid, UID: uid},
			EndTime: svcDeletionTimestamp.Time,
		}
		_, err := dgraph.MutateNode(updatedService, dgraph.UPDATE)
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

	services := retrieveServicesFromServicesXIDs(destinationServicesXIDs)
	source := Service{
		ID:        dgraph.ID{UID: uid, Xid: sourceServiceXID},
		Interacts: services,
	}
	_, err := dgraph.MutateNode(source, dgraph.UPDATE)
	return err
}

// StorePodServiceEdges saves pods in Services object in the dgraph
func StorePodServiceEdges(svcXID string, podsXIDsInService []string) error {
	svcUID := dgraph.GetUID(svcXID, IsService)
	if svcUID != "" {
		svcPods := retrievePodsFromPodsXIDs(podsXIDsInService)
		updatedService := Service{
			ID:  dgraph.ID{UID: svcUID, Xid: svcXID},
			Pod: svcPods,
		}
		_, err := dgraph.MutateNode(updatedService, dgraph.UPDATE)
		return err
	}
	return fmt.Errorf("Service with xid: (%s) not in dgraph", svcXID)
}

// RetrieveAllServices returns all pods in the dgraph
func RetrieveAllServices() ([]Service, error) {
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

// RetrieveAllServicesWithDstPods returns all pods in the dgraph
func RetrieveAllServicesWithDstPods() ([]Service, error) {
	const q = `query {
		services(func: has(isService)) {
			xid
			name
			pod {
				name
				interacts @facets {
					name
				}
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

// RetrieveServiceList ...
func RetrieveServiceList() ([]Service, error) {
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

func retrieveServicesFromServicesXIDs(svcsXIDs []string) []*Service {
	services := []*Service{}
	for _, svcXID := range svcsXIDs {
		svcUID := dgraph.GetUID(svcXID, IsService)
		if svcUID == "" {
			continue
		}

		service := &Service{
			ID: dgraph.ID{UID: svcUID, Xid: svcXID},
		}
		services = append(services, service)
	}
	return services
}
