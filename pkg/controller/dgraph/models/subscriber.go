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
	IsSubscriberCRD bool      `json:"isSubscriberCRD,omitempty"`
	Name            string    `json:"name,omitempty"`
	StartTime       time.Time `json:"startTime,omitempty"`
	EndTime         time.Time `json:"endTime,omitempty"`
	Type            string    `json:"type,omitempty"`
}

func createSubscriberCRDObject(subscriber subscribers_v1.Subscriber) SubscriberCRD {
	newSubscriber := SubscriberCRD{
		Name:            subscriber.Name,
		IsSubscriberCRD: true,
		Type:            "kuber.input",
		ID:              dgraph.ID{Xid: subscriber.Name},
		StartTime:       subscriber.GetCreationTimestamp().Time,
	}

	deletionTimestamp := subscriber.GetDeletionTimestamp()
	if !deletionTimestamp.IsZero() {
		newSubscriber.EndTime = deletionTimestamp.Time
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
