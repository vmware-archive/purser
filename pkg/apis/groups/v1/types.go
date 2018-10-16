package v1

import (
	"github.com/vmware/purser/pkg/controller/metrics"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CRD Group attributes
const (
	CRDPlural   string = "groups"
	CRDGroup    string = "vmware.kuber"
	CRDVersion  string = "v1"
	FullCRDName string = CRDPlural + "." + CRDGroup
)

// Group describes our custom Group resource
type Group struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               GroupSpec   `json:"spec"`
	Status             GroupStatus `json:"status,omitempty"`
}

// GroupSpec is the spec for the Group resource
type GroupSpec struct {
	Name               string                      `json:"name"`
	Type               string                      `json:"type,omitempty"`
	Labels             map[string]string           `json:"labels,omitempty"`
	AllocatedResources *metrics.Metrics            `json:"metrics,omitempty"`
	PodsMetrics        map[string]*metrics.Metrics `json:"pods,omitempty"`
	PodsDetails        map[string]*PodDetails      `json:"podDetails,omitempty"`
}

// GroupList is the list of Group resources
type GroupList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []*Group `json:"items"`
}

// GroupStatus holds the status information for each Group resource
type GroupStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// PodDetails information for the pods associated with the Group resource
type PodDetails struct {
	Name            string
	StartTime       meta_v1.Time
	EndTime         meta_v1.Time
	Containers      []*Container
	PodVolumeClaims map[string]*PersistentVolumeClaim
}

// PersistentVolumeClaim information for the pods associated with the Group resource
// A PVC can bound and unbound to a pod many times, so maintaining
// BoundTimes and UnboundTimes as lists.
// A PVC can be upgraded or downgraded, so maintaining capacityAllocated as a list
// Whenever a PVC capacity changes will update UnboundTime for old capacity, and
// append new capacity to capacityAllocated with bound time appended to BoundTimes
// The i-th capacity alloacted corresponds to the i-th bound time and to i-th unbound time.
// Similarly for RequestSizeInGB
type PersistentVolumeClaim struct {
	Name                string
	VolumeName          string
	RequestSizeInGB     []float64
	CapacityAllotedInGB []float64
	BoundTimes          []meta_v1.Time
	UnboundTimes        []meta_v1.Time
}

// Container information for the pods associated with the Group resource
type Container struct {
	Name    string
	Metrics *metrics.Metrics
}
