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

import (
	"github.com/vmware/purser/pkg/apis/subscriber/v1"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// SubscriberInterface has client methods we need to access Subscriber object
type SubscriberInterface interface {
	CreateSubscriber(obj *v1.Subscriber) (*v1.Subscriber, error)
	UpdateSubscriber(obj *v1.Subscriber) (*v1.Subscriber, error)
	DeleteSubscriber(name string, options *meta_v1.DeleteOptions) error
	GetSubscriber(name string) (*v1.Subscriber, error)
	ListSubscriber(opts meta_v1.ListOptions) (*v1.SubscriberList, error)
	NewListWatchSubscriber() *cache.ListWatch
}

// SubscriberClient structure
type SubscriberClient struct {
	client *rest.RESTClient
	ns     string
	plural string
	codec  runtime.ParameterCodec
}

// NewSubscriber creates a new intance of the group CRD client.
func NewSubscriber(client *rest.RESTClient, scheme *runtime.Scheme, namespace string) *SubscriberClient {
	return &SubscriberClient{
		client: client,
		ns:     namespace,
		plural: v1.SubscriberPlural,
		codec:  runtime.NewParameterCodec(scheme),
	}
}

// CreateSubscriber creates a CRD subscriber.
func (c *SubscriberClient) CreateSubscriber(obj *v1.Subscriber) (*v1.Subscriber, error) {
	result := v1.Subscriber{}
	err := c.client.Post().
		Namespace(c.ns).
		Resource(c.plural).
		Body(obj).
		Do().
		Into(&result)
	return &result, err
}

// UpdateSubscriber modifies the subscriber.
func (c *SubscriberClient) UpdateSubscriber(obj *v1.Subscriber) (*v1.Subscriber, error) {
	result := v1.Subscriber{}
	err := c.client.Put().
		Name((obj.Name)).
		Namespace(c.ns).
		Resource(c.plural).
		Body(obj).
		Do().
		Into(&result)
	return &result, err
}

// DeleteSubscriber removes the subscriber.
func (c *SubscriberClient) DeleteSubscriber(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource(c.plural).
		Name(name).
		Body(options).
		Do().
		Error()
}

// GetSubscriber returns the subscriber
func (c *SubscriberClient) GetSubscriber(name string) (*v1.Subscriber, error) {
	result := v1.Subscriber{}
	err := c.client.Get().
		Namespace(c.ns).
		Resource(c.plural).
		Name(name).
		Do().
		Into(&result)
	return &result, err
}

// ListSubscriber fetches the list of subscriber CRD clients.
func (c *SubscriberClient) ListSubscriber(opts meta_v1.ListOptions) (*v1.SubscriberList, error) {
	result := v1.SubscriberList{}
	err := c.client.Get().
		Namespace(c.ns).
		Resource(c.plural).
		VersionedParams(&opts, c.codec).
		Do().
		Into(&result)
	return &result, err
}

// NewListWatchSubscriber creates a new List watch for our TPR
func (c *SubscriberClient) NewListWatchSubscriber() *cache.ListWatch {
	return cache.NewListWatchFromClient(c.client, c.plural, c.ns, fields.Everything())
}
