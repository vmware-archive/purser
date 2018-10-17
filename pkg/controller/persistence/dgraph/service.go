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

package dgraph

import (
	"encoding/json"
	"log"
	"time"

	api_v1 "k8s.io/api/core/v1"
)

// Dgraph Model Constants
const (
	IsService = "isService"
)

// Service model structure in Dgraph
type Service struct {
	ID
	IsService bool       `json:"isService,omitempty"`
	Name      string     `json:"name,omitempty"`
	StartTime time.Time  `json:"startTime,omitempty"`
	EndTime   time.Time  `json:"endTime,omitempty"`
	Pod       []*Pod     `json:"servicePods,omitempty"`
	Interacts []*Service `json:"interacts,omitempty"`
}

// PersistService create a new node in the Dgraph or if it exists updates the original node
func PersistService(service api_v1.Service) error {
	xid := service.Namespace + ":" + service.Name
	uid := GetUID(Client, xid, IsService)

	if uid == "" {
		newService := Service{
			Name:      service.Name,
			IsService: true,
			ID:        ID{Xid: xid},
		}
		bytes, err := json.Marshal(newService)
		if err != nil {
			return err
		}
		_, err = MutateNode(Client, bytes)
		return err
	}
	return nil
}

// PersistServicesInteractionGraph stores the service interaction data in the Dgraph
func PersistServicesInteractionGraph(sourceService string, destinationServices []string) error {
	uid := GetUID(Client, sourceService, IsService)
	if uid == "" {
		log.Println("Source Service " + sourceService + " is not persisted yet.")
		return nil
	}

	services := []*Service{}
	for _, destinationService := range destinationServices {
		uid = GetUID(Client, destinationService, IsService)
		if uid == "" {
			continue
		}

		service := &Service{
			ID: ID{UID: uid},
		}
		services = append(services, service)
	}
	source := Service{
		ID:        ID{UID: uid},
		Interacts: services,
	}

	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	_, err = MutateNode(Client, bytes)
	return err

}
