package datastructures

import (
	"expvar"
	"fmt"
	"net"
	"net/http"
	"testing"
)

func TestName(t *testing.T) {
	sock, err := net.Listen("tcp", "localhost:8123")
	if err != nil {
		panic("sock error")
	}
	go func() {
		fmt.Println("HTTP now available at port 8123")
		http.Serve(sock, nil)
	}()
	fmt.Println("hello")

	println(counts.String())
	select {}
}

var (
	counts = expvar.NewMap("counters")
	newInt = expvar.NewInt("daada")
)

func init() {
	counts.Add("a", 10)
	counts.Add("b", 10)
	newInt.Add(1)
}
