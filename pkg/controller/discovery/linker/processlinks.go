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
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph/models"
)

// Process holds the details for the executing processes inside the container
type Process struct {
	ID, Name string
}

// StoreProcessInteractions stores process, container to process edge, process to pods edge
func StoreProcessInteractions(containerProcessInteraction map[string][]string, processPodInteraction map[string](map[string]bool), creationTime time.Time) {
	for containerXID, procsXIDs := range containerProcessInteraction {
		for _, procXID := range procsXIDs {
			podsXIDs := []string{}
			for podXID := range processPodInteraction[procXID] {
				podsXIDs = append(podsXIDs, podXID)
			}

			err := models.StoreProcess(procXID, containerXID, podsXIDs, creationTime)
			if err != nil {
				log.Errorf("failed to store process details: %s, err: (%v)", procXID, err)
			}
		}
		err := models.StoreContainerProcessEdge(containerXID, procsXIDs)
		if err != nil {
			log.Errorf("failed to store edge from container: %s to procs, err: (%v)", containerXID, err)
		}
	}
}

func populateContainerProcessTable(containerXID, procXID string, interactions *InteractionsWrapper) {
	if _, isPresent := interactions.ContainerProcessInteraction[containerXID]; !isPresent {
		interactions.ContainerProcessInteraction[containerXID] = []string{}
	}
	interactions.ContainerProcessInteraction[containerXID] = append(interactions.ContainerProcessInteraction[containerXID], procXID)
}

func updatePodProcessInteractions(procXID, dstName string, interactions *InteractionsWrapper) {
	if dstName != "" {
		if _, isPresent := interactions.ProcessToPodInteraction[procXID]; !isPresent {
			interactions.ProcessToPodInteraction[procXID] = make(map[string]bool)
		}
		interactions.ProcessToPodInteraction[procXID][dstName] = true
	}
}
