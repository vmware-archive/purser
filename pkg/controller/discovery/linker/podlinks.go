package linker

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"

	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

// podIPTable: maps pod name with pod IP address
// podToPodTable: maps src pod to the interacting dest pod along with the interaction frequency count.
var (
	podIPTable    = make(map[string]string)
	podToPodTable = make(map[string](map[string]float64))
)

// Process holds the details for the executing processes inside the container
type Process struct {
	ID, Name string
}

// PopulatePodIPTable populates the podIP<->podName map
func PopulatePodIPTable(pods *corev1.PodList) {
	for _, pod := range pods.Items {
		podName := pod.GetName()
		podIP := pod.Status.PodIP
		podIPTable[podIP] = pod.Namespace + ":" + podName
	}
}

// GenerateAndStorePodInteractions generates source to destination Pod mapping and stores it in Dgraph.
func GenerateAndStorePodInteractions() {
	for srcPodName, communication := range podToPodTable {
		dstPods := []string{}
		counts := []float64{}
		for dstPodName, count := range communication {
			dstPods = append(dstPods, dstPodName)
			counts = append(counts, count)
		}
		err := models.StorePodsInteraction(srcPodName, dstPods, counts)
		if err != nil {
			log.Errorf("failed to store pod interaction in Dgraph %v", err)
		}
	}
}

// PopulateMappingTables updates PodToPodTable
func PopulateMappingTables(tcpDump []string, pod corev1.Pod, containerName string) {
	for _, address := range tcpDump {
		address := strings.Split(address, ":")
		srcIP, dstIP := address[0], address[2]
		srcName, dstName := podIPTable[srcIP], podIPTable[dstIP]
		updatePodToPodTable(srcName, dstName)
	}
}

func updatePodToPodTable(srcName, dstName string) {
	if dstName != "" && srcName != "" {
		if _, ok := podToPodTable[srcName]; !ok {
			podToPodTable[srcName] = make(map[string]float64)
		}

		if _, isPresent := podToPodTable[srcName][dstName]; !isPresent {
			podToPodTable[srcName][dstName] = 1
		} else {
			podToPodTable[srcName][dstName]++
		}
	}
}
