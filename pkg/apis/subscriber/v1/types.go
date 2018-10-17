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

package v1

import meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// CRD Subscriber attributes
const (
	SubscriberPlural   string = "subscribers"
	SubscriberGroup    string = "kuber.input"
	SubscriberVersion  string = "v1"
	SubscriberFullName string = SubscriberPlural + "." + SubscriberGroup
)

// Subscriber information
type Subscriber struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               SubscriberSpec   `json:"spec"`
	Status             SubscriberStatus `json:"status,omitempty"`
}

// SubscriberSpec definition details
type SubscriberSpec struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster"`
	OrgID       string `json:"orgId"`
	URL         string `json:"url"`
	AuthType    string `json:"authType,omitempty"`
	AuthToken   string `json:"authToken,omitempty"`
}

// SubscriberStatus definition
type SubscriberStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// SubscriberList type
type SubscriberList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []Subscriber `json:"items"`
}
