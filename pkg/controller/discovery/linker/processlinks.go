/*
 * Copyright (c) 2018 VMware Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
