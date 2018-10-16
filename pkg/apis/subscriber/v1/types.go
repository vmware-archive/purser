package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CRD Subscriber attributes
const (
	SubscriberPlural   string = "subscribers"
	SubscriberGroup    string = "kuber.input"
	SubscriberVersion  string = "v1"
	SubscriberFullName string = SubscriberPlural + "." + SubscriberGroup
)

// Subscriber information
type Subscriber struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               SubscriberSpec   `json:"spec"`
	Status             SubscriberStatus `json:"status,omitempty"`
}

// SubscriberSpec definition details
type SubscriberSpec struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster"`
	OrgID       string `json:"orgId"`
	URL         string `json:"url"`
	AuthType    string `json:"authType,omitempty"`
	AuthToken   string `json:"authToken,omitempty"`
}

// SubscriberStatus definition
type SubscriberStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// SubscriberList type
type SubscriberList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []Subscriber `json:"items"`
}
