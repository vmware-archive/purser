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

package buffering_test

import (
	"sync"
	"testing"

	"github.com/vmware/purser/pkg/purser_controller/buffering"
	"github.com/vmware/purser/test/utils"
)

func TestPut(t *testing.T) {
	// use Put to add one more, return from Put should be True
	r := &buffering.RingBuffer{Size: 2, Mutex: &sync.Mutex{}}

	testValue1 := 1
	ret1 := r.Put(testValue1)
	utils.Assert(t, ret1, "inserting into not full buffer")

	testValue2 := 38
	ret2 := r.Put(testValue2)
	utils.Assert(t, !ret2, "inserting into full buffer")
}

func TestGet(t *testing.T) {
	// use Put to add one more, return from Put should be True
	r := &buffering.RingBuffer{Size: 2, Mutex: &sync.Mutex{}}

	ret1 := r.Get()
	utils.Assert(t, ret1 == nil, "get elements of empty buffer")

	testValue := 1
	r.Put(testValue)
	ret2 := r.Get()
	utils.Assert(t, (*ret2).(int) == testValue, "get elements from non empty buffer")
}
