package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	//defer profile.Start(profile.MemProfile).Stop()
	testForTime()
	M := 10

	wg := sync.WaitGroup{}
	for i := 0; i < M; i++ {
		wg.Add(1)
		go func(i int) {
			runOneSec()
			fmt.Printf("end.M %d\n", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func testForTime() {
	t1 := time.Now() // get current time
	for i := 0; i < 10; i++ {
		//logic handlers
		runOneSec()
		elapsed := time.Since(t1)
		fmt.Println("App elapsed: ", elapsed)
	}
}

func runOneSec() {
	S := 1024 * 20
	b := make([]byte, S)
	for j := 0; j < S; j++ {
		b[j] = byte(1)
	}
	for i := 0; i < 100000000; i++ {

	}
}
