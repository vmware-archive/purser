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

package models

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/purser/pkg/controller/dgraph"

	api_v1 "k8s.io/api/core/v1"
)

// Dgraph Model Constants
const (
	IsNamespace = "isNamespace"
	defaultCPUCostPerCPUPerHour    = "0.024"
	defaultMemCostPerGBPerHour     = "0.01"
	defaultStorageCostPerGBPerHour = "0.00013888888"
)

// Namespace schema in dgraph
type Namespace struct {
	dgraph.ID
	IsNamespace bool      `json:"isNamespace,omitempty"`
	Name        string    `json:"name,omitempty"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
	Type        string    `json:"type,omitempty"`
	Deployments []*Deployment `json:"deployment,omitempty"`
	Statefulsets []*Statefulset `json:"statefulset,omitempty"`
	Jobs []*Job `json:"job,omitempty"`
	Daemonsets []*Daemonset `json:"daemonset,omitempty"`
	Replicasets []*Replicaset `json:"replicaset,omitempty"`
}

type Children struct {
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	CPU    float64    `json:"cpu,omitempty"`
	Memory float64    `json:"memory,omitempty"`
	Storage float64    `json:"storage,omitempty"`
	CPUCost float64    `json:"cpuCost,omitempty"`
	MemoryCost float64    `json:"memoryCost,omitempty"`
	StorageCost float64    `json:"storageCost,omitempty"`
}

type Parent struct {
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	Children []Children  `json:"children,omitempty"`
	CPU    float64    `json:"cpu,omitempty"`
	Memory float64    `json:"memory,omitempty"`
	Storage float64    `json:"storage,omitempty"`
	CPUCost float64    `json:"cpuCost,omitempty"`
	MemoryCost float64    `json:"memoryCost,omitempty"`
	StorageCost float64    `json:"storageCost,omitempty"`
}

type ParentWrapper struct {
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	Children []Children `json:"children,omitempty"`
	Parent []Parent  `json:"parent,omitempty"`
	CPU float64    `json:"cpu,omitempty"`
	Memory float64    `json:"memory,omitempty"`
	Storage float64    `json:"storage,omitempty"`
	CPUCost float64    `json:"cpuCost,omitempty"`
	MemoryCost float64    `json:"memoryCost,omitempty"`
	StorageCost float64    `json:"storageCost,omitempty"`
}

type JsonDataWrapper struct {
	Data ParentWrapper `json:"data,omitempty"`
}

func newNamespace(namespace api_v1.Namespace) Namespace {
	ns := Namespace{
		ID:          dgraph.ID{Xid: namespace.Name},
		Name:        "namespace-" + namespace.Name,
		IsNamespace: true,
		Type:        "namespace",
		StartTime:   namespace.GetCreationTimestamp().Time.Format(time.RFC3339),
	}
	nsDeletionTimestamp := namespace.GetDeletionTimestamp()
	if !nsDeletionTimestamp.IsZero() {
		ns.EndTime = nsDeletionTimestamp.Time.Format(time.RFC3339)
	}
	return ns
}

// CreateOrGetNamespaceByID returns the uid of namespace if exists,
// otherwise creates the namespace and returns uid.
func CreateOrGetNamespaceByID(xid string) string {
	if xid == "" {
		log.Error("Namespace is empty")
		return ""
	}
	uid := dgraph.GetUID(xid, IsNamespace)

	if uid != "" {
		return uid
	}

	ns := Namespace{
		ID:          dgraph.ID{Xid: xid},
		Name:        xid,
		IsNamespace: true,
	}
	assigned, err := dgraph.MutateNode(ns, dgraph.CREATE)
	if err != nil {
		log.Error(err)
		return ""
	}
	log.Infof("Namespace with xid: (%s) persisted", xid)
	return assigned.Uids["blank-0"]
}

// StoreNamespace create a new namespace in the Dgraph  if it is not present.
func StoreNamespace(namespace api_v1.Namespace) (string, error) {
	xid := namespace.Name
	uid := dgraph.GetUID(xid, IsNamespace)

	ns := newNamespace(namespace)
	if uid != "" {
		ns.UID = uid
	}
	assigned, err := dgraph.MutateNode(ns, dgraph.CREATE)
	if err != nil {
		return "", err
	}

	if uid == "" {
		log.Infof("Namespace with xid: (%s) persisted", xid)
	}
	return assigned.Uids["blank-0"], nil
}

// RetrieveAllNamespaces ...
func RetrieveAllNamespaces() ([]byte, error) {
	const q = `query {
		cluster(func: has(isNamespace)) {
			name
			type
			deployment: ~namespace @filter(has(isDeployment)) {
				name
				type
				replicaset: ~deployment @filter(has(isReplicaset)) {
					name
					type
				}
			}
			statefulset: ~namespace @filter(has(isStatefulset)) {
				name
				type
				statefulset: ~statefulset @filter(has(isPod)) {
					name
					type
				}
			}
			job: ~namespace @filter(has(isJob)) {
				name
				type
				job: ~job @filter(has(isPod)) {
					name
					type
				}
			}
			daemonset: ~namespace @filter(has(isDaemonset)) {
				name
				type
				pod: ~daemonset @filter(has(isPod)) {
					name
					type
				}
			}
			replicaset: ~namespace @filter(has(isReplicaset) AND (NOT has(deployment))) {
				name
				type
				pod: ~replicaset @filter(has(isPod)) {
					name
					type
				}
			}
		}
	}`

	result, err := dgraph.ExecuteQueryRaw(q)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RetrieveNamespace ...
func RetrieveNamespace(name string) ([]byte, error) {
	q := `query {
		namespace(func: has(isNamespace)) @filter(eq(name, "` + name + `")) {
			name
			type
			deployment: ~namespace @filter(has(isDeployment)) {
				name
				type
				replicaset: ~deployment @filter(has(isReplicaset)) {
					name
					type
				}
			}
			statefulset: ~namespace @filter(has(isStatefulset)) {
				name
				type
				~statefulset @filter(has(isPod)) {
					name
					type
				}
			}
			job: ~namespace @filter(has(isJob)) {
				name
				type
				pod: ~job @filter(has(isPod)) {
					name
					type
				}
			}
			daemonset: ~namespace @filter(has(isDaemonset)) {
				name
				type
				pod: ~daemonset @filter(has(isPod)) {
					name
					type
				}
			}
			replicaset: ~namespace @filter(has(isReplicaset) AND (NOT has(deployment))) {
				name
				type
				pod: ~replicaset @filter(has(isPod)) {
					name
					type
				}
			}
		}
	}`


	result, err := dgraph.ExecuteQueryRaw(q)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RetrieveAllNamespacesWithMetrics ...
func RetrieveAllNamespacesWithMetrics() (JsonDataWrapper, error) {
	const q = `query {
		ns as var(func: has(isNamespace)) {
			~namespace @filter(has(isPod)){
				namespacePodCpu as cpuRequest
				namespacePodMem as memoryRequest
				namespacePvcStorage as storageRequst
			}
			namespaceCpu as sum(val(namespacePodCpu))
			namespaceMem as sum(val(namespacePodMem))
			namespaceStorage as sum(val(namespacePvcStorage))
        }

		children(func: uid(ns)) {
			name
            type
			cpu: val(namespaceCpu)
			memory: val(namespaceMem)
			storage: val(namespaceStorage)
			cpuCost: math(namespaceCpu * ` + defaultCPUCostPerCPUPerHour + `)
			memoryCost: math(namespaceMem * ` + defaultMemCostPerGBPerHour + `)
			storageCost: math(namespaceStorage * ` + defaultStorageCostPerGBPerHour + `)
        }
    }`
	parentRoot := ParentWrapper{}
	err := dgraph.ExecuteQuery(q, &parentRoot)
	calculateTotal(&parentRoot)
	root := JsonDataWrapper{}
	root.Data = ParentWrapper{
		Name: "cluster",
		Type: "cluster",
		Children: parentRoot.Children,
		CPU: parentRoot.CPU,
		Memory: parentRoot.Memory,
		Storage: parentRoot.Storage,
		CPUCost: parentRoot.CPUCost,
		MemoryCost: parentRoot.MemoryCost,
		StorageCost: parentRoot.StorageCost,
	}
	log.Debugf("data: (%v)", root.Data)
	return root, err
}

// RetrieveNamespaceWithMetrics ...
func RetrieveNamespaceWithMetrics(name string) (JsonDataWrapper, error) {
	q := `query {
		ns as var(func: has(isNamespace)) @filter(eq(name, "` + name + `")) {
			childs as ~namespace @filter(has(isDeployment) OR has(isStatefulset) OR has(isJob) OR has(isDaemonset) OR (has(isReplicaset) AND (NOT has(deployment)))) {
				name
				type
				~deployment @filter(has(isReplicaset)) {
                    name
                    type
					~replicaset @filter(has(isPod)) {
						name
						type
				        replicasetPodCpu as cpuRequest
				        replicasetPodMemory as memoryRequest
						replicasetPvcStorage as storageRequest
			        }
    	            deploymentReplicasetCpu as sum(val(replicasetPodCpu))
			        deploymentReplicasetMemory as sum(val(replicasetPodMemory))
					deploymentReplicasetStorage as sum(val(replicasetPvcStorage))
                }
				~statefulset @filter(has(isPod)) {
                    name
                    type
                    statefulsetPodCpu as cpuRequest
                    statefulsetPodMemory as memoryRequest
					statefulsetPvcStorage as storageRequest
                }
				~job @filter(has(isPod)) {
                    name
                    type
                    jobPodCpu as cpuRequest
                    jobPodMemory as memoryRequest
					jobPvcStorage as jobRequest
                }
				~daemonset @filter(has(isPod)) {
                    name
                    type
                    daemonsetPodCpu as cpuRequest
                    daemonsetPodMemory as memoryRequest
					daemonsetPvcStorage as daemonsetRequest
                }
				~replicaset @filter(has(isPod)) {
                    name
                    type
                    replicasetSimplePodCpu as cpuRequest
                    replicasetSimplePodMemory as memoryRequest
					replicasetSimplePvcStorage as replicasetRequest
                }
				sumReplicasetSimplePodCpu as sum(val(replicasetSimplePodCpu))
				sumDaemonsetPodCpu as sum(val(daemonsetPodCpu))
				sumJobPodCpu as sum(val(jobPodCpu))
				sumStatefulsetPodCpu as sum(val(statefulsetPodCpu))
				sumDeploymentPodCpu as sum(val(deploymentReplicasetCpu))
				namespaceChildCpu as math(sumReplicasetSimplePodCpu + sumDaemonsetPodCpu + sumJobPodCpu + sumStatefulsetPodCpu + sumDeploymentPodCpu)

				sumReplicasetSimplePodMemory as sum(val(replicasetSimplePodMemory))
				sumDaemonsetPodMemory as sum(val(daemonsetPodMemory))
				sumJobPodMemory as sum(val(jobPodMemory))
				sumStatefulsetPodMemory as sum(val(statefulsetPodMemory))
				sumDeploymentPodMemory as sum(val(deploymentReplicasetMemory))
				namespaceChildMemory as math(sumReplicasetSimplePodMemory + sumDaemonsetPodMemory + sumJobPodMemory + sumStatefulsetPodMemory + sumDeploymentPodMemory)
        	
				sumReplicasetSimplePvcStorage as sum(val(replicasetSimplePvcStorage))
				sumDaemonsetPvcStorage as sum(val(daemonsetPvcStorage))
				sumJobPvcStorage as sum(val(jobPvcStorage))
				sumStatefulsetPvcStorage as sum(val(statefulsetPvcStorage))
				sumDeploymentPvcStorage as sum(val(deploymentReplicasetStorage))
				namespaceChildStorage as math(sumReplicasetSimplePvcStorage + sumDaemonsetPvcStorage + sumJobPvcStorage + sumStatefulsetPvcStorage + sumDeploymentPvcStorage)
			}
			namespaceCpu as sum(val(namespaceChildCpu))
			namespaceMemory as sum(val(namespaceChildMemory))
			namespaceStorage as sum(val(namespaceChildStorage))
		}

		parent(func: uid(ns)) {
			name
            type
			children: ~namespace @filter(uid(childs)) {
				name
				type
				cpu: namespaceChildCpu
				memory: namespaceChildMemory
				storage: namespaceChildStorage
				cpuCost: math(namespaceChildCpu * ` + defaultCPUCostPerCPUPerHour + `)
				memoryCost: math(namespaceChildMemory * ` + defaultMemCostPerGBPerHour + `)
				storageCost: math(namespaceChildStorage * ` + defaultStorageCostPerGBPerHour + `)
			}
			cpu: namespaceCpu
			memory: namespaceMemory
			storage: namespaceStorage
			cpuCost: math(namespaceCpu * ` + defaultCPUCostPerCPUPerHour + `)
			memoryCost: math(namespaceMemory * ` + defaultMemCostPerGBPerHour + `)
			storageCost: math(namespaceStorage * ` + defaultStorageCostPerGBPerHour + `)
        }
    }`
	parentRoot := ParentWrapper{}
	err := dgraph.ExecuteQuery(q, &parentRoot)
	root := JsonDataWrapper{}
	if len(parentRoot.Parent) == 0 {
		return root, err
	}
	root.Data = ParentWrapper{
		Name: parentRoot.Parent[0].Name,
		Type: parentRoot.Parent[0].Type,
		Children: parentRoot.Parent[0].Children,
		CPU: parentRoot.Parent[0].CPU,
		Memory: parentRoot.Parent[0].Memory,
		Storage: parentRoot.Parent[0].Storage,
		CPUCost: parentRoot.Parent[0].CPUCost,
		MemoryCost: parentRoot.Parent[0].MemoryCost,
		StorageCost: parentRoot.Parent[0].StorageCost,
	}
	return root, err
}

func calculateTotal(objRoot *ParentWrapper) {
	for _, obj := range objRoot.Children {
		objRoot.CPU += obj.CPU
		objRoot.Memory += obj.Memory
		objRoot.Storage += obj.Storage
		objRoot.CPUCost += obj.CPUCost
		objRoot.MemoryCost += obj.MemoryCost
		objRoot.StorageCost += obj.StorageCost
	}
}