package nat

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNat(t *testing.T) {
	srvPort := 9981
	srvAddr := "127.0.0.1"
	go srv(srvPort)
	go peer("peerA", 9983, srvPort, srvAddr)
	go peer("peerB", 9984, srvPort, srvAddr)

	c := make(chan int)
	c <- 1
}

func srv(srvPort int) {
	srvaddr, err := net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(srvPort))
	listener, err := net.ListenUDP("udp", srvaddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("local addr: <%s> \n", listener.LocalAddr().String())
	peers := make([]net.UDPAddr, 0, 2)
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])
		peers = append(peers, *remoteAddr)

		fmt.Println("remoteAddr:" + (*remoteAddr).String())
		if len(peers) == 2 {
			log.Printf("udp nat %s <--> %s\n", peers[0].String(), peers[1].String())
			listener.WriteToUDP([]byte(peers[1].String()), &peers[0])
			listener.WriteToUDP([]byte(peers[0].String()), &peers[1])
			time.Sleep(time.Second * 8)
			log.Println("exit.")
			return
		}
	}
}

const HAND_SHAKE_MSG = "hello hole"

func peer(name string, selfPort int, srvPort int, srvIp string) {
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: selfPort} // 注意端口必须固定
	dstAddr := &net.UDPAddr{IP: net.ParseIP(srvIp), Port: srvPort}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	if _, err = conn.Write([]byte("hello, I'm new peer:" + name)); err != nil {
		log.Panic(err)
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	conn.Close()
	anotherPeer := parseAddr(string(data[:n]))
	fmt.Printf("local:%s server:%s another:%s\n", srcAddr, remoteAddr, anotherPeer.String())
	// start hole
	bidirectionHole(srcAddr, &anotherPeer, name)

}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}
func bidirectionHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr, name string) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	// send udp hole message
	if _, err = conn.Write([]byte(HAND_SHAKE_MSG)); err != nil {
		log.Println("send handshake:", err)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from [" + name + "]udp")); err != nil {
				log.Println("send udp msg fail", err)
			}
		}
	}()


	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf(name+"recv data:%s\n", data[:n])
		}
	}
}
