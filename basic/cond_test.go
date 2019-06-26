package basic

import (
	"sync"
	"testing"
	"time"
)

func TestCond(t *testing.T) {
	rwMu := sync.RWMutex{}
	cond := sync.NewCond(rwMu.RLocker())

	go func() {
		rwMu.RLock()
		cond.Wait()
		rwMu.RUnlock()
		println("end1")
	}()
	go func() {
		rwMu.RLock()
		cond.Wait()
		rwMu.RUnlock()
		println("end2")
	}()
	go func() {
		rwMu.RLock()
		cond.Wait()
		rwMu.RUnlock()
		println("end3")
	}()
	go func() {
		rwMu.RLock()
		cond.Wait()
		rwMu.RUnlock()
		println("end4")
	}()
	go func() {
		rwMu.RLock()
		cond.Wait()
		rwMu.RUnlock()
		println("end5")
	}()
	time.Sleep(2 * time.Second)
	println("start")
	rwMu.RLock()
	cond.Broadcast()
	rwMu.RUnlock()

	println("end")
}
