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
	"reflect"
	"time"

	log "github.com/Sirupsen/logrus"

	subscriber_v1 "github.com/vmware/purser/pkg/apis/subscriber/v1"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

// NewSubscriberClient returns an instance of the Subscriber Client
func NewSubscriberClient(clientset apiextcs.Interface, config *rest.Config) *SubscriberClient {
	err := createSubscriberCRD(clientset)
	if err != nil {
		log.Fatalf("failed to create CRD subscriber %v", err)
	}

	// Wait for the CRD to be created before we use it (only needed if its a new one)
	time.Sleep(3 * time.Second)

	// Create a new clientset which include our CRD schema
	crdcs, scheme, err := newClient(config)
	if err != nil {
		log.Fatalf("failed to add CRD subscriber schema to clientset %v", err)
	}

	// Create a CRD client interface
	return Subscriber(crdcs, scheme, "default")
}

// Subscriber returns an instance of the subscriber client
func Subscriber(client *rest.RESTClient, scheme *runtime.Scheme, namespace string) *SubscriberClient {
	return &SubscriberClient{
		client: client,
		ns:     namespace,
		plural: subscriber_v1.SubscriberPlural,
		codec:  runtime.NewParameterCodec(scheme),
	}
}

func createSubscriberCRD(clientset apiextcs.Interface) error {
	crd := &apiextv1beta1.CustomResourceDefinition{
		ObjectMeta: meta_v1.ObjectMeta{Name: subscriber_v1.SubscriberFullName},
		Spec: apiextv1beta1.CustomResourceDefinitionSpec{
			Group:   subscriber_v1.SubscriberGroup,
			Version: subscriber_v1.SubscriberVersion,
			//TODO: make cluster scoped?
			Scope: apiextv1beta1.NamespaceScoped,
			Names: apiextv1beta1.CustomResourceDefinitionNames{
				Plural: subscriber_v1.SubscriberPlural,
				Kind:   reflect.TypeOf(subscriber_v1.Subscriber{}).Name(),
			},
		},
	}
	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	// Ignore error if it already exists
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func newClient(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	config := *cfg
	scheme, err := setConfigDefaults(&config)
	if err != nil {
		return nil, nil, err
	}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, nil, err
	}
	return client, scheme, nil
}

func setConfigDefaults(config *rest.Config) (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(subscriber_v1.AddKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	config.GroupVersion = &subscriber_v1.SubscriberGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: serializer.NewCodecFactory(scheme)}
	return scheme, nil
}
