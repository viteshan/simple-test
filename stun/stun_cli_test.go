package stun

import (
	"fmt"
	"github.com/ccding/go-stun/stun"
	"log"
	"net"
	"testing"
	"time"
	//"github.com/libp2p/go-reuseport"
	"strconv"
)

func TestTcpConn(t *testing.T) {
	go func() {
		error := tcpConn("150.109.101.200", 8888, 9005)
		if error != nil {
			t.Error(error)
		}
	}()

	//go func() {
	//	error := tcpConn("127.0.0.1", 9002, 9001)
	//	if error != nil {
	//		t.Error(error)
	//	}
	//}()

	c := make(chan int)
	c <- 1
}

func tcpConn(peerPubAddr string, peerPubPort int, localPort int) error {
	peerAddr := &net.TCPAddr{IP: net.ParseIP(peerPubAddr), Port: peerPubPort} //peer addr
	localTcpAddr := &net.TCPAddr{IP: net.IPv4zero, Port: localPort}
	//send, er := reuseport.Dial("tcp", net.UDPAddr{IP: net.IPv4zero, Port: localPort}.String(), peerAddr.String())

	send, er := net.DialTCP("tcp", localTcpAddr, peerAddr)
	if er != nil {
		log.Printf("error connect: %s\n", er)
		return er
	}
	go func() {
		for {
			if _, err := send.Write([]byte("handshake message from " + peerPubAddr + ":" + strconv.Itoa(peerPubPort) + "\n")); err != nil {
				log.Println("send msg:", err)
			}
			time.Sleep(6 * time.Second)
		}
	}()
	for {
		time.Sleep(time.Second * 5)
		data := make([]byte, 1024)
		n, err := send.Read(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf("recv data:%s\n", data[:n])
		}
	}
}
func TestUdpConn(t *testing.T) {
	go func() {
		error := udpConn("111.204.124.34", 9000, 9001)
		if error != nil {
			t.Error(error)
		}
	}()

	//go func() {
	//	error := udpConn("127.0.0.1", 9002, 9001)
	//	if error != nil {
	//		t.Error(error)
	//	}
	//}()

	c := make(chan int)
	c <- 1

}

func udpConn(peerPubAddr string, peerPubPort int, localPort int) error {
	peerAddr := &net.UDPAddr{IP: net.ParseIP(peerPubAddr), Port: peerPubPort} //peer addr
	localUdpAddr := &net.UDPAddr{IP: net.IPv4zero, Port: localPort}
	//send, er := reuseport.Dial("tcp", net.UDPAddr{IP: net.IPv4zero, Port: localPort}.String(), peerAddr.String())

	send, er := net.DialUDP("udp", localUdpAddr, peerAddr)
	if er != nil {
		log.Printf("error connect: %s\n", er)
		return er
	}
	go func() {
		for {
			if _, err := send.Write([]byte("handshake message from " + peerPubAddr + ":" + strconv.Itoa(peerPubPort) + "\n")); err != nil {
				log.Println("send msg:", err)
			}
			time.Sleep(6 * time.Second)
		}
	}()
	for {
		time.Sleep(time.Second * 5)
		data := make([]byte, 1024)
		n, err := send.Read(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf("recv data:%s\n", data[:n])
		}
	}
}

func TestPubAddr(t *testing.T) {
	localPort := 9000
	pubAddr, e := pubAddr(localPort)
	if e != nil {
		t.Error(e)
	}
	fmt.Println(pubAddr)
}
func pubAddr(localPort int) (string, error) {
	localAddr := &net.UDPAddr{IP: net.IPv4zero, Port: localPort}
	srvAddr := "stun.freeswitch.org:3478" //stun srv
	conn, _ := net.ListenUDP("udp", localAddr)
	//conn, _ := reuseport.ListenPacket("udp", localAddr.String())
	client := stun.NewClientWithConnection(conn)
	client.SetServerAddr(srvAddr)
	nat, host, err := client.Discover()

	fmt.Printf("nat info: [%v],[%v],[%v]\n", nat, host, err)
	if err != nil {
		return "", err
	}
	return host.String(), nil
}
