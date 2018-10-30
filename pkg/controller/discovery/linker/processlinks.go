package linker

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

// Process holds the details for the executing processes inside the container
type Process struct {
	ID, Name string
}

func storeProcessInteractions(containerProcessInteraction map[string][]string, processPodInteraction map[string](map[string]bool), creationTime time.Time) {
	for containerXID, procsXIDs := range containerProcessInteraction {
		for _, procXID := range procsXIDs {
			podsXIDs := []string{}
			for podXID := range processPodInteraction[procXID] {
				podsXIDs = append(podsXIDs, podXID)
			}
			// fetch the 4th field from ns : podName : containerName : procID : procName
			procName := strings.Split(procXID, KeySpliter)[4]
			err := models.StoreProcess(procName, containerXID, podsXIDs, creationTime)
			if err != nil {
				log.Errorf("failed to store process details: %s", procXID)
			}
		}
		err := models.StoreContainerProcessEdge(containerXID, procsXIDs)
		if err != nil {
			log.Errorf("failed to store edge from container: %s to procs", containerXID)
		}
	}
}

func populateContainerProcessTable(containerXID, procXID string) map[string][]string {
	containerProcessInteraction := make(map[string][]string)
	if _, isPresent := containerProcessInteraction[containerXID]; !isPresent {
		containerProcessInteraction[containerXID] = []string{}
	}
	containerProcessInteraction[containerXID] = append(containerProcessInteraction[containerXID], procXID)
	return containerProcessInteraction
}

func updatePodProcessInteractions(procXID, dstName string, processPodInteraction map[string](map[string]bool)) {
	if dstName != "" {
		if _, isPresent := processPodInteraction[procXID]; !isPresent {
			processPodInteraction[procXID] = make(map[string]bool)
		}
		processPodInteraction[procXID][dstName] = true
	}
}
