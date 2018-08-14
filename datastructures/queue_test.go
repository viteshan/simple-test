package datastructures

import (
	"sync"
	"testing"
	"time"

	"github.com/Workiva/go-datastructures/queue"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {

	q := queue.New(10)

	// should be able to Poll() before anything is present, without breaking future Puts
	q.Poll(1, time.Millisecond)

	q.Put(`test`)
	result, err := q.Poll(2, 0)
	if !assert.Nil(t, err) {
		return
	}

	assert.Len(t, result, 1)
	assert.Equal(t, `test`, result[0])
	assert.Equal(t, int64(0), q.Len())

	q.Put(`1`)
	q.Put(`2`)

	q.Peek()

	result, err = q.Poll(1, time.Millisecond)
	if !assert.Nil(t, err) {
		return
	}

	assert.Len(t, result, 1)
	assert.Equal(t, `1`, result[0])
	assert.Equal(t, int64(1), q.Len())

	result, err = q.Poll(2, time.Millisecond)
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, `2`, result[0])

	before := time.Now()
	_, err = q.Poll(1, 5*time.Millisecond)
	// This delta is normally 1-3 ms but running tests in CI with -race causes
	// this to run much slower. For now, just bump up the threshold.
	assert.InDelta(t, 5, time.Since(before).Seconds()*1000, 10)
	assert.Equal(t, queue.ErrTimeout, err)
}

func TestRingBuffer(t *testing.T) {
	buffer := queue.NewRingBuffer(10)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_, e := buffer.Get()
		if e == nil {
			t.Error("dispose error.")
		}
		e.Error()
		wg.Done()
	}()

	buffer.Dispose()
	wg.Wait()

	buffer = queue.NewRingBuffer(10)
	buffer.Put(1)

	_, e := buffer.Get()
	if e != nil {
		t.Error("error for get logic.")
	}
	if buffer.Len() != 0 {
		t.Error("error len for buffer.", buffer.Len())
	}

	buffer.Put(1)
	buffer.Put(2)
	if buffer.Len() != 2 {
		t.Error("error len for buffer.", buffer.Len())
	}

	i, e := buffer.Get()
	if e != nil {
		t.Error("error for get logic.")
	}
	if i.(int) != 1 {
		t.Error("number should be 1.", i)
	}
	i, e = buffer.Get()
	if e != nil {
		t.Error("error for get logic.")
	}
	if i.(int) != 2 {
		t.Error("number should be 2.", i)
	}
}
