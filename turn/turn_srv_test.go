package turn

import (
	"github.com/pions/pkg/stun"
	"github.com/pions/turn"
	"testing"
)

type myTurnServer struct {
}

func (m *myTurnServer) AuthenticateRequest(username string, srcAddr *stun.TransportAddr) (password string, ok bool) {
	return password, true
}

func TestSrv1(t *testing.T) {
	m := &myTurnServer{}

	//realm := "43.241.213.212"
	realm := "10.2.16.76"
	udpPort := 9981

	turn.Start(turn.StartArguments{
		Server:  m,
		Realm:   realm,
		UDPPort: udpPort,
	})

}
