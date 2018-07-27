package crd

import (
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"reflect"
)

const (
	SubscriberPlural   string = "subscribers"
	SubscriberGroup    string = "kuber.input"
	SubscriberVersion  string = "v1"
	SubscriberFullName string = SubscriberPlural + "." + SubscriberGroup
)

type Subscriber struct {
	meta_v1.TypeMeta        `json:",inline"`
	meta_v1.ObjectMeta      `json:"metadata"`
	Spec   SubscriberSpec   `json:"spec"`
	Status SubscriberStatus `json:"status,omitempty"`
}

type SubscriberSpec struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster"`
	Url         string `json:url`
	AuthType    string `json:"authType,omitempty"`
	AuthToken   string `json:"authToken,omitempty"`
}

type SubscriberStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

type SubscriberList struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ListMeta   `json:"metadata"`
	Items []Subscriber `json:"items"`
}

// Create a  Rest client with the new CRD Schema
var SubscriberGroupVersion = schema.GroupVersion{Group: SubscriberGroup, Version: SubscriberVersion}

func CreateSubscriberCRD(clientset apiextcs.Interface) error {
	return CreateCRD(clientset, SubscriberFullName, SubscriberGroup, SubscriberVersion, SubscriberPlural, reflect.TypeOf(Subscriber{}).Name())
}

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *Subscriber) DeepCopyInto(out *Subscriber) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Subscriber) DeepCopyObject() runtime.Object {
	out := Subscriber{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *SubscriberList) DeepCopyObject() runtime.Object {
	out := SubscriberList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Subscriber, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SubscriberGroupVersion,
		&Subscriber{},
		&SubscriberList{},
	)
	meta_v1.AddToGroupVersion(scheme, SubscriberGroupVersion)
	return nil
}

func NewSubscriberClient(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	config := *cfg
	config.GroupVersion = &SubscriberGroupVersion
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
