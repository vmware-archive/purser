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
	"github.com/vmware/purser/pkg/apis/groups/v1"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// GroupInterface has client methods we need to access Group object
type GroupInterface interface {
	CreateGroup(obj *v1.Group) (*v1.Group, error)
	UpdateGroup(obj *v1.Group) (*v1.Group, error)
	DeleteGroup(name string, options *meta_v1.DeleteOptions) error
	GetGroup(name string) (*v1.Group, error)
	ListGroups(opts meta_v1.ListOptions) (*v1.GroupList, error)
	NewListWatchGroup() *cache.ListWatch
}

// GroupClient defines the CRD Group structure
type GroupClient struct {
	client *rest.RESTClient
	ns     string
	plural string
	codec  runtime.ParameterCodec
}

// NewGroup creates a new intance of the group CRD client.
func NewGroup(client *rest.RESTClient, scheme *runtime.Scheme, namespace string) *GroupClient {
	return &GroupClient{
		client: client,
		ns:     namespace,
		plural: v1.CRDPlural,
		codec:  runtime.NewParameterCodec(scheme),
	}
}

// CreateGroup creates a new group.
func (c *GroupClient) CreateGroup(obj *v1.Group) (*v1.Group, error) {
	result := v1.Group{}
	err := c.client.Post().
		Namespace(c.ns).
		Resource(c.plural).
		Body(obj).
		Do().
		Into(&result)
	return &result, err
}

// UpdateGroup modifies the group specification.
func (c *GroupClient) UpdateGroup(obj *v1.Group) (*v1.Group, error) {
	result := v1.Group{}
	err := c.client.Put().
		Name((obj.Name)).
		Namespace(c.ns).
		Resource(c.plural).
		Body(obj).
		Do().
		Into(&result)
	return &result, err
}

// DeleteGroup removes the group.
func (c *GroupClient) DeleteGroup(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource(c.plural).
		Name(name).
		Body(options).
		Do().
		Error()
}

// GetGroup fetches the group
func (c *GroupClient) GetGroup(name string) (*v1.Group, error) {
	result := v1.Group{}
	err := c.client.Get().
		Namespace(c.ns).
		Resource(c.plural).
		Name(name).
		Do().
		Into(&result)
	return &result, err
}

// ListGroups fetches the list of groups.
func (c *GroupClient) ListGroups(opts meta_v1.ListOptions) (*v1.GroupList, error) {
	result := v1.GroupList{}
	err := c.client.Get().
		Namespace(c.ns).
		Resource(c.plural).
		VersionedParams(&opts, c.codec).
		Do().
		Into(&result)
	return &result, err
}

// NewListWatchGroup creates a new List watch for our TPR
func (c *GroupClient) NewListWatchGroup() *cache.ListWatch {
	return cache.NewListWatchFromClient(c.client, c.plural, c.ns, fields.Everything())
}
