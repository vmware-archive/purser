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
	"fmt"
	"net/http"
	"time"

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
	notifiers := getNotifiers(subscribers)

	for _, n := range notifiers {
		go func(n *notifier) {
			req, err := n.createNewRequest(payload)
			if err != nil {
				log.Errorf("Failed to unmarshal payload and create new request %v", err)
			} else {
				err := retry(3, time.Minute, func() error {
					return sendData(req)
				})
				if err != nil {
					log.Errorf("Notification to subscriber %v failed after 3 retries %v", n.url, err)
				}
			}
		}(n)
	}
}

func (n notifier) createNewRequest(payload []*interface{}) (*http.Request, error) {
	payloadWrapper := controller.PayloadWrapper{
		Data:    payload,
		OrgID:   n.orgID,
		Cluster: n.cluster,
	}

	jsonStr, err := json.Marshal(payloadWrapper)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling payload %v", err)
	}

	req, err := http.NewRequest("POST", n.url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	n.setAuthHeaders(req)
	return req, nil
}

func sendData(req *http.Request) error {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending data to %v: %v", req.URL, err)
	}

	if resp != nil {
		if resp.StatusCode != 200 {
			return fmt.Errorf("payload data posting failed for %v, %s", req.URL, resp.Status)
		}
		log.Debugf("Payload data posted successfully for %v", req.URL)
	}
	return nil
}

func (n *notifier) setAuthHeaders(r *http.Request) {
	// TODO: add support for other auth types.
	if n.authType != "" {
		if n.authType == "access-token" {
			r.Header.Set("Authorization", "Bearer "+n.authCode)
		}
	}
}

func getNotifiers(subscribers *subscriber_v1.SubscriberList) []*notifier {
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
	} else {
		log.Debug("No subscribers available.")
	}
	return notifiers
}

func retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}
