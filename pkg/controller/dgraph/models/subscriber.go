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
	"time"

	subscribers_v1 "github.com/vmware/purser/pkg/apis/subscriber/v1"
	"github.com/vmware/purser/pkg/controller/dgraph"
)

// Dgraph Model Constants
const (
	IsSubscriber = "isSubscriber"
)

// SubscriberCRD schema in dgraph
type SubscriberCRD struct {
	dgraph.ID
	IsSubscriber bool           `json:"isSubscriber,omitempty"`
	Name         string         `json:"name,omitempty"`
	StartTime    string         `json:"startTime,omitempty"`
	EndTime      string         `json:"endTime,omitempty"`
	Type         string         `json:"type,omitempty"`
	Spec         SubscriberSpec `json:"spec"`
}

// SubscriberSpec definition details
type SubscriberSpec struct {
	Name    string            `json:"name"`
	Headers map[string]string `json:"headers"`
	URL     string            `json:"url"`
}

func createSubscriberCRDObject(subscriber subscribers_v1.Subscriber) SubscriberCRD {
	newSubscriber := SubscriberCRD{
		Name:         subscriber.Name,
		IsSubscriber: true,
		Type:         subscribers_v1.SubscriberGroup,
		ID:           dgraph.ID{Xid: "subscriber-" + subscriber.Name},
		StartTime:    subscriber.GetCreationTimestamp().Time.Format(time.RFC3339),
		Spec: SubscriberSpec{
			Name: subscriber.Spec.Name,
			Headers: subscriber.Spec.Headers,
			URL: subscriber.Spec.URL,
		},
	}

	deletionTimestamp := subscriber.GetDeletionTimestamp()
	if !deletionTimestamp.IsZero() {
		newSubscriber.EndTime = deletionTimestamp.Time.Format(time.RFC3339)
	}
	return newSubscriber
}

// StoreSubscriberCRD create a new subscriber CRD in the Dgraph and updates if already present.
func StoreSubscriberCRD(subscriber subscribers_v1.Subscriber) (string, error) {
	xid := "subscriber-" + subscriber.Name
	uid := dgraph.GetUID(xid, IsSubscriber)

	if uid != "" {
		return uid, nil
	}

	newSubscriber := createSubscriberCRDObject(subscriber)
	assigned, err := dgraph.MutateNode(newSubscriber, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	logrus.Infof("Subscriber: (%v) persisted in dgraph", subscriber.Name)
	return assigned.Uids["blank-0"], nil
}
