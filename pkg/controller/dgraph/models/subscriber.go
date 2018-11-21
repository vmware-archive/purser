package models

import (
	"time"

	subscribers_v1 "github.com/vmware/purser/pkg/apis/subscriber/v1"
	"github.com/vmware/purser/pkg/controller/dgraph"
)

// Dgraph Model Constants
const (
	IsSubscriberCRD = "isSubscriberCRD"
)

// SubscriberCRD schema in dgraph
type SubscriberCRD struct {
	dgraph.ID
	IsSubscriberCRD bool   `json:"isSubscriberCRD,omitempty"`
	Name            string `json:"name,omitempty"`
	StartTime       string `json:"startTime,omitempty"`
	EndTime         string `json:"endTime,omitempty"`
	Type            string `json:"type,omitempty"`
}

func createSubscriberCRDObject(subscriber subscribers_v1.Subscriber) SubscriberCRD {
	newSubscriber := SubscriberCRD{
		Name:            subscriber.Name,
		IsSubscriberCRD: true,
		Type:            "subscriber",
		ID:              dgraph.ID{Xid: subscriber.Name},
		StartTime:       subscriber.GetCreationTimestamp().Time.Format(time.RFC3339),
	}

	deletionTimestamp := subscriber.GetDeletionTimestamp()
	if !deletionTimestamp.IsZero() {
		newSubscriber.EndTime = deletionTimestamp.Time.Format(time.RFC3339)
	}
	return newSubscriber
}

// StoreSubscriberCRD create a new persistent volume in the Dgraph and updates if already present.
func StoreSubscriberCRD(subscriber subscribers_v1.Subscriber) (string, error) {
	xid := subscriber.Name
	uid := dgraph.GetUID(xid, IsSubscriberCRD)

	newSubscriber := createSubscriberCRDObject(subscriber)
	if uid != "" {
		newSubscriber.UID = uid
	}
	assigned, err := dgraph.MutateNode(newSubscriber, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	return assigned.Uids["blank-0"], nil
}
