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
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/vmware/purser/pkg/controller/utils"
)

// Dgraph Model Constants
const (
	IsProc = "isProc"
)

// Proc schema in dgraph
type Proc struct {
	ID
	IsProc    bool      `json:"isProc,omitemtpy"`
	Name      string    `json:"name,omitempty"`
	Interacts []*Pod    `json:"interacts,omitempty"`
	Container Container `json:"container,omitempty"`
}

func createProc(procXID, procName, containerUID string) (*api.Assigned, error) {
	newProc := Proc{
		ID:        ID{Xid: procXID},
		IsProc:    true,
		Name:      procName,
		Container: Container{ID: ID{UID: containerUID}},
	}
	bytes := utils.JsonMarshal(newProc)
	return MutateNode(Client, bytes)
}

// PersistProc ...
func PersistProc(procXID, procName string, podsXIDs []string, containerXID string) error {
	containerUID := GetUID(Client, containerXID, IsContainer)
	if containerUID == "" {
		return fmt.Errorf("Container not persisted yet")
	}

	procUID := GetUID(Client, procXID, IsProc)
	if procUID == "" {
		assigned, err := createProc(procXID, procName, containerUID)
		if err != nil {
			logrus.Errorf("Unable to create proc: %s", procXID)
			return err
		}
		procUID = assigned.Uids["blank-0"]
	}

	pods := []*Pod{}
	for _, podXID := range podsXIDs {
		podUID := GetUID(Client, podXID, IsPod)
		if podUID != "" {
			pods = append(pods, &Pod{ID: ID{UID: podUID}})
		}
	}

	updatedProc := Proc{
		ID:        ID{UID: procUID},
		Interacts: pods,
	}
	bytes := utils.JsonMarshal(updatedProc)
	_, err := MutateNode(Client, bytes)
	return err
}
