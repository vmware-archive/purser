package crd

import (
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"kuber-controller/metrics"
	"reflect"
)

const (
	GroupPlural   string = "groups"
	GroupGroup    string = "vmware.kuber"
	GroupVersion  string = "v1"
	GroupFullName string = GroupPlural + "." + GroupGroup
)

// Definition of our CRD Group class
type Group struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec   GroupSpec   `json:"spec"`
	Status GroupStatus `json:"status,omitempty"`
}

type GroupSpec struct {
	Name               string                      `json:"name"`
	CustomGroup        bool                        `json:"custom,omitempty"`
	Type               string                      `json:"type,omitempty"`
	Labels             map[string]string           `json:"labels,omitempty"`
	AllocatedResources *metrics.Metrics            `json:"metrics,omitempty"`
	PodsMetrics        map[string]*metrics.Metrics `json:"pods,omitempty"`
}

type GroupStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

type GroupList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items []Group    `json:"items"`
}

// Create a  Rest client with the new CRD Schema
var GroupSchemeGroupVersion = schema.GroupVersion{Group: GroupGroup, Version: GroupVersion}

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
		out.Items = make([]Group, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
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
