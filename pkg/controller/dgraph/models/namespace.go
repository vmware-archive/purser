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
	CPU    float64    `json:"cpu,omitempty"`
	Memory float64    `json:"memory,omitempty"`
	Children []*Children `json:"children,omitempty"`
}

type Children struct {
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	CPU    float64    `json:"cpu,omitempty"`
	Memory float64    `json:"memory,omitempty"`
}

// NamespacesWithMetrics ...
type NamespacesWithMetrics struct {
	Namespace []Namespace  `json:"namespace,omitempty"`
	CPU    float64    `json:"cpu,omitempty"`
	Memory float64    `json:"memory,omitempty"`
}

func newNamespace(namespace api_v1.Namespace) Namespace {
	ns := Namespace{
		ID:          dgraph.ID{Xid: namespace.Name},
		Name:        namespace.Name,
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
		namespace(func: has(isNamespace)) {
			name
			type
			deployment: ~namespace @filter(has(isDeployment)) {
				name
				type
				~deployment @filter(has(isReplicaset)) {
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
				~job @filter(has(isPod)) {
					name
					type
				}
			}
			daemonset: ~namespace @filter(has(isDaemonset)) {
				name
				type
				~daemonset @filter(has(isPod)) {
					name
					type
				}
			}
			replicaset: ~namespace @filter(has(isReplicaset) AND (NOT has(deployment))) {
				name
				type
				~replicaset @filter(has(isPod)) {
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
				~deployment @filter(has(isReplicaset)) {
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
				~job @filter(has(isPod)) {
					name
					type
				}
			}
			daemonset: ~namespace @filter(has(isDaemonset)) {
				name
				type
				~daemonset @filter(has(isPod)) {
					name
					type
				}
			}
			replicaset: ~namespace @filter(has(isReplicaset) AND (NOT has(deployment))) {
				name
				type
				~replicaset @filter(has(isPod)) {
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
func RetrieveAllNamespacesWithMetrics() (NamespacesWithMetrics, error) {
	const q = `query {
		ns as var(func: has(isNamespace)) {
			~namespace @filter(has(isPod)){
				namespacePodCpu as cpuRequest
				namespacePodMem as memoryRequest
			}
			namespaceCpu as sum(val(namespacePodCpu))
			namespaceMem as sum(val(namespacePodMem))
        }

		namespace(func: uid(ns)) {
			name
            type
			cpu: val(namespaceCpu)
			memory: val(namespaceMem)
        }
    }`
	namespaceRoot := NamespacesWithMetrics{}
	err := dgraph.ExecuteQuery(q, &namespaceRoot)
	calculateTotalNamespaceMetrics(&namespaceRoot)
	return namespaceRoot, err
}

// RetrieveNamespaceWithMetrics ...
func RetrieveNamespaceWithMetrics(name string) (NamespacesWithMetrics, error) {
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
			        }
    	            deploymentReplicasetCpu as sum(val(replicasetPodCpu))
			        deploymentReplicasetMemory as sum(val(replicasetPodMemory))
                }
				~statefulset @filter(has(isPod)) {
                    name
                    type
                    cpu: statefulsetPodCpu as cpuRequest
                    memory: statefulsetPodMemory as memoryRequest
                }
				~job @filter(has(isPod)) {
                    name
                    type
                    jobPodCpu as cpuRequest
                    jobPodMemory as memoryRequest
                }
				~daemonset @filter(has(isPod)) {
                    name
                    type
                    daemonsetPodCpu as cpuRequest
                    daemonsetPodMemory as memoryRequest
                }
				~replicaset @filter(has(isPod)) {
                    name
                    type
                    replicasetSimplePodCpu as cpuRequest
                    replicasetSimplePodMemory as memoryRequest
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
        	}
			namespaceCpu as sum(val(namespaceChildCpu))
			namespaceMemory as sum(val(namespaceChildMemory))
		}

		namespace(func: uid(ns)) {
			name
            type
			children: ~namespace @filter(uid(childs)) {
				name
				type
				cpu: val(namespaceChildCpu)
				memory: val(namespaceChildMemory)
			}
			cpu: val(namespaceCpu)
			memory: val(namespaceMemory)
        }
    }`
	namespaceRoot := NamespacesWithMetrics{}
	err := dgraph.ExecuteQuery(q, &namespaceRoot)
	calculateTotalNamespaceMetrics(&namespaceRoot)
	return namespaceRoot, err
}

func calculateTotalNamespaceMetrics(objRoot *NamespacesWithMetrics) {
	for _, obj := range objRoot.Namespace {
		objRoot.CPU += obj.CPU
		objRoot.Memory += obj.Memory
	}
}