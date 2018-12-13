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

// Services model structure in Dgraph
type Service struct {
	dgraph.ID
	IsService bool       `json:"isService,omitempty"`
	Name      string     `json:"name,omitempty"`
	StartTime string     `json:"startTime,omitempty"`
	EndTime   string     `json:"endTime,omitempty"`
	Pod       []*Pod     `json:"pod,omitempty"`
	Services  []*Service `json:"service,omitempty"`
	Namespace *Namespace `json:"namespace,omitempty"`
	Type      string     `json:"type,omitempty"`
}

func newService(svc api_v1.Service) (*api.Assigned, error) {
	newService := Service{
		Name:      "service-" + svc.Name,
		IsService: true,
		Type:      "service",
		ID:        dgraph.ID{Xid: svc.Namespace + ":" + svc.Name},
		StartTime: svc.GetCreationTimestamp().Time.Format(time.RFC3339),
	}
	namespaceUID := CreateOrGetNamespaceByID(svc.Namespace)
	if namespaceUID != "" {
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
		log.Infof("Services with xid: (%s) persisted in dgraph", xid)
		uid = assigned.Uids["blank-0"]
	}

	svcDeletionTimestamp := service.GetDeletionTimestamp()
	if !svcDeletionTimestamp.IsZero() {
		updatedService := Service{
			ID:      dgraph.ID{Xid: xid, UID: uid},
			EndTime: svcDeletionTimestamp.Time.Format(time.RFC3339),
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
		log.Debugf("Source Service " + sourceServiceXID + " is not persisted yet.")
		return fmt.Errorf("source service: (%s) is not persisted yet", sourceServiceXID)
	}

	services := retrieveServicesFromServicesXIDs(destinationServicesXIDs)
	log.Debugf("source service: %s, dstServicesCount: %d", sourceServiceXID, len(services))
	source := Service{
		ID:       dgraph.ID{UID: uid, Xid: sourceServiceXID},
		Services: services,
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
	return fmt.Errorf("Services with xid: (%s) not in dgraph", svcXID)
}

func retrieveServicesFromServicesXIDs(svcsXIDs []string) []*Service {
	services := []*Service{}
	for _, svcXID := range svcsXIDs {
		svcUID := dgraph.GetUID(svcXID, IsService)
		if svcUID == "" {
			log.Debugf("dst svc with xid: (%s) not in dgraph", svcXID)
			continue
		}

		service := &Service{
			ID: dgraph.ID{UID: svcUID, Xid: svcXID},
		}
		services = append(services, service)
	}
	return services
}
