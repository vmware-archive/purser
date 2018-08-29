package eventprocessor

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/purser_controller/config"
)

func ProcessEvents(conf *config.Config) {

	//TODO: listen for subscriber crd updates and update in memory copy.
	subscribers := getSubscribers(conf)

	for true {
		conf.RingBuffer.PrintDetails()

		for true {
			data, size := conf.RingBuffer.ReadN(READ_SIZE)

			if size == 0 {
				log.Debug("There are no events to process.")
				break
			}

			// Post data to subscribers.
			NotifySubscribers(data, subscribers)

			// Update groups
			//controller.UpdateCustomGroups(conf.Groupcrdclient, data)

		}
		time.Sleep(10 * time.Second)
	}
}
