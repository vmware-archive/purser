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
	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
)

// Dgraph Model Constants
const (
	Islabel = "isLabel"
)

// Label structure for Key:Value
type Label struct {
	dgraph.ID
	IsLabel bool   `json:"isLabel,omitempty"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
}

// GetLabel if label is not in dgraph it creates and returns Label object
func GetLabel(key, value string) *Label {
	xid := getXIDOfLabel(key, value)
	uid := CreateOrGetLabelByID(key, value)
	return &Label{
		ID: dgraph.ID{Xid: xid, UID: uid},
	}
}

// CreateOrGetLabelByID if label is not in dgraph it creates and returns uid of label
func CreateOrGetLabelByID(key, value string) string {
	xid := getXIDOfLabel(key, value)
	uid := dgraph.GetUID(xid, Islabel)
	if uid == "" {
		// create new label and get its uid
		uid = createLabelObject(key, value)
	}
	return uid
}

func getXIDOfLabel(key, value string) string {
	return "label-" + key + "-" + value
}

func createLabelObject(key, value string) string {
	xid := getXIDOfLabel(key, value)
	newLabel := Label{
		ID:      dgraph.ID{Xid: xid},
		IsLabel: true,
		Key:     key,
		Value:   value,
	}
	assigned, err := dgraph.MutateNode(newLabel, dgraph.CREATE)
	if err != nil {
		logrus.Fatal(err)
		return ""
	}
	logrus.Debugf("created label in dgraph: (%v)", newLabel)
	return assigned.Uids["blank-0"]
}
