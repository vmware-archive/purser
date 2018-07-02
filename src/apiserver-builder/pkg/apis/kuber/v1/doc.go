
/*
 * licensed to vmware.
*/


// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=apiserver-builder/pkg/apis/kuber
// +k8s:defaulter-gen=TypeMeta
// +groupName=kuber.kuber
package v1 // import "apiserver-builder/pkg/apis/kuber/v1"

