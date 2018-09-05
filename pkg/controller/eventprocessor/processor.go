package eventprocessor

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/config"
)

// ProcessEvents processes the event and notifies the subscribers.
func ProcessEvents(conf *config.Config) {

	//TODO: listen for subscriber crd updates and update in memory copy.
	subscribers := getSubscribers(conf)

	for {
		conf.RingBuffer.PrintDetails()

		for {
			data, size := conf.RingBuffer.ReadN(ReadSize)

			if size == 0 {
				log.Debug("There are no events to process.")
				break
			}

			// Post data to subscribers.
			NotifySubscribers(data, subscribers)
		}
		time.Sleep(10 * time.Second)
	}
}
