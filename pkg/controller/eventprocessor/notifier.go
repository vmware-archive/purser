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
	"github.com/vmware/purser/pkg/controller"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReadSize defines the default payload read size
const ReadSize uint32 = 50

type subscriber struct {
	url      string
	authType string
	authCode string
	cluster  string
	orgID    string
}

// NotifySubscribers notifies subscribers of the process event.
func NotifySubscribers(payload []*interface{}, subscribers []*subscriber) {
	if subscribers == nil {
		return
	}

	for _, subscriber := range subscribers {
		subscriber.sendData(payload)
	}
}

func (subscriber *subscriber) sendData(payload []*interface{}) {
	payloadWrapper := controller.PayloadWrapper{Data: payload, OrgID: subscriber.orgID, Cluster: subscriber.cluster}
	jsonStr, err := json.Marshal(payloadWrapper)
	if err != nil {
		log.Error("Error while unmarshalling payload ", err)
	}

	req, err := http.NewRequest("POST", subscriber.url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Error("Error while creating the http request ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	subscriber.setAuthHeaders(req)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error while sending data to subscriber "+subscriber.url, err)
	} else if resp != nil {
		if resp.StatusCode == 200 {
			log.Info("Data is posted successfully for subscriber " + subscriber.url)
		} else {
			log.Info("Data posting failed for subscriber " + subscriber.url + " " + resp.Status)
		}
	}
}

func (subscriber *subscriber) setAuthHeaders(r *http.Request) {
	//TODO: add support for other auth types.
	if subscriber.authType != "" {
		if subscriber.authType == "access-token" {
			r.Header.Set("Authorization", "Bearer "+subscriber.authCode)
		}
	}
}

func getSubscribers(conf *controller.Config) []*subscriber {
	subscribers := []*subscriber{}
	list, err := conf.Subscriberclient.ListSubscriber(meta_v1.ListOptions{})
	if err != nil {
		log.Error("Error while fetching subscribers list ", err)
		return nil
	}

	if list != nil && len(list.Items) > 0 {
		for _, sub := range list.Items {
			subscriber := &subscriber{
				url:      sub.Spec.URL,
				authType: sub.Spec.AuthType,
				authCode: sub.Spec.AuthToken,
				cluster:  sub.Spec.ClusterName,
				orgID:    sub.Spec.OrgID,
			}
			subscribers = append(subscribers, subscriber)
		}
	} else {
		log.Debug("There are no subscribers")
		return nil
	}
	return subscribers
}
