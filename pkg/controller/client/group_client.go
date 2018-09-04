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

package client

import (
	"github.com/vmware/purser/pkg/controller/crd"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// GroupCrdClient definition
type GroupCrdClient struct {
	cl     *rest.RESTClient
	ns     string
	plural string
	codec  runtime.ParameterCodec
}

// This file implement all the (CRUD) client methods we need to access Group CRD object

// CreateGroupCrdClient creates a new intance of the group CRD client.
func CreateGroupCrdClient(cl *rest.RESTClient, scheme *runtime.Scheme, namespace string) *GroupCrdClient {
	return &GroupCrdClient{cl: cl, ns: namespace, plural: crd.GroupPlural,
		codec: runtime.NewParameterCodec(scheme)}
}

// CreateGroup creates a new group.
func (f *GroupCrdClient) CreateGroup(obj *crd.Group) (*crd.Group, error) {
	var result crd.Group
	err := f.cl.Post().
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

// UpdateGroup modifies the group specification.
func (f *GroupCrdClient) UpdateGroup(obj *crd.Group) (*crd.Group, error) {
	var result crd.Group
	err := f.cl.
		//Put().
		Put().Name((obj.Name)).
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

// DeleteGroup removes the group.
func (f *GroupCrdClient) DeleteGroup(name string, options *meta_v1.DeleteOptions) error {
	return f.cl.Delete().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Body(options).Do().
		Error()
}

// GetGroup fetches the group
func (f *GroupCrdClient) GetGroup(name string) (*crd.Group, error) {
	var result crd.Group
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Do().Into(&result)
	return &result, err
}

// ListGroups fetches the list of groups.
func (f *GroupCrdClient) ListGroups(opts meta_v1.ListOptions) (*crd.GroupList, error) {
	var result crd.GroupList
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		VersionedParams(&opts, f.codec).
		Do().Into(&result)
	return &result, err
}

// NewListWatchGroup creates a new List watch for our TPR
func (f *GroupCrdClient) NewListWatchGroup() *cache.ListWatch {
	return cache.NewListWatchFromClient(f.cl, f.plural, f.ns, fields.Everything())
}
