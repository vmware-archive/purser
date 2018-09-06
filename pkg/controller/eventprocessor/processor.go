package eventprocessor

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/config"
	"github.com/vmware/purser/pkg/controller/controller"
)

// ProcessEvents processes the event and notifies the subscribers.
func ProcessEvents(conf *config.Config) {

	for {
		conf.RingBuffer.PrintDetails()

		for {
			// TODO: listen for subscriber and group crd updates and update
			// in memory copy instead of querying everytime.
			subscribers := getSubscribers(conf)
			groups := controller.GetAllGroups(conf.Groupcrdclient)

			if len(subscribers) == 0 && len(groups) == 0 {
				// Avoid processing events if there are no groups and subscribers.
				break
			}

			data, size := conf.RingBuffer.ReadN(ReadSize)

			if size == 0 {
				log.Debug("There are no events to process.")
				break
			}

			// Post data to subscribers.
			NotifySubscribers(data, subscribers)

			// Update user created groups.
			controller.UpdateCustomGroups(data, groups, conf.Groupcrdclient)

			conf.RingBuffer.RemoveN(size)
			conf.RingBuffer.PrintDetails()
		}
		time.Sleep(10 * time.Second)
	}
}
