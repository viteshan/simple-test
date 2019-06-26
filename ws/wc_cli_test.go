package ws

import (
	"log"
	"net/url"
	"strconv"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func TestCli2(t *testing.T) {
	origin := "*"
	//url := "ws://localhost:12345/ws"
	//ws, err := websocket.Dial(url, "", origin)
	//if err != nil {
	//	panic(err)
	//}
	//if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
	//	panic(err)
	//}
	//var msg = make([]byte, 512)
	//var n int
	//if n, err = ws.Read(msg); err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Received: %s.\n", msg[:n])
	//

	u := url.URL{Scheme: "ws", Host: "192.168.31.146:8080", Path: "/websocket/viteshan"}
	log.Printf("connecting to %s", u.String())

	c, err := websocket.Dial(u.String(), "", origin)
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
				var msg = make([]byte, 512)
				var n int
				if n, err = c.Read(msg); err != nil {
					panic(err)
				}

				log.Printf("recv: %s", msg[:n])
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
				c.Write([]byte("hello" + strconv.Itoa(i)))
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
