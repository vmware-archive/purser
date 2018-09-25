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

package database_test

import (
	"testing"
	"time"

	db "github.com/vmware/purser/pkg/database"
	"github.com/vmware/purser/test/utils"
	bolt "go.etcd.io/bbolt"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TestPodDetails struct {
	Name       string
	StartTime  meta_v1.Time
	EndTime    meta_v1.Time
	Containers []*TestContainer
}

type TestContainer struct {
	Name       string
	Allocation TestMetrics
}

type TestMetrics struct {
	CPU    float64
	Memory float64
}

const (
	dbName     = "purser_test.db"
	bucketName = "test_bucket"
	testKey1   = "test_key_1"
)

func getTestData() TestPodDetails {
	cMetrics := TestMetrics{
		CPU:    8.4,
		Memory: 243.7,
	}

	container := &TestContainer{
		Name:       "container-1",
		Allocation: cMetrics,
	}

	podDetails := TestPodDetails{
		Name:       "pod-1",
		StartTime:  meta_v1.Now(),
		EndTime:    meta_v1.Now(),
		Containers: []*TestContainer{container},
	}
	return podDetails
}

func TestPutAndGet(t *testing.T) {
	boltDB, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer boltDB.Close()
	utils.Ok(t, err)

	key := testKey1
	testValue1 := getTestData()
	err = db.PutInDB(boltDB, bucketName, key, testValue1)
	utils.Ok(t, err)

	key = testKey1
	outValue := &TestPodDetails{}
	err = db.GetValueFromDB(boltDB, bucketName, key, outValue)
	checkEqualForTestPodDetails(t, testValue1, *outValue)
}

func checkEqualForTestPodDetails(t *testing.T, obj1, obj2 TestPodDetails) {
	utils.Equals(t, obj1.Name, obj2.Name)
	utils.Equals(t, obj1.Containers, obj2.Containers)

	// Marshalling and Unmarshalling of type Time doesn't give exact value, but
	// gives very close value. Precision upto 1 second.
	deltaS := obj1.StartTime.Time.Sub(obj2.StartTime.Time).Seconds()
	deltaD := obj1.EndTime.Time.Sub(obj2.EndTime.Time).Seconds()
	utils.Assert(t, deltaS < 1, "start times didn't match")
	utils.Assert(t, deltaD < 1, "end times didnt' match")
}
