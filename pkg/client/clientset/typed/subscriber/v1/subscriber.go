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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

// SubscriberInterface has client methods we need to access Subscriber object
type SubscriberInterface interface {
	Create(obj *v1.Subscriber) (*v1.Subscriber, error)
	Update(obj *v1.Subscriber) (*v1.Subscriber, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	Get(name string) (*v1.Subscriber, error)
	List(opts meta_v1.ListOptions) (*v1.SubscriberList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
}

// SubscriberClient structure
type SubscriberClient struct {
	client *rest.RESTClient
	ns     string
	plural string
	codec  runtime.ParameterCodec
}

// Create creates a CRD subscriber.
func (c *SubscriberClient) Create(obj *v1.Subscriber) (*v1.Subscriber, error) {
	result := v1.Subscriber{}
	err := c.client.Post().
		Namespace(c.ns).
		Resource(c.plural).
		Body(obj).
		Do().
		Into(&result)
	return &result, err
}

// Update modifies the subscriber.
func (c *SubscriberClient) Update(obj *v1.Subscriber) (*v1.Subscriber, error) {
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

// Delete removes the subscriber.
func (c *SubscriberClient) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource(c.plural).
		Name(name).
		Body(options).
		Do().
		Error()
}

// Get returns the subscriber
func (c *SubscriberClient) Get(name string) (*v1.Subscriber, error) {
	result := v1.Subscriber{}
	err := c.client.Get().
		Namespace(c.ns).
		Resource(c.plural).
		Name(name).
		Do().
		Into(&result)
	return &result, err
}

// List fetches the list of subscriber CRD clients.
func (c *SubscriberClient) List(opts meta_v1.ListOptions) (*v1.SubscriberList, error) {
	result := v1.SubscriberList{}
	err := c.client.Get().
		Namespace(c.ns).
		Resource(c.plural).
		VersionedParams(&opts, c.codec).
		Do().
		Into(&result)
	return &result, err
}

// Watch watches for the subcriber CRD
func (c *SubscriberClient) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.
		Get().
		Namespace(c.ns).
		Resource(c.plural).
		VersionedParams(&opts, c.codec).
		Watch()
}
