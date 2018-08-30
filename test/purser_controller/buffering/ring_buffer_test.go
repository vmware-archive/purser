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
