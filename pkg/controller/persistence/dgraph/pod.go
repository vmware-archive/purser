package dgraph

import (
	"encoding/json"
	"time"

	api_v1 "k8s.io/api/core/v1"
)

type ID struct {
	Xid string `json:"xid,omitempty"`
	UID string `json:"uid,omitempty"`
}

type Pod struct {
	ID
	IsPod      string       `json:"isPod"`
	Name       string       `json:"name,omitempty"`
	StartTime  time.Time    `json:"startTime,omitempty"`
	EndTime    time.Time    `json:"endTime,omitempty"`
	Containers []*Container `json:"containers,omitempty"`
}

type Container struct {
	ID
	IsContainer string    `json:"isContainer"`
	Name        string    `json:"name,omitempty"`
	StartTime   time.Time `json:"startTime,omitempty"`
	EndTime     time.Time `json:"endTime,omitempty"`
	Pod         Pod       `json:"pod,omitempty"`
}

func PersistPod(pod api_v1.Pod) error {
	xid := pod.Namespace + ":" + pod.Name
	uid, _ := GetUId(Client, xid, "isPod")

	if uid == "" {
		// If pod is not present, persist it.
		newPod := Pod{
			Name:  pod.Name,
			IsPod: "",
			ID:    ID{Xid: xid},
		}
		bytes, err := json.Marshal(newPod)
		if err != nil {
			return err
		}
		return MutateNode(Client, bytes)
	}
	return nil
}
