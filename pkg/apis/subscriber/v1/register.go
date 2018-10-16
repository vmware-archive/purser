package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SchemeBuilder parameters
var (
	SchemeBuilder = runtime.NewSchemeBuilder(AddKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// GroupName is the group name use in this package
const GroupName = "kuber.input"

// SubscriberGroupVersion is group version used to register these objects
var SubscriberGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SubscriberGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SubscriberGroupVersion.WithResource(resource).GroupResource()
}

// AddKnownTypes ...
func AddKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SubscriberGroupVersion,
		&Subscriber{},
		&SubscriberList{},
	)
	meta_v1.AddToGroupVersion(scheme, SubscriberGroupVersion)

	return nil
}
