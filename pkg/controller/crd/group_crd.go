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

package crd

import (
	"reflect"

	"github.com/vmware/purser/pkg/controller/metrics"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

const (
	// GroupPlural plural form value
	GroupPlural string = "groups"
	// GroupGroup value
	GroupGroup string = "vmware.kuber"
	// GroupVersion value
	GroupVersion string = "v1"
	// GroupFullName value
	GroupFullName string = GroupPlural + "." + GroupGroup
)

// Group definition of our CRD Group class
type Group struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               GroupSpec   `json:"spec"`
	Status             GroupStatus `json:"status,omitempty"`
}

// GroupSpec specifications for Group
type GroupSpec struct {
	Name               string                      `json:"name"`
	Type               string                      `json:"type,omitempty"`
	Labels             map[string]string           `json:"labels,omitempty"`
	AllocatedResources *metrics.Metrics            `json:"metrics,omitempty"`
	PodsMetrics        map[string]*metrics.Metrics `json:"pods,omitempty"`
	PodsDetails        map[string]*PodDetails      `json:"podDetails,omitempty"`
}

// PodDetails information
type PodDetails struct {
	Name       string
	StartTime  meta_v1.Time
	EndTime    meta_v1.Time
	Containers []*Container
}

// Container information
type Container struct {
	Name    string
	Metrics *metrics.Metrics
}

// GroupStatus information
type GroupStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// GroupList specification
type GroupList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []*Group `json:"items"`
}

// GroupSchemeGroupVersion creates a Rest client with the new CRD Schema
var GroupSchemeGroupVersion = schema.GroupVersion{Group: GroupGroup, Version: GroupVersion}

// CreateGroupCRD returns a group CRD instance.
func CreateGroupCRD(clientset apiextcs.Interface) error {
	return CreateCRD(clientset, GroupFullName, GroupGroup, GroupVersion, GroupPlural, reflect.TypeOf(Group{}).Name())
}

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *Group) DeepCopyInto(out *Group) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Group) DeepCopyObject() runtime.Object {
	out := Group{}
	in.DeepCopyInto(&out)
	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *GroupList) DeepCopyObject() runtime.Object {
	out := GroupList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]*Group, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(out.Items[i])
		}
	}
	return &out
}

func addGroupKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupSchemeGroupVersion,
		&Group{},
		&GroupList{},
		&Subscriber{},
		&SubscriberList{},
	)
	meta_v1.AddToGroupVersion(scheme, GroupSchemeGroupVersion)
	return nil
}

// NewGroupClient returns an instance of group REST client.
func NewGroupClient(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addGroupKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	config := *cfg
	config.GroupVersion = &GroupSchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, nil, err
	}
	return client, scheme, nil
}
