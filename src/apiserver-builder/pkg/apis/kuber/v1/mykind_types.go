
/*
 * licensed to vmware.
*/


package v1

import (
	"log"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"apiserver-builder/pkg/apis/kuber"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MyKind
// +k8s:openapi-gen=true
// +resource:path=mykinds,strategy=MyKindStrategy
type MyKind struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyKindSpec   `json:"spec,omitempty"`
	Status MyKindStatus `json:"status,omitempty"`
}

// MyKindSpec defines the desired state of MyKind
type MyKindSpec struct {
}

// MyKindStatus defines the observed state of MyKind
type MyKindStatus struct {
}

// Validate checks that an instance of MyKind is well formed
func (MyKindStrategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*kuber.MyKind)
	log.Printf("Validating fields for MyKind %s\n", o.Name)
	errors := field.ErrorList{}
	// perform validation here and add to errors using field.Invalid
	return errors
}

// DefaultingFunction sets default MyKind field values
func (MyKindSchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*MyKind)
	// set default field values here
	log.Printf("Defaulting fields for MyKind %s\n", obj.Name)
}
