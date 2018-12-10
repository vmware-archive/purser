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
	"github.com/vmware/purser/pkg/controller/dgraph/models"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	subscriber_v1 "github.com/vmware/purser/pkg/apis/subscriber/v1"
	"github.com/vmware/purser/pkg/controller"
)

// ReadSize defines the default payload read size
const ReadSize uint32 = 50

type notifier struct {
	url     string
	headers map[string]string
}

func notifySubscribers(payload []*interface{}, subscribers []models.SubscriberCRD) {
	notifiers := getNotifiers(subscribers)

	for _, n := range notifiers {
		req, err := n.createNewRequest(payload)
		if err != nil {
			log.Errorf("Failed to unmarshal payload and create new request %v", err)
		} else {
			err := retry(3, time.Second, func() error {
				return sendData(req)
			})
			if err != nil {
				log.Errorf("Notification to subscriber %v failed after 3 retries %v", n.url, err)
			}
		}
	}
}

func (n notifier) createNewRequest(payload []*interface{}) (*http.Request, error) {
	payloadWrapper := controller.PayloadWrapper{Data: payload}

	jsonStr, err := json.Marshal(payloadWrapper)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling payload %v", err)
	}

	req, err := http.NewRequest("POST", n.url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request %v", err)
	}
	n.setReqHeaders(req)
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

func (n *notifier) setReqHeaders(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	for key, value := range n.headers {
		r.Header.Set(key, value)
	}
}

func getNotifiers(subscribers []models.SubscriberCRD) []*notifier {
	var notifiers []*notifier
	if len(subscribers) > 0 {
		for _, sub := range subscribers {
			notifier := &notifier{
				url:     sub.Spec.URL,
				headers: sub.Spec.Headers,
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
