package turn

import (
	"log"
	"net"

	"fmt"
	"github.com/pions/turn"
	"github.com/willscott/goturn/client"
	"testing"
)

func TestCli(t *testing.T) {

	args := turn.ClientArguments{
		BindingAddress: "0.0.0.0:0",
		// IP and port of stun1.l.google.com
		//ServerIP:   net.IPv4(74, 125, 143, 127),
		//ServerPort: 19302,
		ServerIP: net.ParseIP("10.2.16.76"),
		ServerPort: 9981,
	}

	resp, err := turn.StartClient(args)

	if err != nil {
		panic(err)
	}

	log.Println(resp)

}

func TestCli3(t *testing.T) {

}

func TestCli2(t *testing.T) {
	// Connect to the stun/turn server
	conn, err := net.Dial("tcp", "10.2.16.76:9981")
	if err != nil {
		log.Fatal("error dialing TURN server: ", err)
	}
	defer conn.Close()

	//credentials := client.LongtermCredentials("username", "password")
	dialer, err := client.NewDialer(nil, conn)
	if err != nil {
		log.Fatal("failed to obtain dialer: ", err)
	}

	addr := dialer.LocalAddr
	fmt.Println(addr)
	//httpClient := &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
	//httpResp, err := httpClient.Get("http://www.google.com/")
	//if err != nil {
	//	log.Fatal("error performing http request: ", err)
	//}
	//defer httpResp.Body.Close()
	//
	//httpBody, err := ioutil.ReadAll(httpResp.Body)
	//if err != nil {
	//	log.Fatal("error reading http response: ", err)
	//}
	//log.Printf("received %d bytes", len(httpBody))
}
