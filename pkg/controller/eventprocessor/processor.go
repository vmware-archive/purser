package eventprocessor

import (
	"encoding/json"
	"time"

	"github.com/vmware/purser/pkg/controller"
	"github.com/vmware/purser/pkg/controller/persistence/dgraph"

	log "github.com/Sirupsen/logrus"
	api_v1 "k8s.io/api/core/v1"
)

// ProcessEvents processes the event and notifies the subscribers.
func ProcessEvents(conf *controller.Config) {

	for {
		conf.RingBuffer.PrintDetails()

		for {
			// TODO: listen for subscriber and group crd updates and update
			// in memory copy instead of querying everytime.
			subscribers := getSubscribers(conf)
			groups := getAllGroups(conf.Groupcrdclient)

			data, size := conf.RingBuffer.ReadN(ReadSize)

			if size == 0 {
				log.Debug("There are no events to process.")
				break
			}

			// Persist in dgraph
			//PersistPayloads(data)

			// Post data to subscribers.
			notifySubscribers(data, subscribers)

			// Update user created groups.
			updateCustomGroups(data, groups, conf.Groupcrdclient)

			conf.RingBuffer.RemoveN(size)
			conf.RingBuffer.PrintDetails()
		}
		time.Sleep(10 * time.Second)
	}
}

// PersistPayloads store payload info in dgraph
func PersistPayloads(payloads []*interface{}) {
	for _, event := range payloads {
		payload := (*event).(*controller.Payload)
		if payload.ResourceType == "Pod" {
			pod := api_v1.Pod{}
			err := json.Unmarshal([]byte(payload.Data), &pod)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			err = dgraph.PersistPod(pod)
			if err != nil {
				log.Errorf("Error while persisting pod %v", err)
			}
		} else if payload.ResourceType == "Service" {
			service := api_v1.Service{}
			err := json.Unmarshal([]byte(payload.Data), &service)
			if err != nil {
				log.Errorf("Error un marshalling payload " + payload.Data)
			}
			err = dgraph.PersistService(service)
			if err != nil {
				log.Errorf("Error while persisting service %v", err)
			}
		}
	}
}
