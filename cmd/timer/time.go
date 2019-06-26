package main

import (
	"time"
)

func main() {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
	}
}

//package main
//
//import (
//	"time"
//)
//
//func main() {
//	for {
//		time.Sleep(time.Millisecond)
//	}
//}
