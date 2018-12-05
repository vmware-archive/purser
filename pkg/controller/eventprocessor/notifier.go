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

package eventprocessor

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/matryer/try.v1"

	subscriber_v1 "github.com/vmware/purser/pkg/apis/subscriber/v1"
	"github.com/vmware/purser/pkg/controller"
)

// ReadSize defines the default payload read size
const ReadSize uint32 = 50

type notifier struct {
	url      string
	authType string
	authCode string
	cluster  string
	orgID    string
}

func notifySubscribers(payload []*interface{}, subscribers *subscriber_v1.SubscriberList) {
	var notifiers []*notifier
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		notifiers, err = getNotifiers(subscribers)
		if err != nil {
			time.Sleep(1 * time.Minute) // wait a minute
		}
		return attempt < 3, err
	})
	if err != nil {
		log.Debugf("Retry unsuccessful. %v", err)
	}

	for _, notifier := range notifiers {
		notifier.sendData(payload)
	}
}

func (n notifier) sendData(payload []*interface{}) {
	payloadWrapper := controller.PayloadWrapper{
		Data:    payload,
		OrgID:   n.orgID,
		Cluster: n.cluster,
	}

	jsonStr, err := json.Marshal(payloadWrapper)
	if err != nil {
		log.Errorf("Error unmarshalling payload %v", err)
	}

	req, err := http.NewRequest("POST", n.url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Errorf("Error creating HTTP request %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	n.setAuthHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending data to subscriber %s: %v", n.url, err)
	}

	if resp != nil {
		if resp.StatusCode == 200 {
			log.Infof("Payload data posted successfully for subscriber %s", n.url)
		} else {
			log.Infof("Payload data posting failed for subscriber %s, %s ", n.url, resp.Status)
		}
	}
}

func (n *notifier) setAuthHeaders(r *http.Request) {
	// TODO: add support for other auth types.
	if n.authType != "" {
		if n.authType == "access-token" {
			r.Header.Set("Authorization", "Bearer "+n.authCode)
		}
	}
}

func getNotifiers(subscribers *subscriber_v1.SubscriberList) ([]*notifier, error) {
	var notifiers []*notifier
	if len(subscribers.Items) > 0 {
		for _, sub := range subscribers.Items {
			notifier := &notifier{
				url:      sub.Spec.URL,
				authType: sub.Spec.AuthType,
				authCode: sub.Spec.AuthToken,
				cluster:  sub.Spec.ClusterName,
				orgID:    sub.Spec.OrgID,
			}
			notifiers = append(notifiers, notifier)
		}
		return notifiers, nil
	}
	return notifiers, errors.New("no notifiers available for subscribers")
}
