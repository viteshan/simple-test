package basic

import (
	"fmt"
	"testing"
	"time"
)

func TestRoutines(t *testing.T) {
	TestFor(t)
	N := 10

	for i := 0; i < N; i++ {
		go func() {
			b := make([]byte, 1024*1024)
			for j := 0; j < 1000; j++ {
				b[j] = byte(1)
			}
			runOneSec()
			fmt.Println("end point.N")
		}()
	}

	M := 1000
	for i := 0; i < M; i++ {
		go func() {
			b := make([]byte, 1024*1024)
			for j := 0; j < 1000000; j++ {
				b[j] = byte(1)
			}
			runOneSec()
			fmt.Println("end point.M")
		}()
	}

	is := make(chan int)
	is <- 10
}

func TestFor(t *testing.T) {
	t1 := time.Now() // get current time
	for i := 0; i < 10; i++ {
		//logic handlers
		runOneSec()
		elapsed := time.Since(t1)
		fmt.Println("App elapsed: ", elapsed)
	}

}

func runOneSec() {
	S := 1024 * 1024
	b := make([]byte, S)
	for i := 0; i < 1500; i++ {
		for j := 0; j < S; j++ {
			b[j] = byte(1)
		}
	}

	//for i := 0; i < 3400000000; i++ {
	//
	//}
}
