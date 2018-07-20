package buffering

import "sync"

const BUFFER_SIZE uint32 = 1000

type RingBuffer struct {
	start, end, size uint32
	buffer [BUFFER_SIZE]*interface{}
	mutex *sync.Mutex
}

/*
 * Puts the item into buffer if there is room in buffer.
 * returns true if item is buffered otherwise false.
 */
func(r *RingBuffer) Put(inp *interface{}) bool {
	ret := false
	r.mutex.Lock()

	if r.isFull() {
		ret = false
	} else {
		next := next(r.end, r.size)
		r.buffer[next] = inp
		r.start = next
		ret = true
	}
	r.mutex.Unlock()
	return ret
}

/*
 * Returns the elements in FIFO manner.
 * If buffer is empty then it returns nil.
 */
func(r *RingBuffer) Get() *interface{} {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.isEmpty() {
		return nil
	} else {
		next := next(r.start, r.size)
		curval := r.buffer[r.start]
		r.buffer[r.start] = nil
		r.start = next
		return curval
	}
}

/*
 * Reads the next n available elements in the buffer.
 * Returns elements and number of elements read.
 */
func(r *RingBuffer) ReadN(n uint32) ([]*interface{}, uint32){
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var elements []*interface{}
	i := uint32(0)
	start := r.start
	for i < n {
		if start == r.end {
			break
		}
		elements = append(elements, r.buffer[start])
		start = next(start, r.size)
		i++
	}
	return elements, uint32(len(elements))
}

/*
 * Removes the first n elements from the buffer.
 */
func(r *RingBuffer) RemoveN(n uint32) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	i := uint32(0)
	start := r.start
	for i < n {
		if start == r.end {
			break
		}
		r.buffer[start] = nil
		start = next(start, r.size)
		r.start = start
		i++
	}
}

func(r *RingBuffer) isEmpty() bool {
	return r.start == r.end
}

func(r *RingBuffer) isFull() bool {
	return next(r.end, r.size) == r.start
}

func next(cur uint32, size uint32) uint32 {
	return (cur + 1) % size
}