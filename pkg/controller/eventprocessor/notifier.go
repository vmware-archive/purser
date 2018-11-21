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
	"net/http"

	log "github.com/Sirupsen/logrus"
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
	if notifiers := getNotifiers(subscribers); notifiers != nil {
		for _, notifier := range notifiers {
			notifier.sendData(payload)
		}
	}
}

func (n notifier) sendData(payload []*interface{}) {
	payloadWrapper := controller.PayloadWrapper{Data: payload, OrgID: n.orgID, Cluster: n.cluster}
	jsonStr, err := json.Marshal(payloadWrapper)
	if err != nil {
		log.Error("Error while unmarshalling payload ", err)
	}

	req, err := http.NewRequest("POST", n.url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Error("Error while creating the http request ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	n.setAuthHeaders(req)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error while sending data to subscriber %v"+n.url, err)
	} else if resp != nil {
		if resp.StatusCode == 200 {
			log.Info("Data is posted successfully for subscriber " + n.url)
		} else {
			log.Info("Data posting failed for subscriber " + n.url + " " + resp.Status)
		}
	}
}

func (n *notifier) setAuthHeaders(r *http.Request) {
	//TODO: add support for other auth types.
	if n.authType != "" {
		if n.authType == "access-token" {
			r.Header.Set("Authorization", "Bearer "+n.authCode)
		}
	}
}

func getNotifiers(subscribers *subscriber_v1.SubscriberList) []*notifier {
	notifiers := []*notifier{}
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
	} else {
		log.Debug("There are no notifiers for subscribers")
	}
	return notifiers
}
