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
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	log "github.com/Sirupsen/logrus"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/vmware/purser/pkg/controller/dgraph"
)

// Dgraph Model Constants
const (
	IsProc = "isProc"
)

// Proc schema in dgraph
type Proc struct {
	dgraph.ID
	IsProc    bool       `json:"isProc,omitemtpy"`
	Name      string     `json:"name,omitempty"`
	Interacts []*Pod     `json:"interacts,omitempty"`
	Container Container  `json:"container,omitempty"`
	StartTime string  `json:"startTime,omitempty"`
	EndTime   string  `json:"endTime,omitempty"`
	Namespace *Namespace `json:"namespace,omitempty"`
	Type      string     `json:"type,omitempty"`
}

func newProc(procXID, procName, containerUID, containerXID string, creationTimeStamp time.Time) (*api.Assigned, error) {
	newProc := Proc{
		ID:        dgraph.ID{Xid: procXID},
		IsProc:    true,
		Type:      "process",
		Name:      procName,
		Container: Container{ID: dgraph.ID{UID: containerUID, Xid: containerXID}},
		StartTime: creationTimeStamp.Format(time.RFC3339),
	}
	return dgraph.MutateNode(newProc, dgraph.CREATE)
}

// StoreProcess ...
func StoreProcess(procXID, containerXID string, podsXIDs []string, creationTimeStamp time.Time) error {
	// fetch the 4th field from ns : podName : containerName : procID : procName
	procName := strings.Join(strings.Split(procXID, ":")[4:], "-")
	containerUID := dgraph.GetUID(containerXID, IsContainer)
	if containerUID == "" {
		return fmt.Errorf("Container not persisted yet")
	}

	procUID := dgraph.GetUID(procXID, IsProc)
	if procUID == "" {
		assigned, err := newProc(procXID, procName, containerUID, containerXID, creationTimeStamp)
		if err != nil {
			logrus.Errorf("Unable to create proc: %s", procXID)
			return err
		}
		log.Infof("Process with xid: (%s) persisted in dgraph", procXID)
		procUID = assigned.Uids["blank-0"]
	}

	pods := retrievePodsFromPodsXIDs(podsXIDs)
	updatedProc := Proc{
		ID:        dgraph.ID{UID: procUID, Xid: procXID},
		Interacts: pods,
	}
	_, err := dgraph.MutateNode(updatedProc, dgraph.UPDATE)
	return err
}

func deleteProcessesInTerminatedContainers(containers []*Container) {
	procs := []Proc{}
	for _, container := range containers {
		for _, proc := range container.Procs {
			updatedProc := Proc{
				ID:      dgraph.ID{UID: proc.ID.UID, Xid: proc.ID.Xid},
				EndTime: container.EndTime,
			}
			procs = append(procs, updatedProc)
		}
	}
	_, err := dgraph.MutateNode(procs, dgraph.UPDATE)
	if err != nil {
		log.Error(err)
	}
}

func retrieveProcessesFromProcessesXIDs(procsXIDs []string) []*Proc {
	procs := []*Proc{}
	for _, procXID := range procsXIDs {
		procUID := dgraph.GetUID(procXID, IsProc)
		if procUID != "" {
			procs = append(procs, &Proc{ID: dgraph.ID{UID: procUID, Xid: procXID}})
		}
	}
	return procs
}

// RetrieveAllProcess ...
func RetrieveAllProcess() ([]byte, error) {
	const q = `query {
		process(func: has(isProc)) {
			name
			type
		}
	}`

	result, err := dgraph.ExecuteQueryRaw(q)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RetrieveProcess ...
func RetrieveProcess(name string) ([]byte, error) {
	q := `query {
		process(func: has(isProc)) @filter(eq(name, "` + name + `")) {
			name
			type
		}
	}`

	result, err := dgraph.ExecuteQueryRaw(q)
	if err != nil {
		return nil, err
	}
	return result, nil
}
