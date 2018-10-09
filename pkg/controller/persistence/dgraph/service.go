package dgraph

import (
	"encoding/json"
	"time"

	api_v1 "k8s.io/api/core/v1"
)

type Service struct {
	ID
	IsService string    `json:"isService"`
	Name      string    `json:"name,omitempty"`
	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
	Pod       []*Pod    `json:"servicePods,omitempty"`
}

func PersistService(service api_v1.Service) error {
	xid := service.Namespace + ":" + service.Name
	uid, _ := GetUId(Client, xid, "isService")

	if uid == "" {
		// If pod is not present, persist it.
		newService := Service{
			Name:      service.Name,
			IsService: "",
			ID:        ID{Xid: xid},
		}
		bytes, err := json.Marshal(newService)
		if err != nil {
			return err
		}
		return MutateNode(Client, bytes)
	}
	return nil
}
