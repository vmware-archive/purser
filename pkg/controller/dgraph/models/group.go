package models

import (
	"time"

	groups_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	"github.com/vmware/purser/pkg/controller/dgraph"
)

// Dgraph Model Constants
const (
	IsPurserGroup = "isPurserGroup"
)

// GroupCRD schema in dgraph
type GroupCRD struct {
	dgraph.ID
	IsPurserGroup bool   `json:"isPurserGroup,omitempty"`
	Name          string `json:"name,omitempty"`
	StartTime     string `json:"startTime,omitempty"`
	EndTime       string `json:"endTime,omitempty"`
	Type          string `json:"type,omitempty"`
}

func createGroupCRDObject(group groups_v1.Group) GroupCRD {
	newGroup := GroupCRD{
		Name:          group.Name,
		IsPurserGroup: true,
		Type:          "purser",
		ID:            dgraph.ID{Xid: group.Name},
		StartTime:     group.GetCreationTimestamp().Time.Format(time.RFC3339),
	}

	deletionTimestamp := group.GetDeletionTimestamp()
	if !deletionTimestamp.IsZero() {
		newGroup.EndTime = deletionTimestamp.Time.Format(time.RFC3339)
	}
	return newGroup
}

// StoreGroupCRD create a new persistent volume in the Dgraph and updates if already present.
func StoreGroupCRD(group groups_v1.Group) (string, error) {
	xid := group.Name
	uid := dgraph.GetUID(xid, IsPurserGroup)

	newGroup := createGroupCRDObject(group)
	if uid != "" {
		newGroup.UID = uid
	}
	assigned, err := dgraph.MutateNode(newGroup, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	return assigned.Uids["blank-0"], nil
}
