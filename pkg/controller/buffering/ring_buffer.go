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

package buffering

import (
	"sync"

	log "github.com/Sirupsen/logrus"
)

// BufferSize the default size for the Ring Buffer
const BufferSize uint32 = 5000

// RingBuffer data structure
type RingBuffer struct {
	start, end, Size uint32
	buffer           [BufferSize]*interface{}
	Mutex            *sync.Mutex
}

// Put adds the item into buffer if there is room in buffer.
// Returns true if item is buffered otherwise false.
func (r *RingBuffer) Put(inp interface{}) bool {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if r.isFull() {
		return false
	}

	next := next(r.end, r.Size)
	r.buffer[r.end] = &inp
	r.end = next
	return true
}

// Get returns the elements in FIFO manner or nil if buffer is empty.
func (r *RingBuffer) Get() *interface{} {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if r.isEmpty() {
		return nil
	}

	next := next(r.start, r.Size)
	curval := r.buffer[r.start]
	r.buffer[r.start] = nil
	r.start = next
	return curval
}

// ReadN reads the next n available elements in the buffer.
// Returns elements and number of elements read.
func (r *RingBuffer) ReadN(n uint32) ([]*interface{}, uint32) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	var elements []*interface{}

	start := r.start
	for i := uint32(0); i < n; i++ {
		if start == r.end {
			break
		}
		elements = append(elements, r.buffer[start])
		start = next(start, r.Size)
	}
	return elements, uint32(len(elements))
}

// RemoveN removes the first n elements from the buffer.
func (r *RingBuffer) RemoveN(n uint32) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	start := r.start
	for i := uint32(0); i < n; i++ {
		if start == r.end {
			break
		}
		r.buffer[start] = nil
		start = next(start, r.Size)
		r.start = start
	}
}

func (r *RingBuffer) isEmpty() bool {
	return r.start == r.end
}

func (r *RingBuffer) isFull() bool {
	return next(r.end, r.Size) == r.start
}

func next(cur uint32, size uint32) uint32 {
	return (cur + 1) % size
}

// PrintDetails diplays details for debugging purpose.
func (r *RingBuffer) PrintDetails() {
	log.Printf("Start Position = %d, End Position = %d, Buffer Size = %d", r.start, r.end, r.Size)
}
