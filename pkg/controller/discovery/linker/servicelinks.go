package linker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph/models"

	corev1 "k8s.io/api/core/v1"
)

var (
	podToSvcTable = make(map[string][]string)
	svcToPodTable = make(map[string][]string)
)

// PopulatePodToServiceTable populates the pod<->service map
func PopulatePodToServiceTable(svc corev1.Service, pods *corev1.PodList) {
	serviceKey := svc.Namespace + KeySpliter + svc.Name
	for _, pod := range pods.Items {
		podKey := pod.Namespace + KeySpliter + pod.Name
		podToSvcTable[podKey] = append(podToSvcTable[podKey], serviceKey)
		svcToPodTable[serviceKey] = append(svcToPodTable[serviceKey], podKey)
	}
	models.StorePodServiceEdges(svcToPodTable)
}

// GenerateAndStoreSvcInteractions parses through pod interactions and generates a source to // destination service interaction.
func GenerateAndStoreSvcInteractions() {
	services, err := models.RetrieveAllServicesWithDstPods()
	if err != nil {
		log.Errorf("Unable to fetch services: %s\n", err)
	}

	for _, service := range services {
		destinationPods := getDestinationPods(service.Pod)
		destinationServices := getServicesXIDsFromPods(destinationPods)
		err = models.StoreServicesInteraction(service.Xid, destinationServices)
		if err != nil {
			log.Error(err)
		}
	}
}

func getDestinationPods(podsInService []*models.Pod) []*models.Pod {
	var destinationPods []*models.Pod
	for _, pod := range podsInService {
		destinationPods = append(destinationPods, pod.Interacts...)
	}
	return destinationPods
}

func getServicesXIDsFromPods(pods []*models.Pod) []string {
	var servicesXIDs []string
	duplicateChecker := make(map[string]bool)
	for _, pod := range pods {
		svcsXIDs := podToSvcTable[pod.Xid]
		for _, svcXID := range svcsXIDs {
			if _, isPresent := duplicateChecker[svcXID]; !isPresent {
				duplicateChecker[svcXID] = true
				servicesXIDs = append(servicesXIDs, svcXID)
			}
		}
	}
	return servicesXIDs
}
