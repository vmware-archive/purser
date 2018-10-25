package processor

import (
	"fmt"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/purser/pkg/controller"
	"github.com/vmware/purser/pkg/controller/discovery/linker"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

var svcwg sync.WaitGroup

// ProcessServiceInteractions parses through the list of services and it's associated pods to
// generate a 1:1 mapping between the communicating services.
func ProcessServiceInteractions(conf controller.Config) {
	services := RetrieveServiceList(conf.Kubeclient, metav1.ListOptions{})

	processServiceDetails(conf.Kubeclient, services)
	linker.GenerateAndStoreSvcInteractions()

	log.Infof("Successfully generated Service To Service mapping.")
}

func processServiceDetails(client *kubernetes.Clientset, services *corev1.ServiceList) {
	svcCount := len(services.Items)
	log.Infof("Processing total of (%d) Services.", svcCount)

	svcwg.Add(svcCount)
	{
		for index, svc := range services.Items {
			log.Debugf("Processing Service (%d/%d): %s ", index+1, svcCount, svc.GetName())

			go func(svc corev1.Service, index int) {
				defer svcwg.Done()

				selectorSet := labels.Set(svc.Spec.Selector)
				if selectorSet != nil {
					options := metav1.ListOptions{
						LabelSelector: selectorSet.AsSelector().String(),
					}
					pods := RetrievePodList(client, options)
					linker.PopulatePodToServiceTable(svc, pods)
				}

				log.Debugf("Finished processing Service (%d/%d)", index+1, svcCount)
			}(svc, index)
		}
	}
	svcwg.Wait()
	fmt.Println("")
}
