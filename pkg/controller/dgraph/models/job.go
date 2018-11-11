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

	"github.com/vmware/purser/pkg/controller/dgraph"
	log "github.com/Sirupsen/logrus"
	batch_v1 "k8s.io/api/batch/v1"
)

// Dgraph Model Constants
const (
	IsJob = "isJob"
)

// Daemonset schema in dgraph
type Job struct {
	dgraph.ID
	IsJob bool       `json:"isJob,omitempty"`
	Name          string     `json:"name,omitempty"`
	StartTime     string  `json:"startTime,omitempty"`
	EndTime       string  `json:"endTime,omitempty"`
	Namespace     *Namespace `json:"namespace,omitempty"`
	Pods          []*Pod     `json:"pods,omitempty"`
	Type          string     `json:"type,omitempty"`
}

func createJobObject(job batch_v1.Job) Job {
	newJob := Job{
		Name:          job.Name,
		IsJob: true,
		Type:          "job",
		ID:            dgraph.ID{Xid: job.Namespace + ":" + job.Name},
		StartTime:     job.GetCreationTimestamp().Time.Format(time.RFC3339),
	}
	namespaceUID := CreateOrGetNamespaceByID(job.Namespace)
	if namespaceUID != "" {
		newJob.Namespace = &Namespace{ID: dgraph.ID{UID: namespaceUID, Xid: job.Namespace}}
	}
	jobDeletionTimestamp := job.GetDeletionTimestamp()
	if !jobDeletionTimestamp.IsZero() {
		newJob.EndTime = jobDeletionTimestamp.Time.Format(time.RFC3339)
	}
	return newJob
}

// StoreJob create a new daemonset in the Dgraph and updates if already present.
func StoreJob(job batch_v1.Job) (string, error) {
	xid := job.Namespace + ":" + job.Name
	uid := dgraph.GetUID(xid, IsJob)

	newJob := createJobObject(job)
	if uid != "" {
		newJob.UID = uid
	}
	assigned, err := dgraph.MutateNode(newJob, dgraph.CREATE)
	if err != nil {
		return "", err
	}
	return assigned.Uids["blank-0"], nil
}

// CreateOrGetJobByID returns the uid of namespace if exists,
// otherwise creates the job and returns uid.
func CreateOrGetJobByID(xid string) string {
	if xid == "" {
		return ""
	}
	uid := dgraph.GetUID(xid, IsJob)

	if uid != "" {
		return uid
	}

	d := Job{
		ID:            dgraph.ID{Xid: xid},
		Name:          xid,
		IsJob: true,
	}
	assigned, err := dgraph.MutateNode(d, dgraph.CREATE)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return assigned.Uids["blank-0"]
}

// RetrieveAllJobs ...
func RetrieveAllJobs() ([]byte, error) {
	const q = `query {
		job(func: has(isJob)) {
			name
			type
			pod: ~job @filter(has(isPod) {
				name
				type
				container: ~pod @filter(has(isContainer)) {
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

// RetrieveJob ...
func RetrieveJob(name string) ([]byte, error) {
	q := `query {
		job(func: has(isJob)) @filter(eq(name, "` + name + `")) {
			name
			type
			pod: ~job @filter(has(isPod)) {
				name
				type
				container: ~pod @filter(has(isContainer)) {
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
