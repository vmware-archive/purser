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
		podToSvcTable[pod.Name] = append(podToSvcTable[pod.Name], serviceKey)
		podKey := pod.Namespace + KeySpliter + pod.Name
		svcToPodTable[serviceKey] = append(svcToPodTable[serviceKey], podKey)
	}
	models.StorePodServiceEdges(svcToPodTable)
}

// GenerateAndStoreSvcInteractions parses through pod interactions and generates a source to // destination service interaction.
func GenerateAndStoreSvcInteractions() {
	duplicateSvcMapChecker := make(map[string](map[string]bool))
	svcMap := make(map[string][]string)

	pods, err := models.RetrieveAllPods()
	if err != nil {
		log.Errorf("Unable to fetch pods: %s\n", err)
	}

	for _, pod := range pods {
		srcSvcs := podToSvcTable[pod.Name]
		for _, dstPod := range pod.Interacts {
			dstSrvcs := podToSvcTable[dstPod.Name]
			for _, srcSvc := range srcSvcs {
				for _, dstSvc := range dstSrvcs {
					if _, ok := duplicateSvcMapChecker[srcSvc]; !ok {
						duplicateSvcMapChecker[srcSvc] = make(map[string]bool)
					}
					if _, isPresent := duplicateSvcMapChecker[srcSvc][dstSvc]; !isPresent {
						svcMap[srcSvc] = append(svcMap[srcSvc], dstSvc)
						duplicateSvcMapChecker[srcSvc][dstSvc] = true
					}
				}
			}
		}
	}
	storeServiceInteractions(svcMap)
}

func storeServiceInteractions(svcMap map[string][]string) {
	for srcSvc, dstSvcs := range svcMap {
		models.StoreServicesInteraction(srcSvc, dstSvcs)
	}
}
