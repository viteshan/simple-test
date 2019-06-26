package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"testing"

	"github.com/gorilla/websocket"
)

func TestCli2(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: "192.168.31.146:8080", Path: "/dev/websocket/viteshan"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"x-forwarded-for": {"192.168.31.47"}})
	//c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			select {
			case <-done:
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("recv: %s", message)
			}
		}
	}()
	go func() {
		defer close(done)
		i := 0
		for {
			i++
			select {
			case <-done:
				return
			default:
				c.WriteMessage(websocket.TextMessage, []byte("hello"+strconv.Itoa(i)))
			}
			time.Sleep(time.Second * 1)
		}
	}()

	ticker := time.NewTicker(1000 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			close(done)
		}
	}
}
